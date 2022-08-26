package strategies

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	"time"
)

type ScalpByProbabilityStrategy struct {
	isInit           bool
	DomainCode       string
	CoinPare         string
	Resolution       string
	CandlesToAnalyse int64
	TimeoutSeconds   time.Duration
}

var candles []structs.DomainCandle
var positions []structs.DomainPosition
var positionsToOpen []structs.DomainPositionRequest

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
	_positions, err := services.GetPositionList(s.DomainCode, s.CoinPare)
	if err != nil {
		return err
	}
	positions = _positions
	if len(positions) >= 2 {
		return nil
	}

	now := time.Now()
	to := now.Unix()
	from := to - (tradeConst.ResolToSec[s.Resolution] * s.CandlesToAnalyse)

	_candles, err := services.GetCandleHistory(s.DomainCode, s.CoinPare, s.Resolution, from, s.CandlesToAnalyse)
	if err != nil {
		return err
	}
	candles = _candles
	return nil
}

func (s *ScalpByProbabilityStrategy) Analyse() error {

	doOpenBuy := true
	doOpenSell := true

	for _, position := range positions {
		if position.Side == tradeConst.SideBuy {
			doOpenBuy = false
		}
		if position.Side == tradeConst.SideSell {
			doOpenSell = false
		}
	}

	if !doOpenBuy && !doOpenSell {
		return nil
	}

	//avgCost := getAverageCost()
	//currentCost := getCurrentCost()
	//avgCostShift := conf.getAvgCostShift()

	//if doOpenBuy {
	//	if currentCost < avgCost-avgCostShift {
	//		println("Add position to buy")
	//	}
	//}
	//
	//if doOpenSell {
	//	if currentCost > avgCost+avgCostShift {
	//		println("Add position to sell")
	//	}
	//}

	//positionToAdd := structs.DomainPositionRequest{
	//	Leverage: 1, //         int64
	//	//PositionId: "1234qwe", //       string
	//	Qty:        0.005,     //              float64
	//	Side:       "Buy",     //             string
	//	TakeProfit: 23600,     //       float64
	//	StopLoss:   22800,     //         float64
	//	Symbol:     "BTCUSDT", //           string
	//	Type:       "Market",
	//}
	//positionsToOpen = append(positionsToOpen, positionToAdd)

	// detect enter point, calculate TP/SL
	return nil
}

func (s *ScalpByProbabilityStrategy) DoAction() error {

	for _, position := range positionsToOpen {
		positionId, err := services.OpenPosition(s.DomainCode, position)
		if err != nil {
			return err
		}
		println("Position id:", positionId)
	}
	positionsToOpen = make([]structs.DomainPositionRequest, 0)
	return nil
}

func (s *ScalpByProbabilityStrategy) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}
