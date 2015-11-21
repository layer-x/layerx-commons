package lxlog
import (
	"github.com/Sirupsen/logrus"
	"fmt"
	"runtime"
)

var log = logrus.New()

func Debugf(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Debugf(format, a)
}

func Infof(fields logrus.Fields, format string, a ...interface{}) {
	fields = addLine(fields)
	log.WithFields(fields).Infof(format, a)
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
	file := fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fn, line)
	fields["file"] = file
	return fields
}

func main(){
	Infof(logrus.Fields{
		"val1": 5,
		"val2": true,
		"val3": "this",
	}, "Here's some %s info", "special")
	Warnf(logrus.Fields{
		"val1": 5,
		"val2": true,
		"val3": "this",
	}, "Here's some %s info", "special")
	Errorf(logrus.Fields{
		"val1": 5,
		"val2": true,
		"val3": "this",
	}, "Here's some %s info", "special")
}