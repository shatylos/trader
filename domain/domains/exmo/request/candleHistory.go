package request

import (
	"github.com/shatylos/trader/utils"
	"strconv"
)

type Candle struct {
	T int64
	O float64
	C float64
	H float64
	L float64
	V float64
}

func LoadCandleHistory(symbol string, resolution string, from int64, to int64) ([]Candle, error) {
	params := make(ApiParams, 0)
	params["symbol"] = symbol
	params["resolution"] = resolution
	params["from"] = strconv.FormatInt(from, 10)
	params["to"] = strconv.FormatInt(to, 10)

	queryResp, er := apiQuery("candles_history", params)
	if er != nil {
		return nil, er
	}

	return mapCandleHistory(queryResp)
}

func mapCandleHistory(queryResp map[string]interface{}) ([]Candle, error) {

	candles, ok := queryResp["candles"].([]interface{})
	if !ok {
		return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candles value."}
	}

	candleResult := make([]Candle, len(candles))

	for i, candleInterface := range candles {
		candle, ok := candleInterface.(map[string]interface{})
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle value."}
		}

		t, ok := candle["t"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"t\" value."}
		}

		o, ok := candle["o"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"o\" value."}
		}

		c, ok := candle["c"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"c\" value."}
		}

		h, ok := candle["h"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"h\" value."}
		}

		l, ok := candle["l"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"l\" value."}
		}

		v, ok := candle["v"].(float64)
		if !ok {
			return nil, utils.AppError{Message: "[Exmo Candle History] Can not parse broker response. Bad candle \"v\" value."}
		}

		candleResult[i] = Candle{
			T: int64(t),
			O: o,
			C: c,
			H: h,
			L: l,
			V: v,
		}
	}

	return candleResult, nil
}
