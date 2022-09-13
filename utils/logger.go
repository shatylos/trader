package utils

import "time"

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

func LogError(msg string) {
	println(colorRed, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}

func LogInfo(msg string) {
	println(colorGrey, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}

func LogSuccess(msg string) {
	println(colorGreen, time.Now().Format("2006-01-02 15:04:05"), msg, colorReset)
}
