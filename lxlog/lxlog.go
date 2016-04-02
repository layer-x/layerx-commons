package lxlog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
	"io"
	"net/http"
	"os/exec"
	"bufio"
)

const (
	default_logger = "default_logger"

	PanicLevel = Level("PanicLevel")
	FatalLevel = Level("FatalLevel")
	ErrorLevel = Level("ErrorLevel")
	WarnLevel = Level("WarnLevel")
	InfoLevel = Level("InfoLevel")
	DebugLevel = Level("DebugLevel")
)

var logLevels = map[Level]logrus.Level{
	PanicLevel: logrus.PanicLevel,
	FatalLevel: logrus.FatalLevel,
	ErrorLevel: logrus.ErrorLevel,
	WarnLevel: logrus.WarnLevel,
	InfoLevel: logrus.InfoLevel,
	DebugLevel: logrus.DebugLevel,
}

type Level string

func (level Level) String() string {
	return string(level)
}

type Fields logrus.Fields

type LxLogger struct {
	loggers map[string]*logrus.Logger
	fields  Fields
	err     error
}

func New() *LxLogger {
	loggers := make(map[string]*logrus.Logger)
	loggers[default_logger] = logrus.New()
	return &LxLogger{
		loggers: loggers,
	}
}

func (lxlog *LxLogger) WithFields(fields Fields) *LxLogger {
	return &LxLogger{
		loggers: lxlog.loggers,
		fields: fields,
		err: lxlog.err,
	}
}

func (lxlog *LxLogger) WithErr(err error) *LxLogger {
	return &LxLogger{
		loggers: lxlog.loggers,
		fields: lxlog.fields,
		err: err,
	}
}

func (lxlog *LxLogger) SetLogLevel(level Level) {
	for _, logrusLogger := range lxlog.loggers {
		logrusLogger.Level = logLevels[level]
	}
}

func (lxlog *LxLogger) AddWriter(name string, level Level, w io.Writer) {
	newLogger := logrus.New()
	newLogger.Out = w
	newLogger.Level = logLevels[level]
	lxlog.loggers[name] = newLogger
}

func (lxlog *LxLogger) DeleteWriter(name string) {
	delete(lxlog.loggers, name)
}

func (lxlog *LxLogger) LogCommand(cmd *exec.Cmd, asDebug bool) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	go func() {
		// read command's stdout line by line
		in := bufio.NewScanner(stdout)

		for in.Scan() {
			if asDebug {
				lxlog.Debugf(in.Text())
			} else {
				lxlog.Infof(in.Text())
			}
		}
	}()
	go func() {
		// read command's stdout line by line
		in := bufio.NewScanner(stderr)

		for in.Scan() {
			lxlog.Errorf(in.Text())
		}
	}()
}

func (lxlog *LxLogger) Infof(format string, a ...interface{}) {
	lxlog.log(InfoLevel, format, a...)
}

func (lxlog *LxLogger) Debugf(format string, a ...interface{}) {
	lxlog.log(DebugLevel, format, a...)
}

func (lxlog *LxLogger) Warnf(format string, a ...interface{}) {
	lxlog.log(WarnLevel, format, a...)
}

func (lxlog *LxLogger) Errorf(format string, a ...interface{}) {
	lxlog.log(ErrorLevel, format, a...)
}

func (lxlog *LxLogger) Fatalf(format string, a ...interface{}) {
	lxlog.log(FatalLevel, format, a...)
}

func (lxlog *LxLogger) Panicf(format string, a ...interface{}) {
	lxlog.log(PanicLevel, format, a...)
}

func (lxlog *LxLogger) log(level Level, format string, a ...interface{}) {
	format = addTrace(format)
	for _, optionalLog := range lxlog.loggers {
		entry := optionalLog.WithFields(logrus.Fields(lxlog.fields))
		if lxlog.err != nil {
			entry = entry.WithError(lxlog.err)
		}
		switch level {
		case PanicLevel:
			entry.Panicf(format, a...)
			break
		case FatalLevel:
			entry.Fatalf(format, a...)
			break
		case ErrorLevel:
			entry.Errorf(format, a...)
			break
		case WarnLevel:
			entry.Warnf(format, a...)
			break
		case InfoLevel:
			entry.Infof(format, a...)
			break
		case DebugLevel:
			entry.Debugf(format, a...)
			break
		}
		if flusher, ok := optionalLog.Out.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

func addTrace(format string) string {
	pc, fn, line, _ := runtime.Caller(3)
	pathComponents := strings.Split(fn, "/")
	var truncatedPath string
	if len(pathComponents) > 3 {
		truncatedPath = strings.Join(pathComponents[len(pathComponents) - 2:], "/")
	} else {
		truncatedPath = strings.Join(pathComponents, "/")
	}
	fnName := runtime.FuncForPC(pc).Name()
	fnNameComponents := strings.Split(fnName, "/")
	truncatedFnName := fnNameComponents[len(fnNameComponents) - 1]

	file := fmt.Sprintf("%s[%s:%d] ", truncatedFnName, truncatedPath, line)

	return file + format
}
