package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
)

type OrderRequest struct {
	TimeInForce    string  //			mandatory		GoodTillCancel | ImmediateOrCancel | FillOrKill | PostOnly
	ReduceOnly     bool    //			mandatory
	CloseOnTrigger bool    //			mandatory
	Side           string  //			mandatory		Buy | Sell
	Symbol         string  //			mandatory
	OrderType      string  //			mandatory		Limit | Market
	Qty            float64 //			mandatory
	Price          float64 //			not mandatory
	TakeProfit     float64 //			not mandatory
	StopLoss       float64 //			not mandatory
	OrderLinkId    string  //			not mandatory
	SlTriggerBy    string  //			not mandatory
	TpTriggerBy    string  //			not mandatory
}

type OrderResponse struct {
	Symbol         string  `json:"symbol"`         // -> BTCUSDT
	Side           string  `json:"side"`           // -> Buy
	Price          string  `json:"price"`          // -> 24553
	AvgPrice       string  `json:"avgPrice"`       // -> 24553
	CumExecValue   string  `json:"cumExecValue"`   // -> 0
	StopLoss       string  `json:"stopLoss"`       // -> 22800
	OrderType      string  `json:"orderType"`      // -> Market
	TimeInForce    string  `json:"timeInForce"`    // -> ImmediateOrCancel
	OrderStatus    string  `json:"orderStatus"`    // -> Created
	ReduceOnly     bool    `json:"reduceOnly"`     // -> false
	CloseOnTrigger bool    `json:"closeOnTrigger"` // -> false
	CreatedTime    string  `json:"createdTime"`    // -> 2022-08-17T20:45:45Z
	TakeProfit     string  `json:"takeProfit"`     // -> 23600
	SlTriggerBy    string  `json:"slTriggerBy"`    // -> LastPrice
	PositionIdx    float64 `json:"positionIdx"`    // -> 1
	CumExecQty     string  `json:"cumExecQty"`     // -> 0
	TpTriggerBy    string  `json:"tpTriggerBy"`    // -> LastPrice
	OrderId        string  `json:"orderId"`        // -> 10772450-26ce-4c92-8bf6-5cf29d190a50
	Qty            string  `json:"qty"`            // -> 0.005
	CumExecFee     string  `json:"cumExecFee"`     // -> 0
	OrderLinkId    string  `json:"orderLinkId"`    // -> 123qwe
	UpdatedTime    string  `json:"updatedTime"`    // -> 2022-08-17T20:45:45Z
}

func CreateOrder(orderRequest OrderRequest, secrets bybitStructs.Secrets) (*OrderResponse, error) {
	params := make(ApiParams, 0)
	params["category"] = "linear"
	params["isLeverage"] = "1"
	params["timeInForce"] = orderRequest.TimeInForce
	params["reduceOnly"] = orderRequest.ReduceOnly
	params["closeOnTrigger"] = orderRequest.CloseOnTrigger

	params["side"] = orderRequest.Side
	params["symbol"] = orderRequest.Symbol
	params["orderType"] = orderRequest.OrderType
	params["qty"] = fmt.Sprintf("%g", orderRequest.Qty)
	params["marketUnit"] = "baseCoin"

	if orderRequest.Price > 0 {
		params["price"] = fmt.Sprintf("%g", orderRequest.Price)
	}
	if orderRequest.TakeProfit > 0 {
		params["takeProfit"] = fmt.Sprintf("%g", orderRequest.TakeProfit)
	}
	if orderRequest.StopLoss > 0 {
		params["stopLoss"] = fmt.Sprintf("%g", orderRequest.StopLoss)
	}
	if orderRequest.OrderLinkId != "" {
		params["orderLinkId"] = orderRequest.OrderLinkId
	}
	if orderRequest.TpTriggerBy != "" {
		params["tpTriggerBy"] = orderRequest.TpTriggerBy
	}
	if orderRequest.SlTriggerBy != "" {
		params["slTriggerBy"] = orderRequest.SlTriggerBy
	}

	queryResp, err := apiQueryPost("/v5/order/create", params, secrets)
	if err != nil {
		return nil, err
	}

	return mapOrderResponse(queryResp)
}

func GetOrderList(coinPare string, orderId string, secrets bybitStructs.Secrets) ([]*OrderResponse, error) {
	params := make(ApiParams, 0)
	params["category"] = "linear"
	if coinPare != "" {
		params["symbol"] = coinPare
	}
	if orderId != "" {
		params["orderId"] = orderId
	}
	params["openOnly"] = "0"
	orders := make([]*OrderResponse, 0)

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
		order, err := mapOrderResponse(queryRespOrderItem)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func mapOrderResponse(queryResp interface{}) (*OrderResponse, error) {

	orderResponseBytes, err := json.Marshal(queryResp)
	if err != nil {
		return nil, tools.AppError{
			Message: "Can not Marshal order response for ByBit",
		}
	}

	orderResponse := OrderResponse{}
	err = json.Unmarshal(orderResponseBytes, &orderResponse)
	if err != nil {
		return nil, tools.AppError{
			Message: fmt.Sprintf("[mapOrderResponse] Can not Unmarshal order response for ByBit: %s", err.Error()),
		}
	}

	return &orderResponse, nil
}
