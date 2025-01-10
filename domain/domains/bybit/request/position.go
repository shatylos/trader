package request

import (
	"encoding/json"
	bybitStructs "github.com/shatylos/trader/domain/domains/bybit/structs"
	"github.com/shatylos/trader/utils"
)

type Position struct {
	BustPrice           float64 `json:"bust_price"`
	CumRealisedPnl      float64 `json:"cum_realised_pnl"`
	FreeQty             float64 `json:"free_qty"`
	TrailingStop        float64 `json:"trailing_stop"`
	UserId              float64 `json:"user_id"`
	Size                float64 `json:"size"`
	PositionValue       float64 `json:"position_value"`
	EntryPrice          float64 `json:"entry_price"`
	OccClosingFee       float64 `json:"occ_closing_fee"`
	UnrealisedPnl       float64 `json:"unrealised_pnl"`
	StopLoss            float64 `json:"stop_loss"`
	TakeProfit          float64 `json:"take_profit"`
	PositionIdx         float64 `json:"position_idx"`
	Mode                string  `json:"mode"`
	Symbol              string  `json:"symbol"`
	Leverage            float64 `json:"leverage"`
	AutoAddMargin       float64 `json:"auto_add_margin"`
	IsIsolated          bool    `json:"is_isolated"`
	DeleverageIndicator float64 `json:"deleverage_indicator"`
	Side                string  `json:"side"`
	LiqPrice            float64 `json:"liq_price"`
	PositionMargin      float64 `json:"position_margin"`
	RealisedPnl         float64 `json:"realised_pnl"`
	TpSlMode            string  `json:"tp_sl_mode"`
	RiskId              float64 `json:"risk_id"`
}

func GetPositionList(coinPare string, secrets bybitStructs.Secrets) ([]Position, error) {
	params := make(ApiParams, 0)
	params["symbol"] = coinPare

	queryResp, er := apiQueryGet("/private/linear/position/list", params, secrets)
	if er != nil {
		return nil, er
	}
	return mapPositionList(queryResp)
}

func mapPositionList(source interface{}) ([]Position, error) {
	resultPositions := make([]Position, 0)

	sourceSlice, ok := source.([]interface{})
	if ok == false {
		return nil, utils.AppError{Message: "[Bybit Position List] Can not parse broker response. Expected slice of positions."}
	}

	for _, sourcePosition := range sourceSlice {
		positionBytes, err := json.Marshal(sourcePosition)
		if err != nil {
			return nil, err
		}
		position := Position{}
		err = json.Unmarshal(positionBytes, &position)
		if err != nil {
			return nil, err
		}
		if position.Size == 0 {
			continue
		}
		resultPositions = append(resultPositions, position)
	}

	return resultPositions, nil
}
