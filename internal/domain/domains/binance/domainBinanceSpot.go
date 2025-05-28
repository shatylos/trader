package binance

import (
	"github.com/shatylos/trader/internal/domain/structs"
)

type DomainBinanceSpot struct {
	code string
	//secrets Secrets
}

func (d *DomainBinanceSpot) GetCode() string {
	return d.code
}

func (d *DomainBinanceSpot) SetConfig(config map[interface{}]interface{}) (err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) GetWallet() (wallet *structs.DomainWallet, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) (candles []structs.DomainCandle, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) GetOrder(domainId string) (order structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) GetOpenOrderList(coinPare string) (orderList []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) OpenOrder(orderRequest structs.DomainOrderRequest) (orderID string, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) CancelOrder(orderId string, coinPare string) (err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceSpot) GetHistoryOrders(limit int64, coinPare string) (domainOrders []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}
