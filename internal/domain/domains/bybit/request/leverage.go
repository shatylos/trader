package request

import (
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
	"strconv"
)

type LeverageRequest struct {
	Symbol       string
	BuyLeverage  int64
	SellLeverage int64
}

func SetLeverage(leverageRequest LeverageRequest, secrets bybitStructs.Secrets) (err error) {
	params := make(ApiParams)
	params["category"] = "linear"
	params["symbol"] = leverageRequest.Symbol
	params["buyLeverage"] = strconv.FormatInt(leverageRequest.BuyLeverage, 10)
	params["sellLeverage"] = strconv.FormatInt(leverageRequest.SellLeverage, 10)

	_, err = apiQueryPost("/v5/position/set-leverage", params, secrets)
	if err != nil {
		err2, ok := err.(tools.AppError)
		// 110043 - set leverage has not been modified
		if ok && err2.Code == 110043 {
			err = nil
		}
	}

	return
}
