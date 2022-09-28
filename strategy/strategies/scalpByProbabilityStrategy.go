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
var orders []structs.DomainOrder
var positions []structs.DomainPosition
var ordersToOpen []structs.DomainOrderRequest
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

	_orders, err := services.GetOpenOrderList(s.DomainCode, s.CoinPare)
	if err != nil {
		return err
	}
	orders = _orders

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

	buyPositionExists := false
	sellPositionExists := false
	buyOrderTPExists := false
	sellOrderTPExists := false
	buyPositionPrice := 0.0
	sellPositionPrice := 0.0

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
			buyPositionExists = true
			buyPositionPrice = position.Price
		}
		if position.Side == tradeConst.SideSell {
			sellPositionExists = true
			sellPositionPrice = position.Price
		}
	}

	for _, order := range orders {
		if order.Side == tradeConst.SideBuy && order.OrderType == "Limit" && order.ReduceOnly {
			buyOrderTPExists = true
		}
		if order.Side == tradeConst.SideSell && order.OrderType == "Limit" && order.ReduceOnly {
			sellOrderTPExists = true
		}
	}

	if buyPositionExists {
		if !sellOrderTPExists {
			// do open sell TP order
			tpPrice := buyPositionPrice + s.TakeProfitSize
			utils.LogInfo(fmt.Sprintf("Add limit order to sell position (as TP). Position cost: %f, TP: %f", buyPositionPrice, tpPrice))
			orderToOpen := structs.DomainOrderRequest{
				Price:       tpPrice,
				Qty:         s.Qty, //
				ReduceOnly:  true,
				Side:        "Sell", // @TODO move the value to a constant
				Symbol:      s.CoinPare,
				TimeInForce: "GoodTillCancel",
				Type:        "Limit", // @TODO move the value to a constant
			}
			ordersToOpen = append(ordersToOpen, orderToOpen)
		}
	} else {
		if currentCost < avgCost-avgCostShift && costDiff < s.CostDiffToStopTrade {
			sl := currentCost - s.StopLossSize
			utils.LogInfo(fmt.Sprintf("Add position to buy. Current cost: %f, SL: %f", currentCost, sl))
			positionToAdd := structs.DomainPositionRequest{
				Leverage:    s.Leverage, //
				Price:       currentCost,
				Qty:         s.Qty, //
				ReduceOnly:  false,
				Side:        "Buy", // @TODO move the value to a constant
				StopLoss:    sl,
				Symbol:      s.CoinPare,
				TimeInForce: "FillOrKill",
				Type:        "Limit", // @TODO move the value to a constant
			}
			positionsToOpen = append(positionsToOpen, positionToAdd)
		}
	}

	if sellPositionExists {
		if !buyOrderTPExists {
			// do open but TP order
			tpPrice := sellPositionPrice - s.TakeProfitSize
			utils.LogInfo(fmt.Sprintf("Add limit order to buy position (as TP). Position cost: %f, TP: %f", sellPositionPrice, tpPrice))
			orderToOpen := structs.DomainOrderRequest{
				Price:       tpPrice,
				Qty:         s.Qty, //
				ReduceOnly:  true,
				Side:        "Buy", // @TODO move the value to a constant
				Symbol:      s.CoinPare,
				TimeInForce: "GoodTillCancel",
				Type:        "Limit", // @TODO move the value to a constant
			}
			ordersToOpen = append(ordersToOpen, orderToOpen)
		}
	} else {
		if currentCost > avgCost+avgCostShift && costDiff < s.CostDiffToStopTrade {
			tp := currentCost - s.TakeProfitSize
			sl := currentCost + s.StopLossSize
			utils.LogInfo(fmt.Sprintf("Add position to sell. Current cost: %f, TP: %f, SL: %f", currentCost, tp, sl))
			positionToAdd := structs.DomainPositionRequest{
				Leverage:    s.Leverage,
				Price:       currentCost,
				Qty:         s.Qty,
				ReduceOnly:  false,
				Side:        "Sell", //@TODO move the value to a constant
				StopLoss:    sl,
				Symbol:      s.CoinPare,
				TimeInForce: "FillOrKill",
				Type:        "Limit", //@TODO move the value to a constant
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

	for _, order := range ordersToOpen {
		orderId, err := services.OpenOrder(s.DomainCode, order)
		if err != nil {
			return err
		}
		utils.LogSuccess(fmt.Sprintf("Order id: %s", orderId))
	}
	ordersToOpen = make([]structs.DomainOrderRequest, 0)

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
