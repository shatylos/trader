package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"strconv"
)

type Position struct {
	Quantity               float64 `json:"quantity"`
	Leverage               int64   `json:"leverage"`
	MarginLevel            float64 `json:"margin_level"`
	IsMarginCall           bool    `json:"is_margin_call"`
	FundingNextPayDt       int64   `json:"funding_next_pay_dt"`
	Pair                   string  `json:"pair"`
	Type                   string  `json:"type"`
	BasePrice              float64 `json:"base_price"`
	Margin                 float64 `json:"margin"`
	RealizedPnl            float64 `json:"realized_pnl"`
	MarginCallPrice        float64 `json:"margin_call_price"`
	FundingQuantity        float64 `json:"funding_quantity"`
	FundingRate            float64 `json:"funding_rate"`
	Created                int64   `json:"created"`
	TakeProfit             float64 `json:"take_profit"`
	StopLoss               float64 `json:"stop_loss"`
	MarginCurrency         string  `json:"margin_currency"`
	Roe                    float64 `json:"roe"`
	FundingNextPayQuantity float64 `json:"funding_next_pay_quantity"`
	PositionId             int64   `json:"position_id"`
	UnrealizedPnl          float64 `json:"unrealized_pnl"`
	PnlCurrency            string  `json:"pnl_currency"`
	BreakEvenPrice         float64 `json:"break_even_price"`
	LiquidationPrice       float64 `json:"liquidation_price"`
	NeedLiquidate          bool    `json:"need_liquidate"`
	FundingCurrency        string  `json:"funding_currency"`
}

func GetMarginPositionList() ([]Position, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQuery("/margin/user/position/list", params)
	if er != nil {
		return nil, er
	}
	return mapPositionList(queryResp)
}

func mapPositionList(queryResp map[string]interface{}) ([]Position, error) {

	positions, ok := queryResp["positions"].([]interface{})
	if !ok {
		return nil, utils.AppError{Message: "[Exmo Position] Can not parse broker response. Bad positions value."}
	}

	positionsResponse := make([]Position, len(positions))

	for i, positionInterface := range positions {
		position, ok := positionInterface.(map[string]interface{})
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Position] Can not parse broker response. Bad position value."}
		}

		quantity, err := posParseFloat(position, "quantity")
		if err != nil {
			return nil, err
		}

		marginLevel, err := posParseFloat(position, "quantity")
		if err != nil {
			return nil, err
		}
		basePrice, err := posParseFloat(position, "base_price")
		if err != nil {
			return nil, err
		}
		margin, err := posParseFloat(position, "margin")
		if err != nil {
			return nil, err
		}
		realizedPnl, err := posParseFloat(position, "realized_pnl")
		if err != nil {
			return nil, err
		}
		marginCallPrice, err := posParseFloat(position, "margin_call_price")
		if err != nil {
			return nil, err
		}
		fundingQuantity, err := posParseFloat(position, "funding_quantity")
		if err != nil {
			return nil, err
		}
		fundingRate, err := posParseFloat(position, "funding_rate")
		if err != nil {
			return nil, err
		}
		takeProfit, err := posParseFloat(position, "take_profit")
		if err != nil {
			return nil, err
		}
		stopLoss, err := posParseFloat(position, "stop_loss")
		if err != nil {
			return nil, err
		}
		roe, err := posParseFloat(position, "roe")
		if err != nil {
			return nil, err
		}
		fundingNextPayQuantity, err := posParseFloat(position, "funding_next_pay_quantity")
		if err != nil {
			return nil, err
		}
		unrealizedPnl, err := posParseFloat(position, "unrealized_pnl")
		if err != nil {
			return nil, err
		}
		breakEvenPrice, err := posParseFloat(position, "break_even_price")
		if err != nil {
			return nil, err
		}
		liquidationPrice, err := posParseFloat(position, "liquidation_price")
		if err != nil {
			return nil, err
		}
		positionId, err := posParseInt(position, "position_id")
		if err != nil {
			return nil, err
		}
		leverage, err := posParseInt(position, "leverage")
		if err != nil {
			return nil, err
		}
		fundingNextPayDt, err := posParseInt(position, "funding_next_pay_dt")
		if err != nil {
			return nil, err
		}
		created, err := posParseInt(position, "created")
		if err != nil {
			return nil, err
		}

		pair, err := posParseStr(position, "pair")
		if err != nil {
			return nil, err
		}
		_type, err := posParseStr(position, "type")
		if err != nil {
			return nil, err
		}
		marginCurrency, err := posParseStr(position, "margin_currency")
		if err != nil {
			return nil, err
		}
		pnlCurrency, err := posParseStr(position, "pnl_currency")
		if err != nil {
			return nil, err
		}
		fundingCurrency, err := posParseStr(position, "funding_currency")
		if err != nil {
			return nil, err
		}
		isMarginCall, err := posParseBool(position, "is_margin_call")
		if err != nil {
			return nil, err
		}
		needLiquidate, err := posParseBool(position, "need_liquidate")
		if err != nil {
			return nil, err
		}

		positionsResponse[i] = Position{
			Quantity:               quantity,
			Leverage:               leverage,
			MarginLevel:            marginLevel,
			IsMarginCall:           isMarginCall,
			FundingNextPayDt:       fundingNextPayDt,
			Pair:                   pair,
			Type:                   _type,
			BasePrice:              basePrice,
			Margin:                 margin,
			RealizedPnl:            realizedPnl,
			MarginCallPrice:        marginCallPrice,
			FundingQuantity:        fundingQuantity,
			FundingRate:            fundingRate,
			Created:                created,
			TakeProfit:             takeProfit,
			StopLoss:               stopLoss,
			MarginCurrency:         marginCurrency,
			Roe:                    roe,
			FundingNextPayQuantity: fundingNextPayQuantity,
			PositionId:             positionId,
			UnrealizedPnl:          unrealizedPnl,
			PnlCurrency:            pnlCurrency,
			BreakEvenPrice:         breakEvenPrice,
			LiquidationPrice:       liquidationPrice,
			NeedLiquidate:          needLiquidate,
			FundingCurrency:        fundingCurrency,
		}

	}

	return positionsResponse, nil
}

func posParseFloat(position map[string]interface{}, key string) (float64, error) {
	resultStr, ok := position[key].(string)
	if !ok {
		return 0, utils.AppError{Message: fmt.Sprintf("[Exmo Position] Can not parse broker response. Bad \"%s\" value.", key)}
	}
	return strconv.ParseFloat(resultStr, 64)
}

func posParseInt(position map[string]interface{}, key string) (int64, error) {
	resultStr, ok := position[key].(string)
	if !ok {
		return 0, utils.AppError{Message: fmt.Sprintf("[Exmo Position] Can not parse broker response. Bad \"%s\" value.", key)}
	}
	return strconv.ParseInt(resultStr, 10, 64)
}

func posParseStr(position map[string]interface{}, key string) (string, error) {
	resultStr, ok := position[key].(string)
	if !ok {
		return "", utils.AppError{Message: fmt.Sprintf("[Exmo Position] Can not parse broker response. Bad \"%s\" value.", key)}
	}
	return resultStr, nil
}

func posParseBool(position map[string]interface{}, key string) (bool, error) {
	resultBool, ok := position[key].(bool)
	if !ok {
		return false, utils.AppError{Message: fmt.Sprintf("[Exmo Position] Can not parse broker response. Bad \"%s\" value.", key)}
	}
	return resultBool, nil
}
