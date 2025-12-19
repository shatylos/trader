package request

import (
	"errors"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools/apperrors"
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

	uri := "/v5/position/set-leverage"
	_, err = apiQueryPost(uri, params, secrets)
	if err != nil {
		if errors.Is(err, LeverageNotModifiedApiError) {
			err = nil
		} else {
			err = apperrors.Wrap(err, "error send api post query, uri: %s, params: %s", uri, params)
			return
		}
	}

	return
}
