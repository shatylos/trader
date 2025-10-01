package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
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

type CreateSpotOrderResponse struct {
	OrderLinkId string `json:"orderLinkId"` //	User-generated order ID
	OrderId     string `json:"orderId"`     //	Order ID
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

func CreateSpotOrder(orderRequest SpotOrderRequest, secrets bybitStructs.Secrets) (orderResponse CreateSpotOrderResponse, err error) {
	params := make(ApiParams, 0)
	params["category"] = "spot"
	params["symbol"] = orderRequest.Symbol
	params["qty"] = orderRequest.OrderQty
	params["side"] = orderRequest.Side
	params["orderType"] = orderRequest.OrderType
	params["timeInForce"] = orderRequest.TimeInForce
	if orderRequest.OrderPrice != "" {
		params["price"] = orderRequest.OrderPrice
	}
	params["orderLinkId"] = orderRequest.OrderLinkId
	params["marketUnit"] = "baseCoin"
	var queryResp interface{}
	queryResp, err = apiQueryPost("/v5/order/create", params, secrets)
	if err != nil {
		return
	}

	return mapCreateSpotOrderResponse(queryResp)
}

func GetSpotOrder(domainId string, secrets bybitStructs.Secrets) (orderResponse SpotOrderResponseTimeStr, err error) {
	params := make(ApiParams, 0)
	params["category"] = "spot"
	params["orderId"] = domainId
	var queryResp interface{}
	queryResp, err = apiQueryGet("/v5/order/realtime", params, secrets)
	if err != nil {
		return
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		msg := "Can not parse order list query response"
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	if queryRespMap["list"] == nil || len(queryRespMap["list"].([]interface{})) == 0 {
		orderResponse.OrderId = domainId
		orderResponse.Status = domainStructs.OrderStatuses.Canceled
		return
	}

	queryResponseOrder := queryRespMap["list"].([]interface{})[0]

	orderResponse, err = mapSpotOrderResponseTimeStr(queryResponseOrder)
	if err != nil {
		return
	}
	return
}

func GetSpotOpenOrderList(coinPare string, secrets bybitStructs.Secrets) ([]SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["symbol"] = coinPare
	params["category"] = "spot"
	params["openOnly"] = "0"
	orders := make([]SpotOrderResponseTimeStr, 0)

	queryResp, er := apiQueryGet("/v5/order/realtime", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		return nil, tools.AppError{Message: "Can not parse order list query response"}
	}

	if queryRespMap["list"] == nil {
		return orders, nil
	}

	queryRespData, ok := queryRespMap["list"].([]interface{})
	if !ok {
		return nil, tools.AppError{Message: "Can not parse order list data values"}
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

func GetSpotOrderHistory(limit int64, secrets bybitStructs.Secrets, coinPare string) ([]SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["limit"] = strconv.FormatInt(limit, 10)
	params["symbol"] = coinPare
	params["category"] = "spot"
	params["orderStatus"] = "Filled"

	orders := make([]SpotOrderResponseTimeStr, 0)

	queryResp, er := apiQueryGet("/v5/order/history", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		return nil, tools.AppError{Message: "Can not parse order list query response"}
	}

	if queryRespMap["list"] == nil {
		return orders, nil
	}

	queryRespData, ok := queryRespMap["list"].([]interface{})
	if !ok {
		return nil, tools.AppError{Message: "Can not parse order list data values"}
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

func mapCreateSpotOrderResponse(queryResp interface{}) (orderResponse CreateSpotOrderResponse, err error) {
	var responseBytes []byte
	responseBytes, err = json.Marshal(queryResp)
	if err != nil {
		msg := "Can not Marshal create order response for ByBit"
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	err = json.Unmarshal(responseBytes, &orderResponse)
	if err != nil {
		msg := fmt.Sprintf("[mapCreateSpotOrderResponse] Can not Unmarshal order response for ByBit: %s", err.Error())
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}

func mapSpotOrderResponseTimeStr(queryResp interface{}) (orderResponse SpotOrderResponseTimeStr, err error) {

	var orderResponseBytes []byte
	orderResponseBytes, err = json.Marshal(queryResp)
	if err != nil {
		msg := "Can not Marshal order response for ByBit"
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	err = json.Unmarshal(orderResponseBytes, &orderResponse)
	if err != nil {
		msg := fmt.Sprintf("[mapSpotOrderResponseTimeStr] Can not Unmarshal order response for ByBit: %s", err.Error())
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	switch strings.ToUpper(orderResponse.Status) {
	case "NEW":
		orderResponse.Status = domainStructs.OrderStatuses.New
		break
	case "OPEN":
		orderResponse.Status = domainStructs.OrderStatuses.Open
		break
	case "FILLED":
		orderResponse.Status = domainStructs.OrderStatuses.Filled
		break
	case "PARTIALLY_FILLED",
		"PARTIALLYFILLED",
		"PARTIALLYFILLEDCANCELED":
		orderResponse.Status = domainStructs.OrderStatuses.PartiallyFilled
		break
	case "CANCELED",
		"CANCELLED":
		orderResponse.Status = domainStructs.OrderStatuses.Canceled
		break
	case "":
		logger.Warning("Bybit order response status is empty")
		orderResponse.Status = ""
		break
	default:
		msg := fmt.Sprintf("Bybit order status (\"%s\") can not be mapped to domain value", orderResponse.Status)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	orderResponse.Side = strings.ToUpper(orderResponse.Side)

	return
}
