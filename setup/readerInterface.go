package setup

import setupStructs "bitbucket.org/shatylos/trader/setup/structs"

type ReaderInterface interface {
	GetSetupList() ([]*setupStructs.Setup, error)
}
