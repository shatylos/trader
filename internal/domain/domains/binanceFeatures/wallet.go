package binanceFeatures

import (
	"encoding/json"
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/request"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	providerStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
)

type ProviderAccountInfo struct {
	TotalInitialMargin          string `json:"totalInitialMargin"`
	TotalMaintMargin            string `json:"totalMaintMargin"`
	TotalWalletBalance          string `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string `json:"totalUnrealizedProfit"`
	TotalMarginBalance          string `json:"totalMarginBalance"`
	TotalPositionInitialMargin  string `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string `json:"totalCrossUnPnl"`
	AvailableBalance            string `json:"availableBalance"`
	MaxWithdrawAmount           string `json:"maxWithdrawAmount"`
	Assets                      []struct {
		Asset                  string `json:"asset"`
		WalletBalance          string `json:"walletBalance"`
		UnrealizedProfit       string `json:"unrealizedProfit"`
		MarginBalance          string `json:"marginBalance"`
		MaintMargin            string `json:"maintMargin"`
		InitialMargin          string `json:"initialMargin"`
		PositionInitialMargin  string `json:"positionInitialMargin"`
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		CrossWalletBalance     string `json:"crossWalletBalance"`
		CrossUnPnl             string `json:"crossUnPnl"`
		AvailableBalance       string `json:"availableBalance"`
		MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
		UpdateTime             int64  `json:"updateTime"`
	} `json:"assets"`
	Positions []struct {
		Symbol           string `json:"symbol"`
		PositionSide     string `json:"positionSide"`
		PositionAmt      string `json:"positionAmt"`
		UnrealizedProfit string `json:"unrealizedProfit"`
		IsolatedMargin   string `json:"isolatedMargin"`
		Notional         string `json:"notional"`
		IsolatedWallet   string `json:"isolatedWallet"`
		InitialMargin    string `json:"initialMargin"`
		MaintMargin      string `json:"maintMargin"`
		UpdateTime       int    `json:"updateTime"`
	} `json:"positions"`
}

func (d *DomainBinanceFutures) GetWallet() (wallet providerStructs.DomainWallet, err error) {

	apiRequest := request.ApiGetRequest{
		Uri:       "/fapi/v3/account",
		ApiParams: binanceStructs.ApiParams{},
		Secrets:   d.secrets,
	}
	var apiResponse binanceStructs.ApiResponse
	apiResponse, err = apiRequest.DoRequest()
	if err != nil {
		return
	}

	var providerAccount ProviderAccountInfo
	err = json.Unmarshal(apiResponse, &providerAccount)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance account info: %s", apiResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	available := make([]providerStructs.DomainWalletCoinItem, 0)
	reserved := make([]providerStructs.DomainWalletCoinItem, 0)

	for _, asset := range providerAccount.Assets {
		var walletBalance, availableBalance float64
		walletBalance, err = _type.ToFloat64(asset.WalletBalance)
		if err != nil {
			return
		}
		if walletBalance == 0 {
			continue
		}
		availableBalance, err = _type.ToFloat64(asset.AvailableBalance)
		if err != nil {
			return
		}
		available = append(available, providerStructs.DomainWalletCoinItem{
			Coin:   asset.Asset,
			Amount: availableBalance,
		})
		reserved = append(reserved, providerStructs.DomainWalletCoinItem{
			Coin:   asset.Asset,
			Amount: walletBalance - availableBalance,
		})
	}

	wallet = providerStructs.DomainWallet{
		DomainCode: d.code,
		Available:  available,
		Reserved:   reserved,
	}
	return
}
