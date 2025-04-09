package tools

import (
	"errors"
	"fmt"
)

type AppError struct {
	Message     string
	ParentError error
	Code        float64
}

var EmptyValueError = errors.New("empty value")

func (t AppError) Error() string {
	if t.ParentError != nil && t.ParentError.Error() != "" {
		return fmt.Sprintf("%s. ParentError: %s", t.Message, t.ParentError.Error())
	}
	return t.Message
}
