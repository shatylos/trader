package strategies

import (
	"time"
)

type ScalpByProbabilityStrategy struct {
	DomainCode  string
	CoinPare    string
	WaitSeconds time.Duration
}

func (s *ScalpByProbabilityStrategy) GetData() error {
	return nil
}

func (s *ScalpByProbabilityStrategy) Analyse() error {
	return nil
}

func (s *ScalpByProbabilityStrategy) DoAction() error {
	return nil
}

func (s *ScalpByProbabilityStrategy) Wait() {
	time.Sleep(time.Second * s.WaitSeconds)
}
