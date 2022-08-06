package _interface

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	GetType() int64
	IsDemoMode() bool
	GetWallet() (*structs.DomainWallet, error)
	LoadCandleHistory(symbol string, resolution string, from int64, to int64) ([]structs.DomainCandle, error)
	GetPositionList() ([]structs.DomainPosition, error)
}
