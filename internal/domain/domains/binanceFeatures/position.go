package binanceFeatures

import (
	"encoding/json"
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/request"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"math"
	"strconv"
)

type ProviderPositionRisk struct {
	Symbol                 string `json:"symbol"`                 // BTCUSDT
	PositionSide           string `json:"positionSide"`           // BOTH
	PositionAmt            string `json:"positionAmt"`            // -0.001
	EntryPrice             string `json:"entryPrice"`             // 105326.5
	BreakEvenPrice         string `json:"breakEvenPrice"`         // 105284.3694
	MarkPrice              string `json:"markPrice"`              // 105336.60000000
	UnRealizedProfit       string `json:"unRealizedProfit"`       // -0.01010000
	LiquidationPrice       string `json:"liquidationPrice"`       // 15045103.95358566
	IsolatedMargin         string `json:"isolatedMargin"`         // 0
	Notional               string `json:"notional"`               // -105.33660000
	MarginAsset            string `json:"marginAsset"`            // USDT
	IsolatedWallet         string `json:"isolatedWallet"`         // 0
	InitialMargin          string `json:"initialMargin"`          // 5.26682999
	MaintMargin            string `json:"maintMargin"`            // 0.42134640
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 5.26682999
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 0
	Adl                    int    `json:"adl"`                    // 3
	BidNotional            string `json:"bidNotional"`            // 0
	AskNotional            string `json:"askNotional"`            // 0
	UpdateTime             int64  `json:"updateTime"`             // 1749309239526
}

type PositionTypesStruct struct {
	Limit              string
	Market             string
	Stop               string
	TakeProfit         string
	StopMarket         string
	TakeProfitMarket   string
	TrailingStopMarket string
}

var PositionTypes = PositionTypesStruct{
	Limit:              "LIMIT",
	Market:             "MARKET",
	Stop:               "STOP",
	TakeProfit:         "TAKE_PROFIT",
	StopMarket:         "STOP_MARKET",
	TakeProfitMarket:   "TAKE_PROFIT_MARKET",
	TrailingStopMarket: "TRAILING_STOP_MARKET",
}

const (
	positionSideShort = "SHORT"
	positionSideLong  = "LONG"
	positionSideBoth  = "BOTH"
)
const (
	positionModeHedge  = "HEDGE"
	positionModeOneWay = "ONEWAY"
)

