package domain

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit"
	"bitbucket.org/shatylos/trader/domain/domains/exmo"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

var (
	domainList          = make([]string, 0)
	domainInterfaceList = make(map[string]DomainInterface, 0)
)

func init() {
	domainList = append(domainList, constant.DomainExmo)
	domainList = append(domainList, constant.DomainExmoMargin)
	domainList = append(domainList, constant.DomainBybitMargin)
	domainList = append(domainList, constant.DomainBybitSpot)
	domainList = append(domainList, constant.DomainBybitSpotDemo)
	domainList = append(domainList, constant.DomainBybitMarginDemo)

	domainInterfaceList[constant.DomainExmo] = &exmo.DomainExmo{}
	domainInterfaceList[constant.DomainExmoMargin] = &exmo.DomainExmoMargin{}
	domainInterfaceList[constant.DomainBybitMargin] = &bybit.DomainBybitMargin{IsDemo: false}
	domainInterfaceList[constant.DomainBybitSpot] = &bybit.DomainBybitSpot{IsDemo: false}
	domainInterfaceList[constant.DomainBybitSpotDemo] = &bybit.DomainBybitSpot{IsDemo: true}
	domainInterfaceList[constant.DomainBybitMarginDemo] = &bybit.DomainBybitMargin{IsDemo: true}
}

func GetDomainList() []string {
	return domainList
}

func GetDomainInterface(domainCode string) (DomainInterface, error) {
	if domainInterfaceObj, ok := domainInterfaceList[domainCode]; ok {
		return domainInterfaceObj, nil
	}
	return nil, utils.AppError{
		Message: fmt.Sprintf("domain with code \"%s\" not implemented", domainCode),
	}
}
