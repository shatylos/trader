package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
	"sort"
)

func GetCandleHistory(domainCode string, symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}

	candles, err := domainInterface.LoadCandleHistory(symbol, resolution, from, limit)
	if err != nil {
		return nil, err
	}

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Time > candles[j].Time
	})

	return candles, nil
}
