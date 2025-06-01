package domain

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures"
	"github.com/shatylos/trader/internal/domain/domains/binanceSpot"
	"github.com/shatylos/trader/internal/domain/domains/bybit"
	"github.com/shatylos/trader/tools"
)

const DomainBybitFutures = "bybitMargin"
const DomainBybitSpot = "bybitSpot"
const DomainBinanceFutures = "binanceFutures"
const DomainBinanceSpot = "binanceSpot"

func GetSpotDomain(domainCode string) (SpotDomainInterface, error) {

	switch domainCode {
	case DomainBybitSpot:
		return &bybit.DomainBybitSpot{}, nil
	case DomainBinanceSpot:
		return &binanceSpot.DomainBinanceSpot{}, nil
	}

	return nil, tools.AppError{
		Message: fmt.Sprintf("spot domain with code \"%s\" not implemented", domainCode),
	}
}

func GetFuturesDomain(domainCode string) (FuturesDomainInterface, error) {

	switch domainCode {
	case DomainBybitFutures:
		return &bybit.DomainBybitFutures{}, nil
	case DomainBinanceFutures:
		return &binanceFeatures.DomainBinanceFutures{}, nil
	}

	return nil, tools.AppError{
		Message: fmt.Sprintf("futures domain with code \"%s\" not implemented", domainCode),
	}
}
