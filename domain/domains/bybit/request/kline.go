package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"strconv"
)

type Candle struct {
	Symbol   string  `json:"symbol"`
	Interval string  `json:"interval"`
	Open     string  `json:"open"`
	High     string  `json:"high"`
	Low      string  `json:"low"`
	OpenTime float64 `json:"open_time"`
	Close    string  `json:"close"`
	Volume   string  `json:"volume"`
	Turnover string  `json:"turnover"`
}

func GetKlineList(symbol string, resolution string, from int64, limit int64, isDemo bool) ([]Candle, error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["interval"] = resolution
	params["from"] = strconv.FormatInt(from, 10)
	params["limit"] = strconv.FormatInt(limit, 10)

	queryResp, er := apiQueryGet("/v2/public/kline/list", params, isDemo)
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
