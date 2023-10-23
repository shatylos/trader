package domain

import "bitbucket.org/shatylos/trader/domain/structs"

type DomainInterface interface {
	GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error)
	GetPositionList(coinPare string) ([]structs.DomainPosition, error)
	GetType() int64
	GetWallet() (*structs.DomainWallet, error)
	LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error)
	OpenPosition(positionRequest structs.DomainPositionRequest) (string, error)
	OpenOrder(orderRequest structs.DomainOrderRequest) (string, error)
	CancelOrder(orderId string) error
	GetHistoryOrders(limit int64) ([]structs.DomainOrder, error)
	SetConfig(map[interface{}]interface{}) error
	GetCode() string
}
