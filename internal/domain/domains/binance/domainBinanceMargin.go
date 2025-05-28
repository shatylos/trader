package binance

import (
	"github.com/shatylos/trader/internal/domain/structs"
)

type DomainBinanceMargin struct {
	code string
	//secrets Secrets
}

func (d *DomainBinanceMargin) GetCode() string {
	return d.code
}

func (d *DomainBinanceMargin) SetConfig(config map[interface{}]interface{}) (err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) GetWallet() (wallet *structs.DomainWallet, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) (candles []structs.DomainCandle, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) GetOrder(domainId string) (order structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) GetOpenOrderList(coinPare string) (orderList []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) GetPosition(coinPare string) (structs.DomainPosition, error) {
	panic("Not implemented")
}

func (d *DomainBinanceMargin) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {
	panic("Not implemented")
}
func (d *DomainBinanceMargin) ModifyTpSl(tpSlRequest structs.TpSlRequest) error {
	panic("Not implemented")
}

func (d *DomainBinanceMargin) OpenOrder(orderRequest structs.DomainOrderRequest) (orderID string, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) CancelOrder(orderId string, coinPare string) (err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceMargin) GetHistoryOrders(limit int64, coinPare string) (domainOrders []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}
