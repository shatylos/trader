package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
)

func GetCandleHistory(domainCode string, symbol string, resolution string, from int64, to int64) ([]structs.DomainCandle, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}

	return domainInterface.LoadCandleHistory(symbol, resolution, from, to)
}
