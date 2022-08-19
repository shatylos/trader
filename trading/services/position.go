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

func OpenPosition(domainCode string, positionRequest structs.DomainPositionRequest) (string, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return "", err
	}
	return domainInterface.OpenPosition(positionRequest)
}
