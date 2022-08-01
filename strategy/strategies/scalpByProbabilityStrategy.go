package strategies

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/utils"
	"time"
)

type ScalpByProbabilityStrategy struct {
	isInit         bool
	DomainCode     string
	CoinPare       string
	TimeoutSeconds time.Duration
}

func (s *ScalpByProbabilityStrategy) IsInit() bool {
	return s.isInit
}

func (s *ScalpByProbabilityStrategy) Initialise() error {
	domainItem, err := domain.GetDomainInterface(s.DomainCode)
	if err != nil {
		return err
	}
	if domainItem.GetType() != constant.DomainTypeMargin {
		return utils.AppError{
			Message: "Strategy ScalpByProbabilityStrategy works only with margin domain type",
		}
	}
	s.isInit = true
	return nil
}

func (s *ScalpByProbabilityStrategy) GetData() error {
	// get last N costs by timeframe T
	// get open orders for setup
	return nil
}

func (s *ScalpByProbabilityStrategy) Analyse() error {
	// check open orders. Return if orders exists
	// detect enter point, calculate TP/SL
	return nil
}

func (s *ScalpByProbabilityStrategy) DoAction() error {
	// open orders
	return nil
}

func (s *ScalpByProbabilityStrategy) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}
