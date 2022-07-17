package services

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
)

func LoadWalletInfo(domainCode string, domainWalletResultChan chan structs.DomainWalletApiResult, bufferParReq chan bool) {
	domainWallet, err := loadWalletInfo(domainCode)

	if err == nil {
		domainWalletResultChan <- structs.DomainWalletApiResult{
			DomainCode:   domainCode,
			DomainWallet: domainWallet,
			Error:        nil,
		}
	} else {
		domainWalletResultChan <- structs.DomainWalletApiResult{
			DomainCode:   domainCode,
			DomainWallet: nil,
			Error:        err,
		}
	}

	<-bufferParReq
}

func loadWalletInfo(domainCode string) (*structs.DomainWallet, error) {
	domainInterface, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return nil, err
	}
	return domainInterface.GetWallet()
}
