package domain

import (
	"bitbucket.org/shatylos/trader/domain/domains/exmo"
	"bitbucket.org/shatylos/trader/domain/interface"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

var (
	domainList          = make([]string, 0)
	domainInterfaceList = make(map[string]_interface.DomainInterface, 0)
)

func init() {
	domainList = append(domainList, "exmo")
	domainInterfaceList["exmo"] = exmo.DomainExmo{}
}

func GetDomainList() []string {
	return domainList
}

func GetDomainInterface(domainCode string) (_interface.DomainInterface, error) {
	if domainInterfaceObj, ok := domainInterfaceList[domainCode]; ok {
		return domainInterfaceObj, nil
	}
	return nil, utils.AppError{
		Message: fmt.Sprintf("domain with code \"%s\" not implemented", domainCode),
	}
}
