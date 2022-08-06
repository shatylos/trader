package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
	"bitbucket.org/shatylos/trader/utils"
)

func GetPositionList(domainCode string, setupId string) ([]structs.DomainPosition, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}

	positionResult := make([]structs.DomainPosition, 0)

	positionIds := []int64{692852959514642257, 692852959514642301}

	allPositions, err := domainInterface.GetPositionList()

	for _, position := range allPositions {
		if utils.ContainsInt(positionIds, position.PositionId) {
			positionResult = append(positionResult, position)
		}
	}

	return positionResult, err
}