func (d *DomainBinanceFutures) OpenPosition(positionRequest domainStructs.DomainPositionRequest) (positionId string, err error) {

	if d.positionMode != positionModeHedge {
		err = d.setPositionMode(positionModeHedge)
		if err != nil {
			return
		}
	}

	var orderSide, positionSide, positionType, price, qty string
	orderSide, positionSide, err = d.positionSideDtoP(positionRequest.Side)
	if err != nil {
		return
	}

	positionType, err = d.positionTypeDtoP(positionRequest.Type)
	if err != nil {
		return
	}

	batchOrders := make([]binanceStructs.ApiParams, 1)
	paramOrder := binanceStructs.ApiParams{
		"symbol":       positionRequest.Symbol,
		"side":         orderSide,
		"positionSide": positionSide, // Default BOTH for One-way Mode ; LONG or SHORT for Hedge Mode. It must be sent in Hedge Mode.
		"type":         positionType,
	}

	qty, err = _type.ToString(positionRequest.Qty)
	if err != nil {
		return
	}
	paramOrder["quantity"] = qty // Cannot be sent with closePosition=true(Close-All)

	if positionRequest.Type == domainStructs.PositionTypes.Limit {
		price, err = _type.ToString(positionRequest.Price)
		if err != nil {
			return
		}
		paramOrder["price"] = price
	}
	batchOrders[0] = paramOrder

	var orderSideReverse string
	orderSideReverse, err = d.reverseOrderSide(orderSide)
	if err != nil {
		return
	}

	if positionRequest.TakeProfit > 0 {
		var tp string
		tp, err = _type.ToString(positionRequest.TakeProfit)
		if err != nil {
			return
		}

		tpOrder := binanceStructs.ApiParams{
			"symbol":        positionRequest.Symbol,
			"side":          orderSideReverse,
			"positionSide":  positionSide, // Default BOTH for One-way Mode ; LONG or SHORT for Hedge Mode. It must be sent in Hedge Mode.
			"type":          PositionTypes.TakeProfitMarket,
			"stopPrice":     tp,     // Used with STOP/STOP_MARKET or TAKE_PROFIT/TAKE_PROFIT_MARKET orders.
			"closePosition": "true", // true, false；Close-All，used with STOP_MARKET or TAKE_PROFIT_MARKET.
		}
		batchOrders = append(batchOrders, tpOrder)
	}
	if positionRequest.StopLoss > 0 {
		var sl string
		sl, err = _type.ToString(positionRequest.StopLoss)
		if err != nil {
			return
		}

		slOrder := binanceStructs.ApiParams{
			"symbol":        positionRequest.Symbol,
			"side":          orderSideReverse,
			"positionSide":  positionSide, // Default BOTH for One-way Mode ; LONG or SHORT for Hedge Mode. It must be sent in Hedge Mode.
			"type":          PositionTypes.StopMarket,
			"stopPrice":     sl,     // Used with STOP/STOP_MARKET or TAKE_PROFIT/TAKE_PROFIT_MARKET orders.
			"closePosition": "true", // true, false；Close-All，used with STOP_MARKET or TAKE_PROFIT_MARKET.
		}
		batchOrders = append(batchOrders, slOrder)
	}

	params := binanceStructs.ApiParams{
		"batchOrders": batchOrders,
	}
	positionApiRequest := request.ApiPostRequest{
		Uri:       "/fapi/v1/batchOrders",
		ApiParams: params,
		Secrets:   d.secrets,
	}

	var apiResponse binanceStructs.ApiResponse
	apiResponse, err = positionApiRequest.DoRequest()
	if err != nil {
		return
	}
	logger.Info(fmt.Sprintf("%s", apiResponse))

	var providerOrders []OrderResponse
	err = json.Unmarshal(apiResponse, &providerOrders)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarhsal Binance GetOrder API response. Raw data: %s", apiResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	var tpSlCreatedIds []int64
	for _, providerOrder := range providerOrders {
		if providerOrder.ErrorCode != 0 {
			msg := fmt.Sprintf("Error creating order for position. Code: %d, Message: %s", providerOrder.ErrorCode, providerOrder.ErrorMessage)
			logger.Error(msg)
			err = tools.AppError{
				Message:     msg,
				ParentError: err,
			}
		} else {
			if providerOrder.Type == orderTypes.Market || providerOrder.Type == orderTypes.Limit {
				positionId = strconv.FormatInt(providerOrder.OrderId, 10)
			} else {
				tpSlCreatedIds = append(tpSlCreatedIds, providerOrder.OrderId)
			}
		}
	}
	if err != nil {
		return
	}
	return
}

