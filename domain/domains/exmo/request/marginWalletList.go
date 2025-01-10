package request

import (
	"github.com/shatylos/trader/utils"
	"strconv"
)

type MarginWalletListStruct struct {
	Wallets map[string]wallet
}

type wallet struct {
	Balance         float64
	Used            float64
	UsedInPositions float64
	UsedInOrders    float64
	Free            float64
	Updated         int64
}

func GetMarginWalletList() (*MarginWalletListStruct, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQuery("margin/user/wallet/list", params)
	if er != nil {
		return nil, er
	}
	return mapWalletList(queryResp)
}

func mapWalletList(source map[string]interface{}) (*MarginWalletListStruct, error) {

	walletsResult := MarginWalletListStruct{
		Wallets: map[string]wallet{},
	}

	wallets, ok := source["wallets"].(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "[Exmo GetMarginWalletList] Can not parse broker response. Bad wallets."}
	}
	for coin, walletInterface := range wallets {
		walletItem, ok := walletInterface.(map[string]interface{})
		if !ok {
			return nil, utils.AppError{Message: "[Exmo GetMarginWalletList] Can not parse broker response. Bad wallet value."}
		}
		walletResult := wallet{}
		for balanceType, value := range walletItem {
			valueStr, ok := value.(string)
			if !ok {
				return nil, utils.AppError{Message: "[Exmo GetMarginWalletList] Can not parse broker response. Bad wallet item value."}
			}

			var er error
			if balanceType == "balance" {
				walletResult.Balance, er = parseFloat(valueStr)
			} else if balanceType == "used" {
				walletResult.Used, er = parseFloat(valueStr)
			} else if balanceType == "used_in_positions" {
				walletResult.UsedInPositions, er = parseFloat(valueStr)
			} else if balanceType == "used_in_orders" {
				walletResult.UsedInOrders, er = parseFloat(valueStr)
			} else if balanceType == "free" {
				walletResult.Free, er = parseFloat(valueStr)
			} else if balanceType == "updated" {
				walletResult.Updated, er = parseInt(valueStr)
			}
			if er != nil {
				return nil, er
			}
		}
		walletsResult.Wallets[coin] = walletResult
	}

	return &walletsResult, nil
}

func parseFloat(sourceValue string) (float64, error) {
	balanceFloat, er := strconv.ParseFloat(sourceValue, 64)
	if er != nil {
		return 0, er
	}
	return balanceFloat, nil
}

func parseInt(sourceValue string) (int64, error) {
	balanceInt, er := strconv.ParseInt(sourceValue, 10, 64)
	if er != nil {
		return 0, er
	}
	return balanceInt, nil
}
