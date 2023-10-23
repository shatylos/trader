package strategy

type StrategyInterface interface {
	SetConfig(interface{}, map[interface{}]interface{}) error
	IsInit() bool
	Initialise() error
	DoAction() error
	Wait()
}