func (d *DomainBinanceFutures) GetPosition(coinPare string) (position domainStructs.DomainPosition, err error) {

	positionApiRequest := request.ApiGetRequest{
		Uri: "/fapi/v3/positionRisk",
		ApiParams: binanceStructs.ApiParams{
			"symbol": coinPare,
		},
		Secrets: d.secrets,
	}

	var positionResponse binanceStructs.ApiResponse
	positionResponse, err = positionApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerPositions := make([]ProviderPositionRisk, 0)
	err = json.Unmarshal(positionResponse, &providerPositions)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance position response data: %s", positionResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	for _, providerPosition := range providerPositions {
		if providerPosition.Symbol == coinPare {
			var markPrice float64
			markPrice, err = _type.ToFloat64(providerPosition.MarkPrice)
			if err != nil {
				return
			}

			var totalPnl float64
			totalPnl, err = _type.ToFloat64(providerPosition.UnRealizedProfit)
			if err != nil {
				return
			}

			var size float64
			size, err = _type.ToFloat64(providerPosition.PositionAmt)
			if err != nil {
				return
			}
			size = math.Abs(size)

			var entryPrice float64
			entryPrice, err = _type.ToFloat64(providerPosition.EntryPrice)
			if err != nil {
				return
			}
			position = domainStructs.DomainPosition{
				Leverage:  0, // check to remove
				AvgPrice:  entryPrice,
				MarkPrice: markPrice,
				Size:      size,
				TotalPnl:  totalPnl,
				Side:      providerPosition.PositionSide,
				Symbol:    providerPosition.Symbol,
			}
			break
		}
	}

	orderApiRequest := request.ApiGetRequest{
		Uri: "/fapi/v1/openOrders",
		ApiParams: binanceStructs.ApiParams{
			"symbol": coinPare,
		},
		Secrets: d.secrets,
	}
	var orderResponse binanceStructs.ApiResponse
	orderResponse, err = orderApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerOrders := make([]OrderResponse, 0)
	err = json.Unmarshal(orderResponse, &providerOrders)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance order response data: %s", orderResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	for _, providerOrder := range providerOrders {
		if providerOrder.Symbol == coinPare && providerOrder.Status == orderStatuses.New && providerOrder.Type == orderTypeTP {
			if position.TakeProfit > 0 {
				msg := fmt.Sprintf("Second take profit value got from provider order list")
				logger.Error(msg)
				err = tools.AppError{
					Message: msg,
				}
				return
			}
			var stopPrice float64
			stopPrice, err = _type.ToFloat64(providerOrder.StopPrice)
			if err != nil {
				return
			}
			position.TakeProfit = stopPrice
		}
		if providerOrder.Symbol == coinPare && providerOrder.Status == orderStatuses.New && providerOrder.Type == orderTypeSL {
			if position.StopLoss > 0 {
				msg := fmt.Sprintf("Second stop loss value got from provider order list")
				logger.Error(msg)
				err = tools.AppError{
					Message: msg,
				}
				return
			}
			var stopPrice float64
			stopPrice, err = _type.ToFloat64(providerOrder.StopPrice)
			if err != nil {
				return
			}
			position.StopLoss = stopPrice
		}
	}

	return
}

func (d *DomainBinanceFutures) positionSideDtoP(domainSide string) (providerOrderSide string, providerPositionSide string, err error) {
	switch domainSide {
	case domainStructs.PositionSideLong:
		providerOrderSide = orderSideBuy
		providerPositionSide = positionSideLong
		break
	case domainStructs.PositionSideShort:
		providerOrderSide = orderSideSell
		providerPositionSide = positionSideShort
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance domain position side value: \"%s\"", domainSide)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}

	return
}

func (d *DomainBinanceFutures) positionTypeDtoP(domainType string) (providerType string, err error) {
	switch domainType {
	case domainStructs.PositionTypes.Market:
		providerType = PositionTypes.Market
		break
	case domainStructs.PositionTypes.Limit:
		providerType = PositionTypes.Limit
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance domain position type value: \"%s\"", domainType)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}
	return
}

func (d *DomainBinanceFutures) setPositionMode(positionMode string) (err error) {
	var dualSidePosition string
	switch positionMode {
	case positionModeHedge:
		dualSidePosition = "true"
		break
	case positionModeOneWay:
		dualSidePosition = "false"
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance position mode: \"%s\"", positionMode)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	apiRequest := request.ApiPostRequest{
		Uri: "/fapi/v1/positionSide/dual",
		ApiParams: binanceStructs.ApiParams{
			"dualSidePosition": dualSidePosition,
		},
		Secrets: d.secrets,
	}
	var rawResponse binanceStructs.ApiResponse
	rawResponse, err = apiRequest.DoRequest()
	if err != nil {
		// skip the error "No need to change position side."
		if request.IsErrorCode(rawResponse, -4059) {
			err = nil
		} else {
			return
		}
	}
	d.positionMode = positionMode
	return
}
