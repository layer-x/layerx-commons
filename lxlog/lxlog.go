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

var GlobalLogLevel Level = InfoLevel

const (
	default_logger = "default_logger"
	default_trace = 3

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

type Logger interface {
	WithFields(fields Fields) Logger
	WithErr(err error) Logger
	SetLogLevel(level Level)
	AddWriter(name string, level Level, w io.Writer)
	DeleteWriter(name string)
	LogCommand(cmd *exec.Cmd, asDebug bool)
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
	Panicf(format string, a ...interface{})
}

type lxLogger struct {
	loggers map[string]*logrus.Logger
	fields  Fields
	err     error
	name	string
	trace   int
}

func New(name string) Logger {
	loggers := make(map[string]*logrus.Logger)
	loggers[default_logger] = logrus.New()
	lxlogger := &lxLogger{
		loggers: loggers,
		name: name,
		trace: 0,
	}
	lxlogger.SetLogLevel(GlobalLogLevel)
	return lxlogger
}

func (lxlog *lxLogger) WithFields(fields Fields) Logger {
	return &lxLogger{
		loggers: lxlog.loggers,
		fields: fields,
		err: lxlog.err,
		name: lxlog.name,
		trace: lxlog.trace,
	}
}

func (lxlog *lxLogger) WithErr(err error) Logger {
	return &lxLogger{
		loggers: lxlog.loggers,
		fields: lxlog.fields,
		err: err,
		name: lxlog.name,
		trace: lxlog.trace,
	}
}

func (lxlog *LxLogger) WithTrace(trace int) *LxLogger {
	return &LxLogger{
		loggers: lxlog.loggers,
		fields: lxlog.fields,
		err: lxlog.err,
		name: lxlog.name,
		trace: trace,
	}
}

func (lxlog *lxLogger) SetLogLevel(level Level) {
	for _, logrusLogger := range lxlog.loggers {
		logrusLogger.Level = logLevels[level]
	}
}

func (lxlog *lxLogger) AddWriter(name string, level Level, w io.Writer) {
	newLogger := logrus.New()
	newLogger.Out = w
	newLogger.Level = logLevels[level]
	lxlog.loggers[name] = newLogger
}

func (lxlog *lxLogger) DeleteWriter(name string) {
	delete(lxlog.loggers, name)
}

func (lxlog *lxLogger) LogCommand(cmd *exec.Cmd, asDebug bool) {
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

func (lxlog *lxLogger) Infof(format string, a ...interface{}) {
	lxlog.log(InfoLevel, format, a...)
}

func (lxlog *lxLogger) Debugf(format string, a ...interface{}) {
	lxlog.log(DebugLevel, format, a...)
}

func (lxlog *lxLogger) Warnf(format string, a ...interface{}) {
	lxlog.log(WarnLevel, format, a...)
}

func (lxlog *lxLogger) Errorf(format string, a ...interface{}) {
	lxlog.log(ErrorLevel, format, a...)
}

func (lxlog *lxLogger) Fatalf(format string, a ...interface{}) {
	lxlog.log(FatalLevel, format, a...)
}

func (lxlog *lxLogger) Panicf(format string, a ...interface{}) {
	lxlog.log(PanicLevel, format, a...)
}

func (lxlog *lxLogger) log(level Level, format string, a ...interface{}) {
	format = lxlog.addTrace(format)
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

func (lxlog *LxLogger) addTrace(format string) string {
	pc, fn, line, _ := runtime.Caller(default_trace+lxlog.trace)
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

	file := fmt.Sprintf("(%s): %s[%s:%d] ", lxlog.name, truncatedFnName, truncatedPath, line)

	return file + format
}
