package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	debugMode = false
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func SetDebug(enabled bool) {
	debugMode = enabled
}

func Debug(args ...interface{}) {
	if debugMode {
		log.Output(2, fmt.Sprint("[DEBUG] ", fmt.Sprint(args...)))
	}
}

func Debugf(format string, args ...interface{}) {
	if debugMode {
		log.Output(2, fmt.Sprintf("[DEBUG] "+format, args...))
	}
}

func Info(args ...interface{}) {
	log.Output(2, fmt.Sprint("[INFO] ", fmt.Sprint(args...)))
}

func Infof(format string, args ...interface{}) {
	log.Output(2, fmt.Sprintf("[INFO] "+format, args...))
}

func Warn(args ...interface{}) {
	log.Output(2, fmt.Sprint("[WARN] ", fmt.Sprint(args...)))
}

func Warnf(format string, args ...interface{}) {
	log.Output(2, fmt.Sprintf("[WARN] "+format, args...))
}

func Error(args ...interface{}) {
	log.Output(2, fmt.Sprint("[ERROR] ", fmt.Sprint(args...)))
}

func Errorf(format string, args ...interface{}) {
	log.Output(2, fmt.Sprintf("[ERROR] "+format, args...))
}

func Fatal(args ...interface{}) {
	log.Output(2, fmt.Sprint("[FATAL] ", fmt.Sprint(args...)))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	log.Output(2, fmt.Sprintf("[FATAL] "+format, args...))
	os.Exit(1)
}
