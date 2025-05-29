package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
	_type "github.com/shatylos/trader/tools/type"
)

type MarginWalletBalance struct {
	TotalEquity           float64 `json:"totalEquity"`
	TotalMarginBalance    float64 `json:"totalMarginBalance"`
	TotalAvailableBalance float64 `json:"totalAvailableBalance"`
	TotalInitialMargin    float64 `json:"totalInitialMargin"`
	TotalWalletBalance    float64 `json:"totalWalletBalance"`
}

func GetFuturesWalletBalance(secrets bybitStructs.Secrets) (response MarginWalletBalance, err error) {
	params := make(ApiParams, 0)
	params["accountType"] = "UNIFIED"
	queryResp, er := apiQueryGet("/v5/account/wallet-balance", params, secrets)
	if er != nil {
		return
	}
	response, err = mapMarginWalletBalance(queryResp)
	return
}

func mapMarginWalletBalance(source interface{}) (response MarginWalletBalance, err error) {
	sourceMap, ok := source.(map[string]interface{})
	if !ok {
		err = tools.AppError{
			Message: "Error marshalling source.",
		}
		return
	}

	sourceList, ok := sourceMap["list"].([]interface{})
	if !ok {
		err = tools.AppError{
			Message: "Error marshalling source.",
		}
		return
	}

	var sourceBytes []byte
	sourceBytes, err = json.Marshal(sourceList[0])

	if err != nil {
		err = tools.AppError{
			Message: fmt.Sprintf("Error marshalling source: %s", err.Error()),
		}
		return
	}

	type MarginWalletBalanceRaw struct {
		TotalEquity           string `json:"totalEquity"`
		TotalMarginBalance    string `json:"totalMarginBalance"`
		TotalAvailableBalance string `json:"totalAvailableBalance"`
		TotalInitialMargin    string `json:"totalInitialMargin"`
		TotalWalletBalance    string `json:"totalWalletBalance"`
	}
	var rawResponse MarginWalletBalanceRaw
	err = json.Unmarshal(sourceBytes, &rawResponse)
	if err != nil {
		err = tools.AppError{
			Message: fmt.Sprintf("Error unmarshalling source: %s", err.Error()),
		}
		return
	}

	var totalEquity, totalMarginBalance, totalAvailableBalance, totalInitialMargin, totalWalletBalance float64

	totalEquity, err = _type.ToFloat64(rawResponse.TotalEquity)
	if err != nil {
		return
	}
	totalMarginBalance, err = _type.ToFloat64(rawResponse.TotalMarginBalance)
	if err != nil {
		return
	}
	totalAvailableBalance, err = _type.ToFloat64(rawResponse.TotalAvailableBalance)
	if err != nil {
		return
	}
	totalInitialMargin, err = _type.ToFloat64(rawResponse.TotalInitialMargin)
	if err != nil {
		return
	}
	totalWalletBalance, err = _type.ToFloat64(rawResponse.TotalWalletBalance)
	if err != nil {
		return
	}

	response = MarginWalletBalance{
		TotalEquity:           totalEquity,
		TotalMarginBalance:    totalMarginBalance,
		TotalAvailableBalance: totalAvailableBalance,
		TotalInitialMargin:    totalInitialMargin,
		TotalWalletBalance:    totalWalletBalance,
	}

	return
}
