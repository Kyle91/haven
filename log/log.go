// @Author Eric
// @Date 2024/6/2 17:36:00
// @Desc 高性能无锁日志，支持并发，自动轮替，utf-8编码存储
package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger struct holds the logger configuration and state
type Logger struct {
	file        *os.File
	writer      *bufio.Writer
	logCh       chan string // channel存储数据
	maxSize     int64
	maxBackups  int
	logDir      string
	currSize    int64
	serviceName string
	wg          sync.WaitGroup
	rotateLock  sync.Mutex
}

var logger *Logger

func init() {
	logDir := getDefaultLogDir()
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create log directory: %w", err))
	}

	logFile := filepath.Join(logDir, "haven.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open log file: %w", err))
	}

	writer := bufio.NewWriter(file)

	logger = &Logger{
		file:        file,
		writer:      writer,
		logCh:       make(chan string, 5000),
		maxSize:     100 * 1024 * 1024, // 100 MB
		maxBackups:  10,
		logDir:      logDir,
		serviceName: getServiceName(),
	}

	// Get the initial size of the log file
	stat, err := file.Stat()
	if err != nil {
		panic(fmt.Errorf("failed to stat log file: %w", err))
	}
	logger.currSize = stat.Size()

	go logger.processLogEntries()

	// Ensure logs are written to the file in case of panic
	defer func() {
		if r := recover(); r != nil {
			logger.flush()
		}
	}()
}

// getDefaultLogDir returns the default log directory based on the operating system
func getDefaultLogDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "haven", "log")
	}
	return filepath.Join(os.Getenv("HOME"), "haven", "log")
}

// getServiceName returns the name of the service by using the process name
func getServiceName() string {
	parts := strings.Split(os.Args[0], string(os.PathSeparator))
	return parts[len(parts)-1]
}

// log logs the message with the specified level
func (l *Logger) log(level, msg string) {
	// Prepare the log entry
	entry := fmt.Sprintf("%s [%s] %s %s:%d %s\n",
		time.Now().Format(time.RFC3339Nano), level, l.serviceName, getFuncName(4), getLine(4), msg)

	// Try to send the log entry to the log channel
	select {
	case l.logCh <- entry:
	default:
		// If the channel is full, write directly to the file
		l.writeLogEntry(entry)
	}
}

func (l *Logger) logf(level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(level, msg)
}

// processLogEntries processes log entries from the log channel
func (l *Logger) processLogEntries() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case entry := <-l.logCh:
			l.writeLogEntry(entry)
		case <-ticker.C:
			l.flush()
		}
	}
}

// writeLogEntry writes a log entry to the log file
func (l *Logger) writeLogEntry(entry string) {
	l.wg.Add(1)
	defer l.wg.Done()

	// Rotate log if necessary
	l.checkRotate()

	// Write to the log file
	n, err := l.writer.WriteString(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to log file: %v\n", err)
	}

	// Write to the console
	fmt.Print(entry)

	l.currSize += int64(n)
}

// checkRotate checks if log rotation is needed and performs the rotation if necessary
func (l *Logger) checkRotate() {
	l.rotateLock.Lock()
	defer l.rotateLock.Unlock()

	if l.currSize >= l.maxSize {
		l.rotateLogs()
	}
}

// rotateLogs rotates the log files
func (l *Logger) rotateLogs() {
	l.flush()
	l.file.Close()

	timestamp := time.Now().Format("20060102150405")
	for i := l.maxBackups - 1; i >= 0; i-- {
		oldPath := filepath.Join(l.logDir, fmt.Sprintf("haven.log.%d", i))
		newPath := filepath.Join(l.logDir, fmt.Sprintf("haven.log.%d", i+1))
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	oldLog := filepath.Join(l.logDir, "haven.log")
	newLog := filepath.Join(l.logDir, fmt.Sprintf("haven.log.%s.0", timestamp))
	os.Rename(oldLog, newLog)

	file, err := os.OpenFile(oldLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open new log file: %v\n", err)
		return
	}

	l.file = file
	l.writer = bufio.NewWriter(file)
	l.currSize = 0
}

// flush flushes the log writer
func (l *Logger) flush() {
	l.wg.Wait()
	l.writer.Flush()
}

// getFuncName returns the name of the function that called the logger
func getFuncName(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	parts := strings.Split(fn.Name(), "/")
	return parts[len(parts)-1]
}

// getLine returns the line number where the logger was called
func getLine(depth int) int {
	_, _, line, ok := runtime.Caller(depth)
	if !ok {
		return 0
	}
	return line
}

// Exported logging functions

func Debug(msg string) {
	logger.log("DEBUG", msg)
}

func Info(msg string) {
	logger.log("INFO", msg)
}

func Warn(msg string) {
	logger.log("WARN", msg)
}

func Error(msg string) {
	logger.log("ERROR", msg)
}

func Fatal(msg string) {
	logger.log("FATAL", msg)
	logger.flush()
	os.Exit(1)
}

func Debugf(format string, args ...interface{}) {
	logger.logf("DEBUG", format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.logf("INFO", format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.logf("WARN", format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.logf("ERROR", format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.logf("FATAL", format, args...)
	logger.flush()
	os.Exit(1)
}
