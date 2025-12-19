package request

import (
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"strconv"
)

type Candle struct {
	StartTime  string `json:"0"`
	OpenPrice  string `json:"1"`
	HighPrice  string `json:"2"`
	LowPrice   string `json:"3"`
	ClosePrice string `json:"4"`
	Volume     string `json:"5"`
	Turnover   string `json:"6"`
}

func GetKlineList(symbol string, resolution string, limit int64, secrets bybitStructs.Secrets) (candles []Candle, err error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["interval"] = resolution
	params["limit"] = strconv.FormatInt(limit, 10)

	var queryResp interface{}
	uri := "/v5/market/kline"
	queryResp, err = apiQueryGet(uri, params, secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error Bybit api query get %s", uri)
		return
	}

	candles, err = mapCandleHistory(queryResp)
	if err != nil {
		err = apperrors.Wrap(err, "error mapping candle history, queryResp: %s", queryResp)
		return
	}
	return
}

func mapCandleHistory(source interface{}) (res []Candle, err error) {

	sourceMap, ok := source.(map[string]interface{})
	if ok == false {
		err = apperrors.New("error parse public kline response for ByBit")
		return
	}

	listSlice, ok := sourceMap["list"].([]interface{})
	if ok == false {
		err = apperrors.New("error parsing public kline list response for ByBit")
		return
	}

	res = make([]Candle, len(listSlice))
	for i, sourceCandle := range listSlice {
		sourceCandleArr, ok := sourceCandle.([]interface{})
		if !ok {
			err = apperrors.New("error parsing public kline item for ByBit")
			return
		}
		if len(sourceCandleArr) != 7 {
			err = apperrors.New("count of items in kline is %d. Expected 7 items", len(sourceCandleArr))
			return
		}

		startTime, ok := sourceCandleArr[0].(string)
		if !ok {
			err = apperrors.New("kline item 0 is not a string")
			return
		}
		openPrice, ok := sourceCandleArr[1].(string)
		if !ok {
			err = apperrors.New("kline item 1 is not a string")
			return
		}
		highPrice, ok := sourceCandleArr[2].(string)
		if !ok {
			err = apperrors.New("kline item 2 is not a string")
			return
		}
		lowPrice, ok := sourceCandleArr[3].(string)
		if !ok {
			err = apperrors.New("kline item 3 is not a string")
			return
		}
		closePrice, ok := sourceCandleArr[4].(string)
		if !ok {
			err = apperrors.New("kline item 4 is not a string")
			return
		}
		volume, ok := sourceCandleArr[5].(string)
		if !ok {
			err = apperrors.New("kline item 5 is not a string")
			return
		}
		turnover, ok := sourceCandleArr[6].(string)
		if !ok {
			err = apperrors.New("kline item 6 is not a string")
			return
		}

		resCandle := Candle{
			StartTime:  startTime,
			OpenPrice:  openPrice,
			HighPrice:  highPrice,
			LowPrice:   lowPrice,
			ClosePrice: closePrice,
			Volume:     volume,
			Turnover:   turnover,
		}

		res[i] = resCandle
	}

	return
}
