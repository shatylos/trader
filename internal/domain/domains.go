package domain

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binance"
	"github.com/shatylos/trader/internal/domain/domains/bybit"
	"github.com/shatylos/trader/tools"
)

const DomainBybitMargin = "bybitMargin"
const DomainBybitSpot = "bybitSpot"
const DomainBinanceMargin = "binanceMargin"
const DomainBinanceSpot = "binanceSpot"

func GetSpotDomain(domainCode string) (SpotDomainInterface, error) {

	switch domainCode {
	case DomainBybitSpot:
		return &bybit.DomainBybitSpot{}, nil
	case DomainBinanceSpot:
		return &binance.DomainBinanceSpot{}, nil
	}

	return nil, tools.AppError{
		Message: fmt.Sprintf("spot domain with code \"%s\" not implemented", domainCode),
	}
}

func GetMarginDomain(domainCode string) (MarginDomainInterface, error) {

	switch domainCode {
	case DomainBybitMargin:
		return &bybit.DomainBybitMargin{}, nil
	case DomainBinanceMargin:
		return &binance.DomainBinanceMargin{}, nil
	}

	return nil, tools.AppError{
		Message: fmt.Sprintf("margin domain with code \"%s\" not implemented", domainCode),
	}
}
