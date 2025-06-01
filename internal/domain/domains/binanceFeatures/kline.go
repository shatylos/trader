package binanceFeatures

import (
	"encoding/json"
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/mapping"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/request"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	providerStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"sort"
	"strconv"
)

func (d *DomainBinanceFutures) LoadCandleHistory(symbol string, resolution string, limit int64) (candles []providerStructs.DomainCandle, err error) {
	var interval string
	interval, err = mapping.ToBinanceInterval(resolution)
	if err != nil {
		return
	}
	apiRequest := request.ApiGetRequest{
		Uri: "/fapi/v1/klines",
		ApiParams: binanceStructs.ApiParams{
			"symbol":   symbol,
			"interval": interval,
			"limit":    strconv.FormatInt(limit, 10),
		},
		Secrets: d.secrets,
	}
	var rawRespnse binanceStructs.ApiResponse
	rawRespnse, err = apiRequest.DoRequest()
	candles, err = unmarshalCandleHistory(rawRespnse)

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Time > candles[j].Time
	})

	return
}

func unmarshalCandleHistory(rawResponse binanceStructs.ApiResponse) (candles []providerStructs.DomainCandle, err error) {

	var rawCandles [][]interface{}
	err = json.Unmarshal(rawResponse, &rawCandles)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshaling raw response: %s", err)
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}

	for _, c := range rawCandles {
		var (
			openTime                          int64
			open, high, low, closeVal, volume float64
		)

		openTime, err = _type.ToInt64(c[0])
		if err != nil {
			return
		}
		openTime = openTime / 1000
		open, err = _type.ToFloat64(c[1])
		if err != nil {
			return
		}
		high, err = _type.ToFloat64(c[2])
		if err != nil {
			return
		}
		low, err = _type.ToFloat64(c[3])
		if err != nil {
			return
		}
		closeVal, err = _type.ToFloat64(c[4])
		if err != nil {
			return
		}
		volume, err = _type.ToFloat64(c[5])
		if err != nil {
			return
		}

		candle := providerStructs.DomainCandle{
			Time:   openTime,
			High:   high,
			Low:    low,
			Open:   open,
			Close:  closeVal,
			Volume: volume,
		}
		candles = append(candles, candle)
	}

	return
}
