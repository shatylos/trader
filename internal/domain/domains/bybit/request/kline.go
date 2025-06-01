package request

import (
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
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

func GetKlineList(symbol string, resolution string, limit int64, secrets bybitStructs.Secrets) ([]Candle, error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["interval"] = resolution
	params["limit"] = strconv.FormatInt(limit, 10)

	queryResp, er := apiQueryGet("/v5/market/kline", params, secrets)
	if er != nil {
		return nil, er
	}

	return mapCandleHistory(queryResp)
}

func mapCandleHistory(source interface{}) (res []Candle, err error) {

	sourceMap, ok := source.(map[string]interface{})
	if ok == false {
		msg := "Error parse public kline response for ByBit"
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	listSlice, ok := sourceMap["list"].([]interface{})
	if ok == false {
		msg := "Error parsing public kline list response for ByBit"
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	res = make([]Candle, len(listSlice))
	for i, sourceCandle := range listSlice {
		sourceCandleArr, ok := sourceCandle.([]interface{})
		if !ok {
			msg := "Error parsing public kline item for ByBit"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		if len(sourceCandleArr) != 7 {
			msg := fmt.Sprintf("Count of items in kline is %d. Expected 7 items", len(sourceCandleArr))
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}

		startTime, ok := sourceCandleArr[0].(string)
		if !ok {
			msg := "kline item 0 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		openPrice, ok := sourceCandleArr[1].(string)
		if !ok {
			msg := "kline item 1 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		highPrice, ok := sourceCandleArr[2].(string)
		if !ok {
			msg := "kline item 2 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		lowPrice, ok := sourceCandleArr[3].(string)
		if !ok {
			msg := "kline item 3 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		closePrice, ok := sourceCandleArr[4].(string)
		if !ok {
			msg := "kline item 4 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		volume, ok := sourceCandleArr[5].(string)
		if !ok {
			msg := "kline item 5 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
			return
		}
		turnover, ok := sourceCandleArr[6].(string)
		if !ok {
			msg := "kline item 6 is not a string"
			logger.Error(msg)
			err = tools.AppError{
				Message: msg,
			}
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
