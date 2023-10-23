package structs

type Config struct {
	StrategyItems []interface{}               `yaml:"strategy"`
	DomainItems   map[interface{}]interface{} `yaml:"domain"`
}
