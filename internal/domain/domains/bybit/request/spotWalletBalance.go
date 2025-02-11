package request

import (
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/type"
	"strconv"
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
		return nil, tools.AppError{
			Message: fmt.Sprintf("Error marshalling source: %s", err.Error()),
		}
	}

	var rawResponse RawWalletBalanceResponse
	err = json.Unmarshal(sourceBytes, &rawResponse)
	if err != nil {
		return nil, tools.AppError{
			Message: fmt.Sprintf("Error unmarshalling source: %s", err.Error()),
		}
	}

	for _, wallet := range rawResponse.List {
		if wallet.AccountType == "UNIFIED" {
			for _, coin := range wallet.Coin {
				total, err := _type.ToFloat64(coin.Total)
				if err != nil {
					return nil, tools.AppError{
						Message: fmt.Sprintf("Error converting total to float64: %s", err.Error()),
					}
				}
				locked, err := _type.ToFloat64(coin.Locked)
				if err != nil {
					return nil, tools.AppError{
						Message: fmt.Sprintf("Error converting locked to float64: %s", err.Error()),
					}
				}
				coin.Free = strconv.FormatFloat(total-locked, 'f', 8, 64)
				result[coin.Coin] = coin
			}
		}
	}

	return &result, nil
}
