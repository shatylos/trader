package domainInterface

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	IsDemoMode() bool
	GetWallet() (*structs.DomainWallet, error)
}
