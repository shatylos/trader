package structs

type SetupInterface interface {
	GetStatus() int64
	SetStatus(status int64)
	NextStep()
}
