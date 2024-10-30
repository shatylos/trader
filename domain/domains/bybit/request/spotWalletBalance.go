package request

import (
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/utils"
	"encoding/json"
	"fmt"
)

type SpotWalletBalance struct {
	Coin   string `json:"coin"`                //	Coin
	Total  string `json:"walletBalance"`       //	Total equity
	Free   string `json:"availableToWithdraw"` //	Available balance
	Locked string `json:"locked"`              //	Reserved for orders
}

type RawWalletBalanceResponse struct {
	List []struct {
		Coin        []SpotWalletBalance `json:"coin"`
		AccountType string              `json:"accountType"`
	} `json:"list"`
}

func GetSpotWalletBalance(secrets bybitStructs.Secrets) (*map[string]SpotWalletBalance, error) {
	params := make(ApiParams, 0)
	params["accountType"] = "UNIFIED"
	queryResp, er := apiQueryGet("/v5/account/wallet-balance", params, secrets)
	if er != nil {
		return nil, er
	}
	return mapSpotWalletBalance(queryResp)
}

func mapSpotWalletBalance(source interface{}) (*map[string]SpotWalletBalance, error) {
	result := map[string]SpotWalletBalance{}
	sourceBytes, err := json.Marshal(source)

	if err != nil {
		return nil, utils.AppError{
			Message: fmt.Sprintf("Error marshalling source: %s", err.Error()),
		}
	}

	var rawResponse RawWalletBalanceResponse
	err = json.Unmarshal(sourceBytes, &rawResponse)
	if err != nil {
		return nil, utils.AppError{
			Message: fmt.Sprintf("Error unmarshalling source: %s", err.Error()),
		}
	}

	for _, wallet := range rawResponse.List {
		if wallet.AccountType == "UNIFIED" {
			for _, coin := range wallet.Coin {
				result[coin.Coin] = coin
			}
		}
	}

	return &result, nil
}
