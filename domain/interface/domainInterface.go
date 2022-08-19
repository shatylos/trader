package _interface

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	GetType() int64
	IsDemoMode() bool
	GetWallet() (*structs.DomainWallet, error)
	LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error)
	GetPositionList(coinPare string) ([]structs.DomainPosition, error)
	OpenPosition(positionRequest structs.DomainPositionRequest) (string, error)
}
