package request

import "encoding/json"

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

func GetWalletBalance(isDemo bool) (*map[string]WalletBalance, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQueryGet("/v2/private/wallet/balance", params, isDemo)
	if er != nil {
		return nil, er
	}
	return mapWalletBalance(queryResp)
}

func mapWalletBalance(source map[string]interface{}) (*map[string]WalletBalance, error) {
	result := map[string]WalletBalance{}

	for coin, coinValues := range source {
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
