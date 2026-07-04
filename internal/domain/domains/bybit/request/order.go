package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools/apperrors"
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
	SlTriggerBy    string  //			not mandatory
	TpTriggerBy    string  //			not mandatory
	TpslMode       string  //			not mandatory	Full | Partial; Partial is required for a limit take profit
	TpOrderType    string  //			not mandatory	Market | Limit
	TpLimitPrice   float64 //			not mandatory	limit price of the take profit order, requires TpslMode Partial and TpOrderType Limit
}

type TpSlRequest struct {
	Symbol     string
	TakeProfit float64
	StopLoss   float64
	TpSlMode   string // Full | Partial
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
	OrderLinkId    string  `json:"orderLinkId"`    //
	Qty            string  `json:"qty"`            // -> 0.005
	CumExecFee     string  `json:"cumExecFee"`     // -> 0
	UpdatedTime    string  `json:"updatedTime"`    // -> 2022-08-17T20:45:45Z
}

func CreateOrder(orderRequest OrderRequest, secrets bybitStructs.Secrets) (orderResponse *OrderResponse, err error) {
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

	if orderRequest.Price > 0 {
		params["price"] = fmt.Sprintf("%g", orderRequest.Price)
	}
	if orderRequest.TakeProfit > 0 {
		params["takeProfit"] = fmt.Sprintf("%g", orderRequest.TakeProfit)
	}
	if orderRequest.StopLoss > 0 {
		params["stopLoss"] = fmt.Sprintf("%g", orderRequest.StopLoss)
	}
	if orderRequest.TpTriggerBy != "" {
		params["tpTriggerBy"] = orderRequest.TpTriggerBy
	}
	if orderRequest.SlTriggerBy != "" {
		params["slTriggerBy"] = orderRequest.SlTriggerBy
	}
	if orderRequest.TpslMode != "" {
		params["tpslMode"] = orderRequest.TpslMode
	}
	if orderRequest.TpOrderType != "" {
		params["tpOrderType"] = orderRequest.TpOrderType
	}
	if orderRequest.TpLimitPrice > 0 {
		params["tpLimitPrice"] = fmt.Sprintf("%g", orderRequest.TpLimitPrice)
	}

	uri := "/v5/order/create"
	var queryResp interface{}
	queryResp, err = apiQueryPost(uri, params, secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error sending post request, uri: %s, params: %s", uri, params)
		return
	}

	orderResponse, err = mapOrderResponse(queryResp)
	if err != nil {
		err = apperrors.Wrap(err, "error mapping order response, query response: %s", queryResp)
		return
	}
	return
}

func ModifyTpSl(orderRequest TpSlRequest, secrets bybitStructs.Secrets) (err error) {

	params := ApiParams{}
	params["category"] = "linear"
	params["symbol"] = orderRequest.Symbol
	params["tpslMode"] = orderRequest.TpSlMode
	params["positionIdx"] = "0"

	if orderRequest.TakeProfit > 0 {
		params["takeProfit"] = fmt.Sprintf("%g", orderRequest.TakeProfit)
	}
	if orderRequest.StopLoss > 0 {
		params["stopLoss"] = fmt.Sprintf("%g", orderRequest.StopLoss)
	}

	uri := "/v5/position/trading-stop"
	_, err = apiQueryPost(uri, params, secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error sending post request, uri: %s, params: %s", uri, params)
		return
	}
	return
}

func GetOrderList(coinPare string, orderId string, secrets bybitStructs.Secrets) (orders []*OrderResponse, err error) {
	params := make(ApiParams, 0)
	params["category"] = "linear"
	if coinPare != "" {
		params["symbol"] = coinPare
	}
	if orderId != "" {
		params["orderId"] = orderId
	}
	params["openOnly"] = "0"

	uri := "/v5/order/realtime"
	var queryResp interface{}
	queryResp, err = apiQueryGet(uri, params, secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error sending get request, uri: %s, params: %s", uri, params)
		return
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if !ok {
		err = apperrors.New("Can not parse order list query response, queryResp: %s", queryResp)
		return
	}

	if queryRespMap["list"] == nil {
		return
	}

	queryRespData, ok := queryRespMap["list"].([]interface{})
	if !ok {
		err = apperrors.New("Can not parse order list data values, queryRespMap: %s", queryRespMap)
		return
	}

	for _, queryRespOrderItem := range queryRespData {
		var order *OrderResponse
		order, err = mapOrderResponse(queryRespOrderItem)
		if err != nil {
			err = apperrors.Wrap(err, "error map order response, queryRespOrderItem: %s", queryRespOrderItem)
			return
		}
		orders = append(orders, order)
	}

	return
}

func mapOrderResponse(queryResp interface{}) (orderResponse *OrderResponse, err error) {

	var orderResponseBytes []byte
	orderResponseBytes, err = json.Marshal(queryResp)
	if err != nil {
		err = apperrors.Wrap(err, "can not marshal order response for ByBit, queryResp: %s", queryResp)
		return
	}

	response := OrderResponse{}
	err = json.Unmarshal(orderResponseBytes, &response)
	if err != nil {
		err = apperrors.Wrap(err, "can not unmarshal order response for ByBit, orderResponseBytes: %s", orderResponseBytes)
		return
	}
	orderResponse = &response

	return
}
