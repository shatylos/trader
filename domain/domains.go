package domain

import (
	"bitbucket.org/shatylos/trader/domain/domainInterface"
	"bitbucket.org/shatylos/trader/domain/domains/exmo"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

var (
	domainList          = make([]string, 0)
	domainInterfaceList = make(map[string]domainInterface.DomainInterface, 0)
)

func init() {
	domainList = append(domainList, "exmo")
	domainInterfaceList["exmo"] = exmo.DomainExmo{}
}

func GetDomainList() []string {
	return domainList
}

func GetDomainInterface(domainCode string) (domainInterface.DomainInterface, error) {
	if domainInterfaceObj, ok := domainInterfaceList[domainCode]; ok {
		return domainInterfaceObj, nil
	}
	return nil, utils.AppError{
		Message: fmt.Sprintf("domain with code \"%s\" not implemented", domainCode),
	}
}
