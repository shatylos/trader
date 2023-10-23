package request

import (
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
)

type SpotWalletBalance struct {
	Coin   string `json:"coin"`   //	Coin
	CoinId string `json:"coinId"` //	Coin ID
	Total  string `json:"total"`  //	Total equity
	Free   string `json:"free"`   //	Available balance
	Locked string `json:"locked"` //	Reserved for orders
}

func GetSpotWalletBalance(secrets bybitStructs.Secrets) (*map[string]SpotWalletBalance, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQueryGet("/spot/v3/private/account", params, secrets)
	if er != nil {
		return nil, er
	}
	return mapSpotWalletBalance(queryResp)
}

func mapSpotWalletBalance(source interface{}) (*map[string]SpotWalletBalance, error) {
	result := map[string]SpotWalletBalance{}

	sourceMap, ok := source.(map[string]interface{})
	if ok == false {
		return nil, utils.AppError{
			Message: "Error parse SpotWalletBalance response for ByBit",
		}
	}

	if sourceMap["balances"] == nil {
		return nil, utils.AppError{
			Message: "balances is not available in response",
		}
	}

	sourceBalances, ok := sourceMap["balances"].([]interface{})
	if !ok {
		return nil, utils.AppError{
			Message: "can not parse balances in response",
		}
	}

	for _, coinValues := range sourceBalances {
		coinValuesBytes, er := json.Marshal(coinValues)
		if er != nil {
			return nil, er
		}
		coinBalance := SpotWalletBalance{}
		er = json.Unmarshal(coinValuesBytes, &coinBalance)
		if er != nil {
			return nil, er
		}
		result[coinBalance.CoinId] = coinBalance
	}

	return &result, nil
}
