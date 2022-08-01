package _interface

type StrategyInterface interface {
	IsInit() bool
	Initialise() error
	GetData() error
	Analyse() error
	DoAction() error
	Wait()
}
