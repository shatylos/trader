package _interface

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	GetType() int64
	IsDemoMode() bool
	GetWallet() (*structs.DomainWallet, error)
}
