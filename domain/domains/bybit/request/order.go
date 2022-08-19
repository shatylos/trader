package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"fmt"
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
}

type OrderResponse struct {
	UserId         float64 `json:"user_id"`          // -> 646256
	Symbol         string  `json:"symbol"`           // -> BTCUSDT
	Side           string  `json:"side"`             // -> Buy
	Price          float64 `json:"price"`            // -> 24553
	CumExecValue   float64 `json:"cum_exec_value"`   // -> 0
	StopLoss       float64 `json:"stop_loss"`        // -> 22800
	OrderType      string  `json:"order_type"`       // -> Market
	TimeInForce    string  `json:"time_in_force"`    // -> ImmediateOrCancel
	OrderStatus    string  `json:"order_status"`     // -> Created
	ReduceOnly     bool    `json:"reduce_only"`      // -> false
	CloseOnTrigger bool    `json:"close_on_trigger"` // -> false
	CreatedTime    string  `json:"created_time"`     // -> 2022-08-17T20:45:45Z
	TakeProfit     float64 `json:"take_profit"`      // -> 23600
	SlTriggerBy    string  `json:"sl_trigger_by"`    // -> LastPrice
	PositionIdx    float64 `json:"position_idx"`     // -> 1
	LastExecPrice  float64 `json:"last_exec_price"`  // -> 0
	CumExecQty     float64 `json:"cum_exec_qty"`     // -> 0
	TpTriggerBy    string  `json:"tp_trigger_by"`    // -> LastPrice
	OrderId        string  `json:"order_id"`         // -> 10772450-26ce-4c92-8bf6-5cf29d190a50
	Qty            float64 `json:"qty"`              // -> 0.005
	CumExecFee     float64 `json:"cum_exec_fee"`     // -> 0
	OrderLinkId    string  `json:"order_link_id"`    // -> 123qwe
	UpdatedTime    string  `json:"updated_time"`     // -> 2022-08-17T20:45:45Z
}

func CreateOrder(orderRequest OrderRequest, isDemo bool) (*OrderResponse, error) {
	params := make(ApiParams, 0)
	params["time_in_force"] = orderRequest.TimeInForce
	params["reduce_only"] = orderRequest.ReduceOnly
	params["close_on_trigger"] = orderRequest.CloseOnTrigger

	params["side"] = orderRequest.Side
	params["symbol"] = orderRequest.Symbol
	params["order_type"] = orderRequest.OrderType
	params["qty"] = fmt.Sprintf("%f", orderRequest.Qty)

	if orderRequest.Price > 0 {
		params["price"] = fmt.Sprintf("%f", orderRequest.Price)
	}
	if orderRequest.TakeProfit > 0 {
		params["take_profit"] = fmt.Sprintf("%f", orderRequest.TakeProfit)
	}
	if orderRequest.StopLoss > 0 {
		params["stop_loss"] = fmt.Sprintf("%f", orderRequest.StopLoss)
	}
	if orderRequest.OrderLinkId != "" {
		params["order_link_id"] = orderRequest.OrderLinkId
	}

	queryResp, err := apiQueryPost("/private/linear/order/create", params, isDemo)
	if err != nil {
		return nil, err
	}

	return mapOrderResponse(queryResp)
}

func mapOrderResponse(queryResp interface{}) (*OrderResponse, error) {

	orderResponseBytes, err := json.Marshal(queryResp)
	if err != nil {
		return nil, utils.AppError{
			Message: "Can not Marshal order response for ByBit",
		}
	}

	orderResponse := OrderResponse{}
	err = json.Unmarshal(orderResponseBytes, &orderResponse)
	if err != nil {
		return nil, utils.AppError{
			Message: "Can not Unmarshal order response for ByBit",
		}
	}

	return &orderResponse, nil
}
