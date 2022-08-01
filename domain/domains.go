package domain

import (
	"bitbucket.org/shatylos/trader/domain/constant"
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
	domainList = append(domainList, constant.DomainExmo)
	domainList = append(domainList, constant.DomainExmoMargin)

	domainInterfaceList[constant.DomainExmo] = &exmo.DomainExmo{}
	domainInterfaceList[constant.DomainExmoMargin] = &exmo.DomainExmoMargin{}
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
