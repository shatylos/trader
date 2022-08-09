package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
)

func GetPositionList(domainCode string, coinPare string) ([]structs.DomainPosition, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}

	return domainInterface.GetPositionList(coinPare)
}
