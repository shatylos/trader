package request

import (
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

/**
Spot order docs: https://bybit-exchange.github.io/docs/spot/trade/get-order
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

type SpotOrderResponseTimeInt struct {
	AccountId           string `json:"accountId"`           //	Account ID
	Symbol              string `json:"symbol"`              //	Name of the trading pair
	OrderLinkId         string `json:"orderLinkId"`         //	User-generated order ID
	OrderId             string `json:"orderId"`             //	Order ID
	OrderPrice          string `json:"orderPrice"`          //	Order price
	OrderQty            string `json:"orderQty"`            //	Order quantity
	ExecQty             string `json:"execQty"`             //	Executed quantity
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"` //	Total order quantity. For some historical data, it might less than 0, and that means it is not applicable
	AvgPrice            string `json:"avgPrice"`            //	Average filled price
	Status              string `json:"status"`              //	Order status
	TimeInForce         string `json:"timeInForce"`         //	Time in force
	OrderType           string `json:"orderType"`           //	Order type
	Side                string `json:"side"`                //	Side. BUY, SELL
	StopPrice           string `json:"stopPrice"`           //	Stop price
	IcebergQty          string `json:"icebergQty"`          //	Please ignore
	CreateTime          int64  `json:"createTime"`          //	Order created time in the match engine
	UpdateTime          int64  `json:"updateTime"`          //	Last time order was updated
	IsWorking           string `json:"isWorking"`           //	Is working. 0：valid, 1：invalid
	OrderCategory       int64  `json:"orderCategory"`       //	Order category. 0：normal order; 1：TP/SL order. TP/SL order has this field
	TriggerPrice        string `json:"triggerPrice"`        //	Trigger price. TP/SL order has this field
	BlockTradeId        string `json:"blockTradeId"`        //	Paradigm block trade ID
	CancelType          string `json:"cancelType"`          //	Cancel type. CancelBySmp
	SmpType             string `json:"smpType"`             //	SMP execution type
	SmpGroup            int64  `json:"smpGroup"`            //	Smp group ID. If the uid has no group, it is 0 by default
	SmpOrderId          string `json:"smpOrderId"`          //	The counterparty's orderID which triggers this SMP execution
}

type SpotOrderResponseTimeStr struct {
	AccountId           string `json:"accountId"`           //	Account ID
	Symbol              string `json:"symbol"`              //	Name of the trading pair
	OrderLinkId         string `json:"orderLinkId"`         //	User-generated order ID
	OrderId             string `json:"orderId"`             //	Order ID
	OrderPrice          string `json:"orderPrice"`          //	Order price
	OrderQty            string `json:"orderQty"`            //	Order quantity
	ExecQty             string `json:"execQty"`             //	Executed quantity
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"` //	Total order quantity. For some historical data, it might less than 0, and that means it is not applicable
	AvgPrice            string `json:"avgPrice"`            //	Average filled price
	Status              string `json:"status"`              //	Order status
	TimeInForce         string `json:"timeInForce"`         //	Time in force
	OrderType           string `json:"orderType"`           //	Order type
	Side                string `json:"side"`                //	Side. BUY, SELL
	StopPrice           string `json:"stopPrice"`           //	Stop price
	IcebergQty          string `json:"icebergQty"`          //	Please ignore
	CreateTime          string `json:"createTime"`          //	Order created time in the match engine
	UpdateTime          string `json:"updateTime"`          //	Last time order was updated
	IsWorking           string `json:"isWorking"`           //	Is working. 0：valid, 1：invalid
	OrderCategory       int64  `json:"orderCategory"`       //	Order category. 0：normal order; 1：TP/SL order. TP/SL order has this field
	TriggerPrice        string `json:"triggerPrice"`        //	Trigger price. TP/SL order has this field
	BlockTradeId        string `json:"blockTradeId"`        //	Paradigm block trade ID
	CancelType          string `json:"cancelType"`          //	Cancel type. CancelBySmp
	SmpType             string `json:"smpType"`             //	SMP execution type
	SmpGroup            int64  `json:"smpGroup"`            //	Smp group ID. If the uid has no group, it is 0 by default
	SmpOrderId          string `json:"smpOrderId"`          //	The counterparty's orderID which triggers this SMP execution
}

func CreateSpotOrder(orderRequest SpotOrderRequest, secrets bybitStructs.Secrets) (*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["symbol"] = orderRequest.Symbol
	params["orderQty"] = orderRequest.OrderQty
	params["side"] = orderRequest.Side
	params["orderType"] = orderRequest.OrderType
	params["timeInForce"] = orderRequest.TimeInForce
	params["orderPrice"] = orderRequest.OrderPrice
	params["orderLinkId"] = orderRequest.OrderLinkId

	queryResp, err := apiQueryPost("/spot/v3/private/order", params, secrets)
	if err != nil {
		return nil, err
	}

	return mapSpotOrderResponseTimeStr(queryResp)
}

func GetSpotOrder(domainId string, secrets bybitStructs.Secrets) (*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["orderId"] = domainId
	queryResp, er := apiQueryGet("/spot/v3/private/order", params, secrets)
	if er != nil {
		return nil, er
	}
	order, err := mapSpotOrderResponseTimeStr(queryResp)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetSpotOpenOrderList(coinPare string, secrets bybitStructs.Secrets) ([]*SpotOrderResponseTimeStr, error) {
	params := make(ApiParams, 0)
	params["symbol"] = coinPare
	orders := make([]*SpotOrderResponseTimeStr, 0)

	queryResp, er := apiQueryGet("/spot/v3/private/open-orders", params, secrets)
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

func CancelSpotOrder(orderId string, secrets bybitStructs.Secrets) error {
	params := make(ApiParams, 0)
	params["orderId"] = orderId

	_, err := apiQueryPost("/spot/v3/private/cancel-order", params, secrets)
	if err != nil {
		return err
	}
	return nil
}

func GetSpotOrderHistory(limit int64, secrets bybitStructs.Secrets) ([]*SpotOrderResponseTimeInt, error) {
	params := make(ApiParams, 0)
	params["limit"] = strconv.FormatInt(limit, 10)
	orders := make([]*SpotOrderResponseTimeInt, 0)

	queryResp, er := apiQueryGet("/spot/v3/private/history-orders", params, secrets)
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
		order, err := mapSpotOrderResponseTimeInt(queryRespOrderItem)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func mapSpotOrderResponseTimeInt(queryResp interface{}) (*SpotOrderResponseTimeInt, error) {

	orderResponseBytes, err := json.Marshal(queryResp)
	if err != nil {
		return nil, utils.AppError{
			Message: "Can not Marshal order response for ByBit",
		}
	}

	orderResponse := SpotOrderResponseTimeInt{}
	err = json.Unmarshal(orderResponseBytes, &orderResponse)
	if err != nil {
		return nil, utils.AppError{
			Message: fmt.Sprintf("[mapSpotOrderResponseTimeInt] Can not Unmarshal order response for ByBit: %s", err.Error()),
		}
	}

	return &orderResponse, nil
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

	return &orderResponse, nil
}
