package domain

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit"
	"bitbucket.org/shatylos/trader/domain/domains/exmo"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

var (
	domainNameList = make([]string, 0)
)

func init() {
	domainNameList = append(domainNameList, constant.DomainExmo)
	domainNameList = append(domainNameList, constant.DomainExmoMargin)
	domainNameList = append(domainNameList, constant.DomainBybitMargin)
	domainNameList = append(domainNameList, constant.DomainBybitSpot)
}

func GetDomainNameList() []string {
	return domainNameList
}

func GetDomainInterface(domainCode string) (DomainInterface, error) {

	switch domainCode {
	case constant.DomainExmo:
		return &exmo.DomainExmo{}, nil
	case constant.DomainExmoMargin:
		return &exmo.DomainExmoMargin{}, nil
	case constant.DomainBybitMargin:
		return &bybit.DomainBybitMargin{}, nil
	case constant.DomainBybitSpot:
		return &bybit.DomainBybitSpot{}, nil
	}

	return nil, utils.AppError{
		Message: fmt.Sprintf("domain with code \"%s\" not implemented", domainCode),
	}
}
