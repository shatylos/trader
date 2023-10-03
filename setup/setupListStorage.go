package setup

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"bitbucket.org/shatylos/trader/strategy/strategies/buyCheapSellHigh"
	tradingConstant "bitbucket.org/shatylos/trader/trading/constant"
	"time"
)

var setupList []*setupStructs.Setup

func init() {
	setupListInit()
}

func setupListInit() {
	setupList = make([]*setupStructs.Setup, 0)

	//setupList = append(setupList, &setupStructs.Setup{
	//	Strategy: &strategies.buyCheapSellHigh{
	//		DomainCode:          constant.DomainBybitMargin,
	//		CoinPare:            "BTCUSDT",
	//		Resolution:          tradingConstant.Resol30m,
	//		CandlesToAnalyse:    10,
	//		TimeoutSeconds:      5,
	//		CostDiffToStopTrade: 250,
	//		AvgCostShift:        50,
	//		Leverage:            1,
	//		Qty:                 0.005,
	//		QtyCoefficient:      2,
	//		TakeProfitSize:      50,
	//		StopLossSize:        200,
	//	},
	//})

	setupList = append(setupList, &setupStructs.Setup{
		Strategy: &buyCheapSellHigh.BuyCheapSellHigh{
			DomainCode:            constant.DomainBybitSpot,
			CoinPare:              "BTCUSDT",
			MainCurrency:          "USDT",
			TradeCurrency:         "BTC",
			TimeoutSeconds:        10,
			Resolution:            tradingConstant.Resol30m,
			CostRanges:            []int64{300, 600, 900, 1200},
			PercentRanges:         []int64{7, 8, 10, 12},
			LongTermMaxPrice:      70000,
			LongTermMinPrice:      10000,
			LongTermPercentBuffer: 10,
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
