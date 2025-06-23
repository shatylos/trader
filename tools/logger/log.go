package logger

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGrey   = "\033[37m"
)

var isInit bool

func logInit() {
	log.SetOutput(os.Stdout)
}

func Error(msg string) {
	if !isInit {
		logInit()
	}
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	stackTrace := string(debug.Stack())
	log.Printf("%sError%s %s [%s:%d %s]\nStack Trace:\n%s",
		colorRed, colorReset, msg, file, line, funcName, stackTrace)
}

func Warning(msg string) {
	if !isInit {
		logInit()
	}
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	//stackTrace := string(debug.Stack())
	//log.Printf("%sWarning%s %s [%s:%d %s]\nStack Trace:\n%s",
	//	colorYellow, colorReset, msg, file, line, funcName, stackTrace)
	log.Printf("%sWarning%s %s [%s:%d %s]\n",
		colorYellow, colorReset, msg, file, line, funcName)
}

func Info(msg string) {
	if !isInit {
		logInit()
	}
	log.Println(colorGrey, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}

func Success(msg string) {
	if !isInit {
		logInit()
	}
	log.Println(colorGreen, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}
