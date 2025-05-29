package binance

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	_type "github.com/shatylos/trader/tools/type"
)

type DomainBinanceFutures struct {
	code string
	//secrets Secrets
}

func (d *DomainBinanceFutures) GetCode() string {
	return d.code
}

func (d *DomainBinanceFutures) SetConfig(config map[interface{}]interface{}) (err error) {
	domainCode, err := _type.ToString(config["code"])
	if err != nil {
		return tools.AppError{
			Message: "The field code is empty or contains not correct value type. In DomainBinanceFutures config. Expects a string",
		}
	}
	d.code = domainCode

	return
}

func (d *DomainBinanceFutures) GetWallet() (wallet *structs.DomainWallet, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceFutures) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) (candles []structs.DomainCandle, err error) {

	panic("Not implemented")

	return
}

func (d *DomainBinanceFutures) GetOrder(domainId string) (order structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceFutures) GetOpenOrderList(coinPare string) (orderList []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceFutures) GetPosition(coinPare string) (structs.DomainPosition, error) {
	panic("Not implemented")
}

func (d *DomainBinanceFutures) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {
	panic("Not implemented")
}
func (d *DomainBinanceFutures) ModifyTpSl(tpSlRequest structs.TpSlRequest) error {
	panic("Not implemented")
}

func (d *DomainBinanceFutures) OpenOrder(orderRequest structs.DomainOrderRequest) (orderID string, err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceFutures) CancelOrder(orderId string, coinPare string) (err error) {
	panic("Not implemented")
	return
}

func (d *DomainBinanceFutures) GetHistoryOrders(limit int64, coinPare string) (domainOrders []structs.DomainOrder, err error) {
	panic("Not implemented")
	return
}
