package mapping

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
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
	structs.OrderTypes.Limit:  bybitOrderTypeLimit,
	structs.OrderTypes.Market: bybitOrderTypeMarket,
}

var bybitOrderSides = map[string]string{
	structs.OrderSideBuy:  bybitSideBuy,
	structs.OrderSideSell: bybitSideSell,
}

func ToBybitOrderType(domainOrderType string) (bybitOrderType string, err error) {
	var ok bool
	bybitOrderType, ok = bybitOrderTypes[domainOrderType]
	if ok == false {
		err = apperrors.New("unexpected value (%s) to map Bybit order type", domainOrderType)
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
		err = apperrors.New("Bybit side %s can not be mapped to domain side", bybitSide)
		return
	}
	return
}

func ToBybitOrderSide(domainSide string) (bybitSide string, err error) {
	var ok bool
	bybitSide, ok = bybitOrderSides[domainSide]
	if ok == false {
		err = apperrors.New("unexpected value (%s) to map Bybit order side", domainSide)
		return
	}
	return
}
