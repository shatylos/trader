package _interface

type StrategyInterface interface {
	IsInit() bool
	Initialise() error
	DoAction() error
	Wait()
}
