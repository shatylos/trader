package _interface

type SetupInterface interface {
	GetStatus() int64
	SetStatus(status int64)
	NextStep()
}
