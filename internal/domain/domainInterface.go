package domain

import (
	"github.com/shatylos/trader/internal/domain/structs"
)

type SpotDomainInterface interface {
	GetOrder(domainId string) (structs.DomainOrder, error)
	GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error)
	GetWallet() (*structs.DomainWallet, error)
	LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error)
	OpenOrder(orderRequest structs.DomainOrderRequest) (string, error)
	CancelOrder(orderId string, coinPare string) error
	GetHistoryOrders(limit int64, coinPare string) ([]structs.DomainOrder, error)
	SetConfig(map[interface{}]interface{}) error
	GetCode() string
}

type FuturesDomainInterface interface {
	GetOrder(domainId string) (structs.DomainOrder, error)
	GetPosition(coinPare string) (structs.DomainPosition, error)
	GetWallet() (*structs.DomainWallet, error)
	LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error)
	OpenPosition(positionRequest structs.DomainPositionRequest) (string, error)
	ModifyTpSl(tpSlRequest structs.TpSlRequest) error
	SetConfig(map[interface{}]interface{}) error
	GetCode() string
}
