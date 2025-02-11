package tools

import "errors"

type AppError struct {
	Message     string
	ParentError error
	Code        float64
}

var EmptyValueError = errors.New("empty value")

func (t AppError) Error() string {
	return t.Message
}
