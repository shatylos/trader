package setup

import "bitbucket.org/shatylos/trader/setup/structs"

type ReaderInterface interface {
	GetConfig() (*structs.Config, error)
}
