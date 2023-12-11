package utils

import "errors"

type AppError struct {
	Message     string
	ParentError error
}

var EmptyValueError = errors.New("empty value")

func (t AppError) Error() string {
	return t.Message
}
