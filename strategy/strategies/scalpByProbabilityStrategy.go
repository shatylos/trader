package strategies

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"time"
)

type ScalpByProbabilityStrategy struct {
	AvgCostShift        float64
	CandlesToAnalyse    int64
	CoinPare            string
	CostDiffToStopTrade float64
	DomainCode          string
	isInit              bool
	Leverage            int64
	Qty                 float64
	Resolution          string
	StopLossSize        float64
	TakeProfitSize      float64
	TimeoutSeconds      time.Duration
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

	avgCost, err := s.getAverageCost()
	if err != nil {
		return err
	}
	currentCost, err := s.getCurrentCost()
	if err != nil {
		return err
	}
	costDiff := s.getCostDiff()

	avgCostShift := s.AvgCostShift

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

	if doOpenBuy {
		if currentCost < avgCost-avgCostShift && costDiff < s.CostDiffToStopTrade {
			tp := currentCost + s.TakeProfitSize
			sl := currentCost - s.StopLossSize
			utils.LogInfo(fmt.Sprintf("Add position to buy. Current cost: %f, TP: %f, SL: %f", currentCost, tp, sl))
			positionToAdd := structs.DomainPositionRequest{
				Leverage:   s.Leverage, //1,         //         int64
				Qty:        s.Qty,      //0.005,     //              float64
				Side:       "Buy",      //             string //@TODO move the value to a constant
				TakeProfit: tp,         //23600,     //       float64
				StopLoss:   sl,         //22800,     //         float64
				Symbol:     s.CoinPare, //"BTCUSDT", //           string
				Type:       "Limit",    //@TODO move the value to a constant
				Price:      currentCost,
			}
			positionsToOpen = append(positionsToOpen, positionToAdd)
		}
	}

	if doOpenSell {
		if currentCost > avgCost+avgCostShift && costDiff < s.CostDiffToStopTrade {
			tp := currentCost - s.TakeProfitSize
			sl := currentCost + s.StopLossSize
			utils.LogInfo(fmt.Sprintf("Add position to sell. Current cost: %f, TP: %f, SL: %f", currentCost, tp, sl))
			positionToAdd := structs.DomainPositionRequest{
				Leverage:   s.Leverage, //1,         //         int64
				Qty:        s.Qty,      //0.005,     //              float64
				Side:       "Sell",     //             string //@TODO move the value to a constant
				TakeProfit: tp,         //23600,     //       float64
				StopLoss:   sl,         //22800,     //         float64
				Symbol:     s.CoinPare, //"BTCUSDT", //           string
				Type:       "Limit",    //@TODO move the value to a constant
				Price:      currentCost,
			}
			positionsToOpen = append(positionsToOpen, positionToAdd)
		}
	}

	return nil
}

func (s *ScalpByProbabilityStrategy) DoAction() error {

	for _, position := range positionsToOpen {
		positionId, err := services.OpenPosition(s.DomainCode, position)
		if err != nil {
			return err
		}
		utils.LogSuccess(fmt.Sprintf("Position id: %s", positionId))
	}
	positionsToOpen = make([]structs.DomainPositionRequest, 0)
	return nil
}

func (s *ScalpByProbabilityStrategy) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}

func (s *ScalpByProbabilityStrategy) getAverageCost() (float64, error) {
	allAvgs := 0.0
	for _, candle := range candles {
		maxCost := candle.High
		minCost := candle.Low
		allAvgs += (maxCost + minCost) / 2
	}
	if allAvgs == 0.0 {
		return 0, utils.AppError{
			Message: "Average cost can not be calculated",
		}
	}
	avg := allAvgs / float64(len(candles))
	return avg, nil
}

func (s *ScalpByProbabilityStrategy) getCurrentCost() (float64, error) {
	if len(candles) == 0 {
		return 0, utils.AppError{
			Message: "Current cost not found",
		}
	}
	return candles[0].Close, nil
}

func (s *ScalpByProbabilityStrategy) getCostDiff() float64 {
	minCost := 0.0
	maxCost := 0.0

	for _, candle := range candles {
		if candle.High > maxCost {
			maxCost = candle.High
		}
		if candle.Low < minCost || minCost == 0.0 {
			minCost = candle.Low
		}
	}

	return maxCost - minCost
}
