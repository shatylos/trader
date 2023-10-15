package scalpByProbability

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"math"
	"time"
)

type ScalpByProbability struct {
	AvgCostShift        float64
	CandlesToAnalyse    int64
	CoinPare            string
	CostDiffToStopTrade float64
	DomainCode          string
	isInit              bool
	Leverage            int64
	Qty                 float64
	QtyCoefficient      float64
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

func (s *ScalpByProbability) IsInit() bool {
	return s.isInit
}

func (s *ScalpByProbability) Initialise() error {
	domainItem, err := domain.GetDomainInterface(s.DomainCode)
	if err != nil {
		return err
	}
	if domainItem.GetType() != constant.DomainTypeMargin {
		return utils.AppError{
			Message: "Strategy ScalpByProbability works only with margin domain type",
		}
	}
	s.isInit = true
	return nil
}

func (s *ScalpByProbability) GetData() error {
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

func (s *ScalpByProbability) Analyse() error {

	positionsToOpen = make([]structs.DomainPositionRequest, 0)
	ordersToOpen = make([]structs.DomainOrderRequest, 0)

	buyPositionExists := false
	sellPositionExists := false
	buyOrderTPExists := false
	sellOrderTPExists := false
	buyPositionPrice := 0.0
	buyPositionQty := 0.0
	sellPositionPrice := 0.0
	sellPositionQty := 0.0

	avgCost, err := s.getAverageCost()
	if err != nil {
		return err
	}
	currentCost, err := s.getCurrentCost()
	if err != nil {
		return err
	}
	costDiff := s.getCandleCostDiff()

	avgCostShift := s.AvgCostShift

	for _, position := range positions {
		if position.Side == tradeConst.SideBuy {
			buyPositionExists = true
			buyPositionPrice = position.Price
			buyPositionQty = position.Qty
		}
		if position.Side == tradeConst.SideSell {
			sellPositionExists = true
			sellPositionPrice = position.Price
			sellPositionQty = position.Qty
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
			tpPrice := math.Round((buyPositionPrice+s.TakeProfitSize)*100) / 100
			utils.LogInfo(fmt.Sprintf("Add limit order to sell position (as TP). Position cost: %f, TP: %f, QTY: %f", buyPositionPrice, tpPrice, buyPositionQty))
			orderToOpen := structs.DomainOrderRequest{
				Price:       tpPrice,
				Qty:         buyPositionQty,
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
			qty, err := s.getQtyByWallet(currentCost)
			if err != nil {
				return err
			}

			positionToAdd := structs.DomainPositionRequest{
				Leverage:    s.Leverage, //
				Price:       currentCost,
				Qty:         qty,
				ReduceOnly:  false,
				Side:        "Buy", // @TODO move the value to a constant
				Symbol:      s.CoinPare,
				TimeInForce: "FillOrKill",
				Type:        "Limit", // @TODO move the value to a constant
			}

			sl := 0.0
			if s.StopLossSize > 0 {
				sl = currentCost - s.StopLossSize
				positionToAdd.StopLoss = sl
			}

			utils.LogInfo(fmt.Sprintf("Add position to buy. Current cost: %f, SL: %f, QTY: %f", currentCost, sl, qty))
			positionsToOpen = append(positionsToOpen, positionToAdd)
		}
	}

	if sellPositionExists {
		if !buyOrderTPExists {
			// do open but TP order
			tpPrice := math.Round((sellPositionPrice-s.TakeProfitSize)*100) / 100
			utils.LogInfo(fmt.Sprintf("Add limit order to buy position (as TP). Position cost: %f, TP: %f, QTY: %f", sellPositionPrice, tpPrice, sellPositionQty))
			orderToOpen := structs.DomainOrderRequest{
				Price:       tpPrice,
				Qty:         sellPositionQty,
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
			qty, err := s.getQtyByWallet(currentCost)
			if err != nil {
				return err
			}

			positionToAdd := structs.DomainPositionRequest{
				Leverage:    s.Leverage,
				Price:       currentCost,
				Qty:         qty,
				ReduceOnly:  false,
				Side:        "Sell", //@TODO move the value to a constant
				Symbol:      s.CoinPare,
				TimeInForce: "FillOrKill",
				Type:        "Limit", //@TODO move the value to a constant
			}

			sl := 0.0
			if s.StopLossSize > 0 {
				sl = currentCost + s.StopLossSize
				positionToAdd.StopLoss = sl
			}

			utils.LogInfo(fmt.Sprintf("Add position to sell. Current cost: %f, SL: %f, QTY: %f", currentCost, sl, qty))

			positionsToOpen = append(positionsToOpen, positionToAdd)
		}
	}

	return nil
}

func (s *ScalpByProbability) DoAction() error {

	for _, position := range positionsToOpen {
		positionId, err := services.OpenPosition(s.DomainCode, position)
		if err != nil {
			return err
		}
		utils.LogSuccess(fmt.Sprintf("Position id: %s", positionId))
	}

	for _, order := range ordersToOpen {
		orderId, err := services.OpenOrder(s.DomainCode, order)
		if err != nil {
			return err
		}
		utils.LogSuccess(fmt.Sprintf("Order id: %s", orderId))
	}

	return nil
}

func (s *ScalpByProbability) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}

func (s *ScalpByProbability) getAverageCost() (float64, error) {
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

func (s *ScalpByProbability) getCurrentCost() (float64, error) {
	if len(candles) == 0 {
		return 0, utils.AppError{
			Message: "Current cost not found",
		}
	}
	return candles[0].Close, nil
}

func (s *ScalpByProbability) getCandleCostDiff() float64 {
	minCost := 0.0
	maxCost := 0.0

	for _, candle := range candles {

		high := candle.Open
		low := candle.Open
		if candle.Close > high {
			high = candle.Close
		}
		if candle.Close < low {
			low = candle.Close
		}

		if high > maxCost {
			maxCost = high
		}
		if low < minCost || minCost == 0.0 {
			minCost = low
		}
	}

	return maxCost - minCost
}

func (s *ScalpByProbability) getQtyByWallet(coinCost float64) (float64, error) {

	walletInfo, err := services.LoadWalletInfo(s.DomainCode)

	if err != nil {
		return 0, err
	}

	mainCoinAmount := 0.0

	for _, coin := range walletInfo.Available {
		if coin.Coin == "USDT" {
			mainCoinAmount = coin.Amount
		}
	}

	if mainCoinAmount == 0.0 {
		return 0.0, utils.AppError{Message: "Can not find valid amount of main coin in the wallet"}
	}

	return math.Round(mainCoinAmount/coinCost*s.QtyCoefficient*1000) / 1000, nil
}
