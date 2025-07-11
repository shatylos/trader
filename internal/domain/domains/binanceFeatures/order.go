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
	"strconv"
)

type OrderResponse struct {
	OrderId                 int64  `json:"orderId"`
	Symbol                  string `json:"symbol"`
	Status                  string `json:"status"`
	ClientOrderId           string `json:"clientOrderId"`
	Price                   string `json:"price"`
	AvgPrice                string `json:"avgPrice"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	CumQuote                string `json:"cumQuote"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	ReduceOnly              bool   `json:"reduceOnly"`
	ClosePosition           bool   `json:"closePosition"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	StopPrice               string `json:"stopPrice"`
	WorkingType             string `json:"workingType"`
	PriceProtect            bool   `json:"priceProtect"`
	OrigType                string `json:"origType"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int    `json:"goodTillDate"`
	Time                    int64  `json:"time"`
	UpdateTime              int64  `json:"updateTime"`
	ErrorCode               int    `json:"code"`
	ErrorMessage            string `json:"msg"`
}

type orderStatusesStruct struct {
	New             string
	PendingNew      string
	PartiallyFilled string
	Filled          string
	Canceled        string
	Rejected        string
	Expired         string
	ExpiredInMatch  string
}

var orderStatuses = orderStatusesStruct{
	New:             "NEW",
	PendingNew:      "PENDING_NEW",
	PartiallyFilled: "PARTIALLY_FILLED",
	Filled:          "FILLED",
	Canceled:        "CANCELED",
	Rejected:        "REJECTED",
	Expired:         "EXPIRED",
	ExpiredInMatch:  "EXPIRED_IN_MATCH",
}

type orderTypesStruct struct {
	Limit           string
	Market          string
	StopLoss        string
	StopLossLimit   string
	TakeProfit      string
	TakeProfitLimit string
	LimitMaker      string
}

var orderTypes = orderTypesStruct{
	Limit:           "LIMIT",
	Market:          "MARKET",
	StopLoss:        "STOP_LOSS",
	StopLossLimit:   "STOP_LOSS_LIMIT",
	TakeProfit:      "TAKE_PROFIT",
	TakeProfitLimit: "TAKE_PROFIT_LIMIT",
	LimitMaker:      "LIMIT_MAKER",
}
var orderTypeTP = "TAKE_PROFIT_MARKET"
var orderTypeSL = "STOP_MARKET"

type orderSidesStruct struct {
	Buy  string
	Sell string
}

var orderSideBuy = "BUY"
var orderSideSell = "SELL"

