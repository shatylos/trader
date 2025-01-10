package setup

import "github.com/shatylos/trader/setup/structs"

type ReaderInterface interface {
	GetConfig() (*structs.Config, error)
}
