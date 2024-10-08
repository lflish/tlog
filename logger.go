package pailog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"sync"
)

type LEVEL int

const (
	TRACE = iota
	DEBUG
	INFO
	ERROR
	PANIC
	FATAL
)

type Log struct {
	Level  string `yaml:"Level"`
	Output string `yaml:"Output"`
}

var level = map[LEVEL]logrus.Level{
	TRACE: logrus.TraceLevel,
	DEBUG: logrus.DebugLevel,
	INFO:  logrus.InfoLevel,
	ERROR: logrus.ErrorLevel,
	PANIC: logrus.PanicLevel,
	FATAL: logrus.FatalLevel,
}

type Option struct {
	level LEVEL
	path  string
}

func (o *Option) Level() logrus.Level {
	return level[o.level]
}

var option Option

type Logger struct {
	mod string
	*logrus.Logger
}

var logWriter *lumberjack.Logger

func SetOption(level LEVEL, path string) {

	option.level = level
	option.path = path

	logWriter = &lumberjack.Logger{
		// 日志文件最大 size, 单位是 MB
		MaxSize: 500, // megabytes
		// 最大过期日志保留的个数
		MaxBackups: 10,
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge: 28, //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: false, // disabled by default
	}

	loggersMu.Lock()
	defer loggersMu.Unlock()

	if option.path != "" {
		logWriter.Filename = option.path
	}

	for _, v := range loggers {
		v.SetOutput(logWriter)
		v.Logger.SetFormatter(v)
		v.SetLevel(option.Level())
		v.SetReportCaller(true)
	}
}

/****************** ctx ********************************/
func (f *Logger) ErrorfEx(ctx context.Context, format string, args ...interface{}) error {
	f.ErrorF(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (f *Logger) ErrorF(ctx context.Context, format string, args ...interface{}) {
	f.WithContext(ctx).Errorf(format, args...)
}

func (f *Logger) DebugF(ctx context.Context, format string, args ...interface{}) {
	f.WithContext(ctx).Debugf(format, args...)
}

func (f *Logger) InfoF(ctx context.Context, format string, args ...interface{}) {
	f.WithContext(ctx).Infof(format, args...)
}

func (f *Logger) PanicF(ctx context.Context, format string, args ...interface{}) {
	f.WithContext(ctx).Panicf(format, args...)
}

func (f *Logger) FatalF(ctx context.Context, format string, args ...interface{}) {
	f.WithContext(ctx).Fatalf(format, args...)
}

/**********************************************/
/*
func (f *Logger) Errorf(format string, args ...interface{}) {
	f.Errorf(format, args...)
}

func (f *Logger) Debugf(format string, args ...interface{}) {
	f.Debugf(format, args...)
}

func (f *Logger) Infof(format string, args ...interface{}) {
	f.Infof(format, args...)
}

func (f *Logger) Panicf(format string, args ...interface{}) {
	f.Panicf(format, args...)
}

func (f *Logger) Fatalf(format string, args ...interface{}) {
	f.Fatalf(format, args...)
}*/

func (f *Logger) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	var newLog string

	if entry.HasCaller() {
		fileName := filepath.Base(entry.Caller.File)
		if logWriter.Filename == "" {
			fmt.Printf("%s [%s] [%s:%d] [%s] [%s] %s\n", timestamp, entry.Level, fileName, entry.Caller.Line, f.mod, getTraceId(entry.Context), entry.Message)
		} else if getTraceId(entry.Context) == "nil" {
			newLog = fmt.Sprintf("%s [%s] [%s:%d] [%s] %s\n", timestamp, entry.Level, fileName, entry.Caller.Line, f.mod, entry.Message)
			b.WriteString(newLog)
		} else {
			newLog = fmt.Sprintf("%s [%s] [%s:%d] [%s] [%s] %s\n", timestamp, entry.Level, fileName, entry.Caller.Line, f.mod, getTraceId(entry.Context), entry.Message)
			b.WriteString(newLog)
		}
	} else {
		if logWriter.Filename == "" {
			fmt.Printf("%s [%s] [%s] [%s] %s\n", timestamp, entry.Level, f.mod, getTraceId(entry.Context), entry.Message)
		} else if getTraceId(entry.Context) == "nil" {
			newLog = fmt.Sprintf("%s [%s] [%s] %s\n", timestamp, entry.Level, f.mod, entry.Message)
			b.WriteString(newLog)
		} else {
			newLog = fmt.Sprintf("%s [%s] [%s] [%s] %s\n", timestamp, entry.Level, f.mod, getTraceId(entry.Context), entry.Message)
			b.WriteString(newLog)
		}
	}

	return b.Bytes(), nil
}

var loggersMu sync.Mutex

var loggers = map[string]*Logger{}

func GetLogger(mod string) *Logger {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	l, ok := loggers[mod]
	if !ok {
		l = &Logger{
			mod:    mod,
			Logger: logrus.New(),
		}
		loggers[mod] = l
	}

	return l
}