func (d *DomainBinanceFutures) GetOrder(orderId string, coinPare string) (order domainStructs.DomainOrder, err error) {
	orderApiRequest := request.ApiGetRequest{
		Uri: "/fapi/v1/order",
		ApiParams: binanceStructs.ApiParams{
			"symbol":  coinPare,
			"orderId": orderId,
		},
		Secrets: d.secrets,
	}
	var apiResponse binanceStructs.ApiResponse
	apiResponse, err = orderApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerOrder := OrderResponse{}
	err = json.Unmarshal(apiResponse, &providerOrder)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarhsal Binance GetOrder API response. Raw data: %s", apiResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	var respOrderId, status, orderType, side string
	respOrderId, err = _type.ToString(providerOrder.OrderId)
	if err != nil {
		return
	}

	status, err = d.orderStatusPtoD(providerOrder.Status)
	if err != nil {
		return
	}
	orderType, err = d.orderTypePtoD(providerOrder.Type)
	if err != nil {
		return
	}
	side, err = d.orderSidePtoD(providerOrder.Side)
	if err != nil {
		return
	}

	var price, qty float64
	price, err = _type.ToFloat64(providerOrder.Price)
	if err != nil {
		msg := fmt.Sprintf("Binance GetOrder price (%s) can not be parsed as float64", providerOrder.Price)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	qty, err = _type.ToFloat64(providerOrder.OrigQty)
	if err != nil {
		msg := fmt.Sprintf("Binance GetOrder qty (%s) can not be parsed as float64", providerOrder.OrigQty)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	order = domainStructs.DomainOrder{
		CreatedTime: providerOrder.Time / 1000,
		OrderId:     respOrderId,
		OrderStatus: status,
		OrderType:   orderType,
		Price:       price,
		Qty:         qty,
		ReduceOnly:  providerOrder.ReduceOnly,
		Side:        side,
		Symbol:      providerOrder.Symbol,
		TimeInForce: providerOrder.TimeInForce,
		UpdatedTime: providerOrder.UpdateTime / 1000,
	}

	return
}

func (d *DomainBinanceFutures) GetOrders(coinPare string) (orders []OrderResponse, err error) {
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

	err = json.Unmarshal(orderResponse, &orders)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance order response data: %s", orderResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}

func (d *DomainBinanceFutures) DeleteOrder(orderId int64, coinPare string) (err error) {
	delApiRequest := request.ApiDeleteRequest{
		Uri: "/fapi/v1/order",
		ApiParams: binanceStructs.ApiParams{
			"symbol":  coinPare,
			"orderId": strconv.FormatInt(orderId, 10),
		},
		Secrets: d.secrets,
	}
	var apiResponse binanceStructs.ApiResponse
	apiResponse, err = delApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerOrder := OrderResponse{}
	err = json.Unmarshal(apiResponse, &providerOrder)
	if err != nil {
		return
	}
	if providerOrder.Status != orderStatuses.Canceled {
		msg := fmt.Sprintf("Order status is \"%s\". Expected status %s", providerOrder.Status, orderStatuses.Canceled)
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}
	return
}

func (d *DomainBinanceFutures) orderStatusPtoD(providerStatus string) (domainStatus string, err error) {
	switch providerStatus {
	case orderStatuses.New:
	case orderStatuses.PendingNew:
		domainStatus = domainStructs.OrderStatuses.Open
		break
	case orderStatuses.PartiallyFilled:
		domainStatus = domainStructs.OrderStatuses.PartiallyFilled
		break
	case orderStatuses.Filled:
		domainStatus = domainStructs.OrderStatuses.Filled
		break
	case orderStatuses.Canceled:
	case orderStatuses.Rejected:
	case orderStatuses.Expired:
	case orderStatuses.ExpiredInMatch:
		domainStatus = domainStructs.OrderStatuses.Canceled
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance provider order status value: \"%s\"", providerStatus)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}
	return
}

func (d *DomainBinanceFutures) orderTypePtoD(providerType string) (domainType string, err error) {
	switch providerType {
	case orderTypes.Limit:
	case orderTypes.StopLoss:
	case orderTypes.StopLossLimit:
	case orderTypes.TakeProfit:
	case orderTypes.TakeProfitLimit:
	case orderTypes.LimitMaker:
		domainType = domainStructs.OrderTypes.Limit
		break
	case orderTypes.Market:
		domainType = domainStructs.OrderTypes.Market
		break
	default:
		msg := fmt.Sprintf("Binance order type (\"%s\") can not be mapped to domain order type", providerType)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}
	return
}

func (d *DomainBinanceFutures) orderSidePtoD(providerSide string) (domainSide string, err error) {
	switch providerSide {
	case orderSideBuy:
		domainSide = domainStructs.OrderSideBuy
		break
	case orderSideSell:
		domainSide = domainStructs.OrderSideSell
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance provider order side value: \"%s\"", providerSide)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}

	return
}

func (d *DomainBinanceFutures) reverseOrderSide(orderSide string) (reverseSide string, err error) {
	switch orderSide {
	case domainStructs.OrderSideBuy:
		reverseSide = domainStructs.OrderSideSell
		break
	case domainStructs.OrderSideSell:
		reverseSide = domainStructs.OrderSideBuy
		break
	default:
		msg := fmt.Sprintf("Unexpected Binance order side value: \"%s\"", orderSide)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
	}
	return
}
