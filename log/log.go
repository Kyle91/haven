// @Author Eric
// @Date 2024/6/2 17:36:00
// @Desc 在 logrus 的基础上进行二次封装
package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime"
)

type Log struct {
	*logrus.Logger
	serviceName string
}

// NewLogger 创建一个新的 Logger 实例,默认日志级别为 Info
func NewLogger(serviceName string) *Log {
	logger := logrus.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}

	logDir := filepath.Join(homeDir, "haven", "log")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	logPath := filepath.Join(logDir, "haven.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   logPath, // 日志文件的位置
		MaxSize:    20,      // 兆为单位
		MaxBackups: 10,      //最大备份数量
		MaxAge:     28,      //最大保存时间（以天为单位）
		Compress:   true,    // 是否压缩日志
	})

	return &Log{
		Logger:      logger,
		serviceName: serviceName,
	}
}

// WithContext 添加上下文信息到日志条目中
func (l *Log) WithContext() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	var fnName string
	if ok {
		fnName = runtime.FuncForPC(pc).Name()
	}
	return l.WithFields(logrus.Fields{
		"service": l.serviceName,
		"file":    file,
		"line":    line,
		"func":    fnName,
	})
}

// Infof 包装 logrus 的 Infof 方法，添加上下文信息
func (l *Log) Infof(format string, args ...interface{}) {
	l.WithContext().Infof(format, args...)
}

// Errorf 包装 logrus 的 Errorf 方法，添加上下文信息
func (l *Log) Errorf(format string, args ...interface{}) {
	l.WithContext().Errorf(format, args...)
}

// Debugf 包装 logrus 的 Debugf 方法，添加上下文信息
func (l *Log) Debugf(format string, args ...interface{}) {
	l.WithContext().Debugf(format, args...)
}

// Warningf 包装 logrus 的 Warningf 方法，添加上下文信息
func (l *Log) Warningf(format string, args ...interface{}) {
	l.WithContext().Warningf(format, args...)
}

// Tracef 包装 logrus 的 Tracef 方法，添加上下文信息
func (l *Log) Tracef(format string, args ...interface{}) {
	l.WithContext().Tracef(format, args...)
}
