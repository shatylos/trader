package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"strconv"
)

type Candle struct {
	Close    float64 `json:"close"`
	High     float64 `json:"high"`
	Id       float64 `json:"id"`
	Interval string  `json:"interval"`
	Low      float64 `json:"low"`
	Open     float64 `json:"open"`
	OpenTime float64 `json:"open_time"`
	Period   string  `json:"period"`
	StartAt  float64 `json:"start_at"`
	Symbol   string  `json:"symbol"`
	Turnover float64 `json:"turnover"`
	Volume   float64 `json:"volume"`
}

func GetKlineList(symbol string, resolution string, from int64, limit int64, isDemo bool) ([]Candle, error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["interval"] = resolution
	params["from"] = strconv.FormatInt(from, 10)
	params["limit"] = strconv.FormatInt(limit, 10)

	queryResp, er := apiQueryGet("/public/linear/kline", params, isDemo)
	if er != nil {
		return nil, er
	}

	return mapCandleHistory(queryResp)
}

func mapCandleHistory(source interface{}) ([]Candle, error) {

	sourceSlice, ok := source.([]interface{})
	if ok == false {
		return nil, utils.AppError{
			Message: "Error parse public kline list response for ByBit",
		}
	}

	res := make([]Candle, len(sourceSlice))

	for i, sourceCandle := range sourceSlice {
		sourceCandleBytes, err := json.Marshal(sourceCandle)
		if err != nil {
			return nil, utils.AppError{
				Message: "Can not Marshal public kline list response for ByBit",
			}
		}
		resCandle := Candle{}
		err = json.Unmarshal(sourceCandleBytes, &resCandle)
		if err != nil {
			return nil, utils.AppError{
				Message: "Can not Unmarshal public kline list response for ByBit",
			}
		}
		res[i] = resCandle
	}

	return res, nil
}
