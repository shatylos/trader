package domain

import (
	"fmt"
	"github.com/shatylos/trader/domain/constant"
	"github.com/shatylos/trader/domain/domains/bybit"
	"github.com/shatylos/trader/domain/domains/exmo"
	"github.com/shatylos/trader/utils"
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
