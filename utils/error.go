package utils

type AppError struct {
	Message string
}

func (t AppError) Error() string {
	return t.Message
}
