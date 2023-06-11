package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
)

func GetOpenOrderList(domainCode string, coinPare string) ([]structs.DomainOrder, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}

	return domainInterface.GetOpenOrderList(coinPare)
}

func OpenOrder(domainCode string, orderRequest structs.DomainOrderRequest) (string, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return "", err
	}
	return domainInterface.OpenOrder(orderRequest)
}

func CancelOrder(domainCode string, orderId string) error {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return err
	}
	return domainInterface.CancelOrder(orderId)
}

func GetHistoryOrders(domainCode string, limit int64) ([]structs.DomainOrder, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}
	return domainInterface.GetHistoryOrders(limit)
}
