package utils

import "os"

func AppConfig(params ...string) string {
	paramKey := params[0]
	defaultValue := ""

	if len(params) > 1 {
		defaultValue = params[1]
	}

	resultValue := os.Getenv(paramKey)

	if resultValue == "" {
		resultValue = defaultValue
	}

	return resultValue
}
