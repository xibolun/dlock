package dlock

import (
	"log"
	"os"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type DLog struct {
	InfoL    *log.Logger
	WarningL *log.Logger
	ErrorL   *log.Logger
	TraceL   *log.Logger
	DebugL   *log.Logger
}

var dlog *DLog

func init() {
	Info := log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stdout, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	Debug := log.New(os.Stdout, "Debug:", log.Ldate|log.Ltime|log.Lshortfile)
	Trace := log.New(os.Stdout, "Trace:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning := log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)

	dlog = &DLog{InfoL: Info, WarningL: Warning, ErrorL: Error, TraceL: Trace, DebugL: Debug,}
}

func Debug(v ...interface{}) {
	dlog.DebugL.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	dlog.DebugL.Printf(format, v...)
}

func Info(v ...interface{}) {
	dlog.InfoL.Println(v...)
}

func Infof(format string, v ...interface{}) {
	dlog.InfoL.Printf(format, v...)
}

func Warning(v ...interface{}) {
	dlog.WarningL.Println(v...)
}

func Warningf(format string, v ...interface{}) {
	dlog.WarningL.Printf(format, v...)
}

func Trace(v ...interface{}) {
	dlog.TraceL.Println(v...)
}

func Tracef(format string, v ...interface{}) {
	dlog.TraceL.Printf(format, v...)
}

func Error(v ...interface{}) {
	dlog.ErrorL.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	dlog.ErrorL.Printf(format, v...)
}
