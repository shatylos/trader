package apperrors

import (
	"fmt"
	"runtime"
)

func New(format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s\n    %s:%d\n", msg, file, line)
	return
}

func Wrap(parent error, format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s\n    %s:%d\n%w", msg, file, line, parent)
	return
}
