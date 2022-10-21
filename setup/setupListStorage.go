package setup

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"bitbucket.org/shatylos/trader/strategy/strategies"
	tradingConstant "bitbucket.org/shatylos/trader/trading/constant"
	"time"
)

var setupList []*setupStructs.Setup

func init() {
	setupListInit()
}

func setupListInit() {
	setupList = make([]*setupStructs.Setup, 0)

	setupList = append(setupList, &setupStructs.Setup{
		Strategy: &strategies.ScalpByProbabilityStrategy{
			DomainCode:          constant.DomainBybitMargin,
			CoinPare:            "BTCUSDT",
			Resolution:          tradingConstant.Resol30m,
			CandlesToAnalyse:    10,
			TimeoutSeconds:      5,
			CostDiffToStopTrade: 250,
			AvgCostShift:        50,
			Leverage:            2,
			Qty:                 0.005,
			TakeProfitSize:      50,
			StopLossSize:        250,
		},
	})
}

func LoadNextSetupStep(setupChanelContext chan *setupStructs.Setup) {
	var setupItemResult *setupStructs.Setup

	for {
		for _, setupListItem := range setupList {
			if setupListItem.GetStatus() == setupStructs.StatusReadyForNext {
				setupListItem.SetStatus(setupStructs.StatusInProgress)
				setupItemResult = setupListItem
				break
			}
		}

		if setupItemResult != nil {
			setupChanelContext <- setupItemResult
			return
		}
		time.Sleep(time.Second / 4)
	}
}
