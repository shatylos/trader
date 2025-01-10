package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/domain/domains/bybit/structs"
	"github.com/shatylos/trader/utils"
	"strconv"
	"strings"
)

/**
Spot order docs: https://bybit-exchange.github.io/docs/v5/intro
*/

type SpotOrderRequest struct {
	Symbol        string //	Required. Name of the trading pair
	OrderQty      string //	Required. Order qty. When you place a MARKET BUY order, this param means quote amount. e.g., MARKET BUY BTCUSDT, orderQty should be 200 USDT
	Side          string //	Required. Side. BUY, SELL
	OrderType     string //	Required. Order type
	TimeInForce   string //	Time in force
	OrderPrice    string //	Order price. When the type field is MARKET, the price field is optional. When the type field is LIMIT or LIMIT_MAKER, the price field is required
	OrderLinkId   string //	User-generated order ID
	OrderCategory int    //	Order category. 0：normal order by default; 1：TP/SL order, Required for TP/SL order.
	TriggerPrice  string //	Trigger price. Used for TP/SL order
	SmpType       string //	Smp execution type. What is SMP?
}

type SpotOrderResponseTimeStr struct {
	Symbol      string `json:"symbol"`      //	Name of the trading pair
	OrderLinkId string `json:"orderLinkId"` //	User-generated order ID
	OrderId     string `json:"orderId"`     //	Order ID
	OrderPrice  string `json:"price"`       //	Order price
	OrderQty    string `json:"qty"`         //	Order quantity
	ExecQty     string `json:"cumExecQty"`  //	Executed quantity
	AvgPrice    string `json:"avgPrice"`    //	Average filled price
	Status      string `json:"orderStatus"` //	Order status
	TimeInForce string `json:"timeInForce"` //	Time in force
	OrderType   string `json:"orderType"`   //	Order type
	Side        string `json:"side"`        //	Side. BUY, SELL
	CreateTime  string `json:"createdTime"` //	Order created time in the match engine
	UpdateTime  string `json:"updatedTime"` //	Last time order was updated
}

func CreateSpotOrder(orderRequest SpotOrderRequest, secrets bybitStructs.Secrets) (*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["category"] = "spot"
	params["symbol"] = orderRequest.Symbol
	params["qty"] = orderRequest.OrderQty
	params["side"] = orderRequest.Side
	params["orderType"] = orderRequest.OrderType
	params["timeInForce"] = orderRequest.TimeInForce
	params["price"] = orderRequest.OrderPrice
	params["orderLinkId"] = orderRequest.OrderLinkId

	queryResp, err := apiQueryPost("/v5/order/create", params, secrets)
	if err != nil {
		return nil, err
	}

	return mapSpotOrderResponseTimeStr(queryResp)
}

func GetSpotOrder(domainId string, secrets bybitStructs.Secrets) (*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["category"] = "spot"
	params["orderId"] = domainId
	queryResp, er := apiQueryGet("/v5/order/realtime", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse order list query response"}
	}
	if queryRespMap["list"] == nil || len(queryRespMap["list"].([]interface{})) == 0 {
		return &SpotOrderResponseTimeStr{
			OrderId:    domainId,
			Status:     "CANCELED",
			AvgPrice:   "0",
			ExecQty:    "0",
			CreateTime: "0",
			UpdateTime: "0",
		}, nil
	}

	queryResponseOrder := queryRespMap["list"].([]interface{})[0]

	order, err := mapSpotOrderResponseTimeStr(queryResponseOrder)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetSpotOpenOrderList(coinPare string, secrets bybitStructs.Secrets) ([]*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["symbol"] = coinPare
	params["category"] = "spot"
	params["openOnly"] = "0"
	orders := make([]*SpotOrderResponseTimeStr, 0)

	queryResp, er := apiQueryGet("/v5/order/realtime", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse order list query response"}
	}

	if queryRespMap["list"] == nil {
		return orders, nil
	}

	queryRespData, ok := queryRespMap["list"].([]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse order list data values"}
	}

	for _, queryRespOrderItem := range queryRespData {
		order, err := mapSpotOrderResponseTimeStr(queryRespOrderItem)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func CancelSpotOrder(orderId string, secrets bybitStructs.Secrets, coinPare string) error {
	params := make(ApiParams, 0)
	params["orderId"] = orderId
	params["category"] = "spot"
	params["symbol"] = coinPare

	_, err := apiQueryPost("/v5/order/cancel", params, secrets)
	if err != nil {
		return err
	}
	return nil
}

func GetSpotOrderHistory(limit int64, secrets bybitStructs.Secrets, coinPare string) ([]*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["limit"] = strconv.FormatInt(limit, 10)
	params["symbol"] = coinPare
	params["category"] = "spot"
	params["orderStatus"] = "Filled"

	orders := make([]*SpotOrderResponseTimeStr, 0)

	queryResp, er := apiQueryGet("/v5/order/history", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse order list query response"}
	}

	if queryRespMap["list"] == nil {
		return orders, nil
	}

	queryRespData, ok := queryRespMap["list"].([]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse order list data values"}
	}

	for _, queryRespOrderItem := range queryRespData {
		order, err := mapSpotOrderResponseTimeStr(queryRespOrderItem)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func mapSpotOrderResponseTimeStr(queryResp interface{}) (*SpotOrderResponseTimeStr, error) {

	orderResponseBytes, err := json.Marshal(queryResp)
	if err != nil {
		return nil, utils.AppError{
			Message: "Can not Marshal order response for ByBit",
		}
	}

	orderResponse := SpotOrderResponseTimeStr{}
	err = json.Unmarshal(orderResponseBytes, &orderResponse)
	if err != nil {
		return nil, utils.AppError{
			Message: fmt.Sprintf("[mapSpotOrderResponseTimeStr] Can not Unmarshal order response for ByBit: %s", err.Error()),
		}
	}

	orderResponse.Status = strings.ToUpper(orderResponse.Status)
	orderResponse.Side = strings.ToUpper(orderResponse.Side)

	if orderResponse.Status == "CANCELLED" {
		orderResponse.Status = "CANCELED"
	}
	if orderResponse.Status == "PARTIALLYFILLEDCANCELED" {
		orderResponse.Status = "PARTIALLY_FILLED"
	}

	return &orderResponse, nil
}
