package _interface

type StrategyInterface interface {
	GetData() error
	Analyse() error
	DoAction() error
	Wait()
}
