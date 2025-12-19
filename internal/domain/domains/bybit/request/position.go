package request

import (
	"encoding/json"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools/apperrors"
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
	uri := "/v5/position/list"
	queryResp, err = apiQueryGet(uri, params, secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error sending get request, uri: %s, params: %s", uri, params)
		return
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if ok == false {
		err = apperrors.Wrap(err, "can not parse broker response, queryResp: %s", queryResp)
		return
	}

	position, err = extractPosition(coinPare, queryRespMap["list"])
	if err != nil {
		err = apperrors.Wrap(err, "error extract position, coinPare: %s, queryRespMap list: %s", coinPare, queryRespMap["list"])
		return
	}
	return
}

func extractPosition(coinPare string, source interface{}) (resultPosition Position, err error) {

	sourceSlice, ok := source.([]interface{})
	if ok == false {
		err = apperrors.Wrap(err, "can not parse broker response. Expected slice of positions. Source: %s", source)
		return
	}

	for _, sourcePosition := range sourceSlice {
		var positionBytes []byte
		positionBytes, err = json.Marshal(sourcePosition)
		if err != nil {
			err = apperrors.Wrap(err, "error marshal source position: %s", sourcePosition)
			return
		}
		position := Position{}
		err = json.Unmarshal(positionBytes, &position)
		if err != nil {
			err = apperrors.Wrap(err, "error unmarshal position: %s", positionBytes)
			return
		}
		if position.Symbol == coinPare {
			resultPosition = position
			return
		}
	}

	err = apperrors.New("position not found for the symbol %s", coinPare)
	return
}
