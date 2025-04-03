package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
)

type Position struct {
	AvgPrice       string `json:"avgPrice"`
	CreatedTime    string `json:"createdTime"`
	CumRealisedPnl string `json:"cumRealisedPnl"`
	CurRealisedPnl string `json:"curRealisedPnl"`
	Leverage       string `json:"leverage"`
	MarkPrice      string `json:"markPrice"`
	PositionIM     string `json:"positionIM"`
	PositionMM     string `json:"positionMM"`
	PositionValue  string `json:"positionValue"`
	Side           string `json:"side"`
	Size           string `json:"size"`
	StopLoss       string `json:"stopLoss"`
	Symbol         string `json:"symbol"`
	TakeProfit     string `json:"takeProfit"`
	UnrealizedPnl  string `json:"unrealisedPnl"`
	UpdatedTime    string `json:"updatedTime"`
}

func GetPosition(coinPare string, secrets bybitStructs.Secrets) (position Position, err error) {
	params := make(ApiParams, 0)
	params["category"] = "linear"
	params["symbol"] = coinPare

	var queryResp interface{}
	queryResp, err = apiQueryGet("/v5/position/list", params, secrets)
	if err != nil {
		return
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if ok == false {
		return position, tools.AppError{Message: "[Bybit Position List] Can not parse broker response."}
	}

	return extractPosition(coinPare, queryRespMap["list"])
}

func extractPosition(coinPare string, source interface{}) (resultPosition Position, err error) {

	sourceSlice, ok := source.([]interface{})
	if ok == false {
		return resultPosition, tools.AppError{Message: "[Bybit Position] Can not parse broker response. Expected slice of positions."}
	}

	for _, sourcePosition := range sourceSlice {
		var positionBytes []byte
		positionBytes, err = json.Marshal(sourcePosition)
		if err != nil {
			return
		}
		position := Position{}
		err = json.Unmarshal(positionBytes, &position)
		if err != nil {
			return
		}
		if position.Symbol == coinPare {
			return position, nil
		}
	}

	return resultPosition, tools.AppError{Message: fmt.Sprintf("[Bybit Position] Position not found for the symbol (%s).", coinPare)}
}
