package utils

type AppError struct {
	Message     string
	ParentError error
}

func (t AppError) Error() string {
	return t.Message
}
