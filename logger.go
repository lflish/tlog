package tlog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"syscall"
)

type Log struct {
	Level  string `yaml:"Level"`
	Output string `yaml:"Output"`
}

var level = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}

type Option struct {
	level string
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

func SetOption(level, path string) {

	option.level = level
	option.path = path

	logWriter = &lumberjack.Logger{
		// 日志输出文件路径
		Filename: option.path,
		// 日志文件最大 size, 单位是 MB
		MaxSize: 500, // megabytes
		// 最大过期日志保留的个数
		MaxBackups: 10,
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge: 28, //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: false, // disabled by default
	}

	stdout, err := os.OpenFile(path+".out", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0664)
	if err != nil {
		fmt.Printf("create stdout %s err:%s", path+".out", err)
	}

	syscall.Dup2(int(stdout.Fd()), 1)
	syscall.Dup2(int(stdout.Fd()), 2)

}

func (f *Logger) ErrorfEx(ctx context.Context, format string, args ...interface{}) error {
	f.Errorf(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (f *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	f.Logger.Errorf(format, args...)
}

func (f *Logger) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	var newLog string
	newLog = fmt.Sprintf("%s [%s] [%s] %s\n", timestamp, entry.Level, f.mod, getTraceId(entry.Context), entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

var loggersMu sync.Mutex

var loggers = map[string]*Logger{}

func GetLogger(mod string) *Logger {
	loggersMu.Lock()

	l, ok := loggers[mod]
	if !ok {
		l = &Logger{
			mod:    mod,
			Logger: logrus.New(),
		}

		if option.path != "" {
			l.SetOutput(logWriter) //调用 logrus 的 SetOutput()函数
		}

		l.Logger.SetFormatter(l)
		l.SetLevel(option.Level())

		loggers[mod] = l
	}
	loggersMu.Unlock()

	return l
}
