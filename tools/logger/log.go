package logger

import (
	"log"
	"os"
	"runtime"
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

func Error(msg string) {
	log.SetOutput(os.Stdout)
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	log.Println(colorRed, "Error", colorReset, msg, file, line, funcName)
}

func Warning(msg string) {
	log.SetOutput(os.Stdout)
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	log.Println(colorYellow, "Warning", colorReset, msg, file, line, funcName)
}

func Info(msg string) {
	log.Println(colorGrey, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}

func Success(msg string) {
	log.Println(colorGreen, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}
