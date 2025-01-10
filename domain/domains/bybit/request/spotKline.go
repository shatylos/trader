package request

import (
	bybitStructs "github.com/shatylos/trader/domain/domains/bybit/structs"
	"github.com/shatylos/trader/utils"
	"strconv"
)

var SpotMinToResol = map[string]string{
	//"1":   60,
	//"5":   300,
	//"15":  900,
	"30": "30",
	//"45":  2700,
	//"60":  3600,
	//"120": 7200,
	//"180": 10800,
	//"240": 14400,
	//"D":   86400,
	//"W":   604800,
	//"M":   2592000,
}

type SpotCandle struct {
	Timestamp int64  `json:"t"` //	Timestamp(ms)
	Close     string `json:"c"` //	Close price
	High      string `json:"h"` //	High price
	Low       string `json:"l"` //	Low price
	Open      string `json:"o"` //	Open price
	Volume    string `json:"v"` //	Trading volume
}

func GetSpotKlineList(symbol string, resolution string, from int64, limit int64, secrets bybitStructs.Secrets) ([]SpotCandle, error) {
	params := make(ApiParams, 0)
	params["category"] = "spot"
	params["symbol"] = symbol
	params["interval"] = SpotMinToResol[resolution]
	params["start"] = strconv.FormatInt(from*1000, 10)
	params["limit"] = strconv.FormatInt(limit, 10)

	queryResp, er := apiQueryGet("/v5/market/kline", params, secrets)
	if er != nil {
		return nil, er
	}

	queryRespMap, ok := queryResp.(map[string]interface{})
	if ok == false {
		return nil, utils.AppError{
			Message: "Error parse public kline list response for ByBit",
		}
	}

	if queryRespMap["list"] == nil {
		return nil, utils.AppError{
			Message: "Empty list response in Kline from ByBit",
		}
	}

	return mapSpotCandleHistory(queryRespMap["list"])
}

func mapSpotCandleHistory(source interface{}) ([]SpotCandle, error) {

	sourceSlice, ok := source.([]interface{})
	if ok == false {
		return nil, utils.AppError{
			Message: "Error parse public kline list response for ByBit",
		}
	}

	res := make([]SpotCandle, len(sourceSlice))

	for i, sourceCandle := range sourceSlice {
		sourceCandleArr, ok := sourceCandle.([]interface{})
		if ok == false {
			return nil, utils.AppError{
				Message: "Error parse candle in public kline list response for ByBit",
			}
		}
		if len(sourceCandleArr) < 6 {
			return nil, utils.AppError{
				Message: "Unexpected candle format in public kline list response for ByBit",
			}
		}
		resCandle := SpotCandle{}
		ts, err := utils.ToInt64(sourceCandleArr[0])
		if err != nil {
			return nil, utils.AppError{
				Message: "Can not convert timestamp to int64 in public kline list response for ByBit",
			}
		}
		resCandle.Timestamp = ts / 1000
		resCandle.Open = sourceCandleArr[1].(string)
		resCandle.High = sourceCandleArr[2].(string)
		resCandle.Low = sourceCandleArr[3].(string)
		resCandle.Close = sourceCandleArr[4].(string)
		resCandle.Volume = sourceCandleArr[5].(string)
		res[i] = resCandle
	}

	return res, nil
}
