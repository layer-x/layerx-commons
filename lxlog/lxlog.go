package lxlog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

var log = logrus.New()

func Debugf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	if len(a) > 0 {
		log.WithFields(fields).Debugf(format, a)
	} else {
		log.WithFields(fields).Debug(format)
	}
}

func Infof(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	if len(a) > 0 {
		log.WithFields(fields).Infof(format, a)
	} else {
		log.WithFields(fields).Info(format)
	}
}

func Warnf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Warnf(format, a)
}

func Errorf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Errorf(format, a)
}

func Fatalf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Fatalf(format, a)
}

func Panicf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Panicf(format, a)
}

func addLine(fields logrus.Fields) logrus.Fields {
	pc, fn, line, _ := runtime.Caller(2)
	pathComponents := strings.Split(fn, "/")
	var truncatedPath string
	if len(pathComponents) > 3 {
		truncatedPath = strings.Join(pathComponents[len(pathComponents)-2:], "/")
	} else {
		truncatedPath = strings.Join(pathComponents, "/")
	}
	fnName := runtime.FuncForPC(pc).Name()
	fnNameComponents :=strings.Split(fnName, "/")
	truncatedFnName := fnNameComponents[len(fnNameComponents)-1]

	file := fmt.Sprintf("%s[%s:%d]", truncatedFnName, truncatedPath, line)
	fields["file"] = file
	return fields
}
