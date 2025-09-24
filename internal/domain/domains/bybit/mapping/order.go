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

const (
	bybitSideBuy  = "BUY"
	bybitSideSell = "SELL"
)

var bybitOrderTypes = map[string]string{
	trading.TypeLimit:  bybitOrderTypeLimit,
	trading.TypeMarket: bybitOrderTypeMarket,
}

var bybitOrderSides = map[string]string{
	trading.SideBuy:  bybitSideBuy,
	trading.SideSell: bybitSideSell,
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

func ToDomainOrderSide(bybitSide string) (domainSide string, err error) {
	for sideD, sideB := range bybitOrderSides {
		if sideB == bybitSide {
			domainSide = sideD
		}
	}
	if domainSide == "" {
		msg := fmt.Sprintf("Bybit side %s can not be mapped to domain side", bybitSide)
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}
	return
}

func ToBybitOrderSide(domainSide string) (bybitSide string, err error) {
	var ok bool
	bybitSide, ok = bybitOrderSides[domainSide]
	if ok == false {
		msg := fmt.Sprintf("Unexpected value (%s) to map Bybit order side", domainSide)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}
