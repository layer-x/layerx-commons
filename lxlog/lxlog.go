package lxlog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

var log = logrus.New()

func ActiveDebugMode() {
	log.Level = logrus.DebugLevel
}

func Debugf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Debugf(format, a)
	} else {
		log.WithFields(fields).Debug(format)
	}
}

func Infof(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Infof(format, a)
	} else {
		log.WithFields(fields).Info(format)
	}
}

func Warnf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Warnf(format, a)
	} else {
		log.WithFields(fields).Warn(format)
	}
}

func Errorf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Errorf(format, a)
	} else {
		log.WithFields(fields).Error(format)
	}
}

func Fatalf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Fatalf(format, a)
	} else {
		log.WithFields(fields).Fatal(format)
	}
}

func Panicf(fields logrus.Fields, format string, a ...interface{}) {
	format = addTrace(format)
	if len(a) > 0 {
		log.WithFields(fields).Panicf(format, a)
	} else {
		log.WithFields(fields).Panic(format)
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
