package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"strconv"
)

var SpotMinToResol = map[string]string{
	//"1":   60,
	//"5":   300,
	//"15":  900,
	"30": "30m",
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
	Timestamp int64  `json:"t"`  //	Timestamp(ms)
	Symbol    string `json:"s"`  //	Name of the trading pair
	Alias     string `json:"sn"` //	Alias
	Close     string `json:"c"`  //	Close price
	High      string `json:"h"`  //	High price
	Low       string `json:"l"`  //	Low price
	Open      string `json:"o"`  //	Open price
	Volume    string `json:"v"`  //	Trading volume
}

func GetSpotKlineList(symbol string, resolution string, from int64, limit int64, isDemo bool) ([]SpotCandle, error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["interval"] = SpotMinToResol[resolution]
	params["startTime"] = strconv.FormatInt(from*1000, 10)
	params["limit"] = strconv.FormatInt(limit, 10)

	queryResp, er := apiQueryGet("/spot/v3/public/quote/kline", params, isDemo)
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
		sourceCandleBytes, err := json.Marshal(sourceCandle)
		if err != nil {
			return nil, utils.AppError{
				Message: "Can not Marshal public kline list response for ByBit",
			}
		}
		resCandle := SpotCandle{}
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
