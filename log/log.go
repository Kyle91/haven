// @Author Eric
// @Date 2024/6/2 17:36:00
// @Desc 日志库，支持并发，自动轮替
package log

import (
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
	sync.Mutex
	file        *os.File
	serviceName string
	maxSize     int64
	maxBackups  int
	logDir      string
	currSize    int64
}

var logger *Logger

// Init initializes the logger
func Init(serviceName, logDir string, maxSize int64, maxBackups int) error {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(logDir, "logfile.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logger = &Logger{
		file:        file,
		serviceName: serviceName,
		maxSize:     maxSize,
		maxBackups:  maxBackups,
		logDir:      logDir,
	}

	// Get the initial size of the log file
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}
	logger.currSize = stat.Size()

	return nil
}

// log logs the message with the specified level
func (l *Logger) log(level, msg string) {
	l.Lock()
	defer l.Unlock()

	// Check if we need to rotate the log
	if l.currSize >= l.maxSize {
		l.rotateLogs()
	}

	// Prepare the log entry
	entry := fmt.Sprintf("%s [%s] %s %s:%d %s\n",
		time.Now().Format(time.RFC3339Nano), level, l.serviceName, getFuncName(), getLine(), msg)

	// Write to the log file
	n, err := l.file.WriteString(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to log file: %v\n", err)
	}

	// Write to the console
	fmt.Print(entry)

	l.currSize += int64(n)
}

// rotateLogs rotates the log files
func (l *Logger) rotateLogs() {
	l.file.Close()

	for i := l.maxBackups - 1; i >= 0; i-- {
		oldPath := filepath.Join(l.logDir, fmt.Sprintf("logfile.log.%d", i))
		newPath := filepath.Join(l.logDir, fmt.Sprintf("logfile.log.%d", i+1))
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	os.Rename(filepath.Join(l.logDir, "logfile.log"), filepath.Join(l.logDir, "logfile.log.0"))

	file, err := os.OpenFile(filepath.Join(l.logDir, "logfile.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open new log file: %v\n", err)
		return
	}

	l.file = file
	l.currSize = 0
}

// getFuncName returns the name of the function that called the logger
func getFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
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
func getLine() int {
	_, _, line, ok := runtime.Caller(2)
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
	os.Exit(1)
}
