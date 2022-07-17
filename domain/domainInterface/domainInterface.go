package domainInterface

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	GetWallet() (*structs.DomainWallet, error)
}
