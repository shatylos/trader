package request

import (
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
)

type WalletBalance struct {
	Equity           float64 `json:"equity"`
	AvailableBalance float64 `json:"available_balance"`
	UsedMargin       float64 `json:"used_margin"`
	OrderMargin      float64 `json:"order_margin"`
	OccClosingFee    float64 `json:"occ_closing_fee"`
	WalletBalance    float64 `json:"wallet_balance"`
	UnrealisedPnl    float64 `json:"unrealised_pnl"`
	CumRealisedPnl   float64 `json:"cum_realised_pnl"`
	ServiceCash      float64 `json:"service_cash"`
	PositionMargin   float64 `json:"position_margin"`
	OccFundingFee    float64 `json:"occ_funding_fee"`
	RealisedPnl      float64 `json:"realised_pnl"`
	GivenCash        float64 `json:"given_cash"`
}

func GetWalletBalance(secrets bybitStructs.Secrets) (*map[string]WalletBalance, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQueryGet("/v2/private/wallet/balance", params, secrets)
	if er != nil {
		return nil, er
	}
	return mapWalletBalance(queryResp)
}

func mapWalletBalance(source interface{}) (*map[string]WalletBalance, error) {
	result := map[string]WalletBalance{}

	sourceMap, ok := source.(map[string]interface{})
	if ok == false {
		return nil, utils.AppError{
			Message: "Error parse WalletBalance response for ByBit",
		}
	}

	for coin, coinValues := range sourceMap {
		coinValuesBytes, er := json.Marshal(coinValues)
		if er != nil {
			return nil, er
		}
		coinBalance := WalletBalance{}
		er = json.Unmarshal(coinValuesBytes, &coinBalance)
		if er != nil {
			return nil, er
		}
		result[coin] = coinBalance
	}

	return &result, nil
}
