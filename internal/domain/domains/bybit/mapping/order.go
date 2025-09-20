package mapping

import (
	"fmt"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
)

const (
	bybitOrderTypeLimit  = "LIMIT"
	bybitOrderTypeMarket = "MARKET"
)

var bybitOrderTypes = map[string]string{
	trading.TypeLimit:  bybitOrderTypeLimit,
	trading.TypeMarket: bybitOrderTypeMarket,
}

func ToBybitOrderType(domainOrderType string) (bybitOrderType string, err error) {
	var ok bool
	bybitOrderType, ok = bybitOrderTypes[domainOrderType]
	if ok == false {
		msg := fmt.Sprintf("Unexpected value (%s) to map Bybit order type", domainOrderType)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}
