package structs

type StrategyInterface interface {
	IsInit() bool
	Initialise() error
	DoAction() error
	Wait()
}
