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

var log = logrus.New()
var optionalLogs = make(map[string]*logrus.Logger)

func ActiveDebugMode() {
	log.Level = logrus.DebugLevel
}

func AddLogger(name string, logLevel logrus.Level, w io.Writer) {
	newLogger := logrus.New()
	newLogger.Out = w
	newLogger.Level = logLevel
	optionalLogs[name] = newLogger
}

func LogCommand(cmd *exec.Cmd, debug bool) {
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
			if debug {
				Debugf(logrus.Fields{}, in.Text())
			} else {
				Infof(logrus.Fields{}, in.Text())
			}
		}
	}()
	go func() {
		// read command's stdout line by line
		in := bufio.NewScanner(stderr)

		for in.Scan() {
			Debugf(logrus.Fields{}, in.Text())
		}
	}()
}

func DeleteLogger(name string) {
	delete(optionalLogs, name)
}

func Debugf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Debugf(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Debugf(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Debug(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Debug(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func Infof(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Infof(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Infof(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Info(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Info(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func Warnf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Warnf(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Warnf(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Warn(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Warn(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func Errorf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Errorf(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Errorf(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Error(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Error(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func Fatalf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Fatalf(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Fatalf(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Fatal(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Fatal(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func Panicf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Panicf(format, a)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Panicf(format, a)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	} else {
		log.WithFields(fields).Panic(format)
		for _, optionalLog := range optionalLogs {
			optionalLog.WithFields(fields).Panic(format)
			if flusher, ok := optionalLog.Out.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func addTrace(format string) string {
	pc, fn, line, _ := runtime.Caller(2)
	pathComponents := strings.Split(fn, "/")
	var truncatedPath string
	if len(pathComponents) > 3 {
		truncatedPath = strings.Join(pathComponents[len(pathComponents)-2:], "/")
	} else {
		truncatedPath = strings.Join(pathComponents, "/")
	}
	fnName := runtime.FuncForPC(pc).Name()
	fnNameComponents := strings.Split(fnName, "/")
	truncatedFnName := fnNameComponents[len(fnNameComponents)-1]

	file := fmt.Sprintf("%s[%s:%d] ", truncatedFnName, truncatedPath, line)

	return file + format
}
