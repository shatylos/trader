package binanceFeatures

import (
	"encoding/json"
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/request"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"math"
)

type ProviderPositionRisk struct {
	Symbol                 string `json:"symbol"`                 // BTCUSDT
	PositionSide           string `json:"positionSide"`           // BOTH
	PositionAmt            string `json:"positionAmt"`            // -0.001
	EntryPrice             string `json:"entryPrice"`             // 105326.5
	BreakEvenPrice         string `json:"breakEvenPrice"`         // 105284.3694
	MarkPrice              string `json:"markPrice"`              // 105336.60000000
	UnRealizedProfit       string `json:"unRealizedProfit"`       // -0.01010000
	LiquidationPrice       string `json:"liquidationPrice"`       // 15045103.95358566
	IsolatedMargin         string `json:"isolatedMargin"`         // 0
	Notional               string `json:"notional"`               // -105.33660000
	MarginAsset            string `json:"marginAsset"`            // USDT
	IsolatedWallet         string `json:"isolatedWallet"`         // 0
	InitialMargin          string `json:"initialMargin"`          // 5.26682999
	MaintMargin            string `json:"maintMargin"`            // 0.42134640
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 5.26682999
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 0
	Adl                    int    `json:"adl"`                    // 3
	BidNotional            string `json:"bidNotional"`            // 0
	AskNotional            string `json:"askNotional"`            // 0
	UpdateTime             int64  `json:"updateTime"`             // 1749309239526
}

func (d *DomainBinanceFutures) GetPosition(coinPare string) (position domainStructs.DomainPosition, err error) {

	positionApiRequest := request.ApiGetRequest{
		Uri: "/fapi/v3/positionRisk",
		ApiParams: binanceStructs.ApiParams{
			"symbol": coinPare,
		},
		Secrets: d.secrets,
	}

	var positionResponse binanceStructs.ApiResponse
	positionResponse, err = positionApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerPositions := make([]ProviderPositionRisk, 0)
	err = json.Unmarshal(positionResponse, &providerPositions)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance position response data: %s", positionResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	for _, providerPosition := range providerPositions {
		if providerPosition.Symbol == coinPare {
			var markPrice float64
			markPrice, err = _type.ToFloat64(providerPosition.MarkPrice)
			if err != nil {
				return
			}

			var totalPnl float64
			totalPnl, err = _type.ToFloat64(providerPosition.UnRealizedProfit)
			if err != nil {
				return
			}

			var size float64
			size, err = _type.ToFloat64(providerPosition.PositionAmt)
			if err != nil {
				return
			}
			size = math.Abs(size)

			var entryPrice float64
			entryPrice, err = _type.ToFloat64(providerPosition.EntryPrice)
			if err != nil {
				return
			}
			position = domainStructs.DomainPosition{
				Leverage:  0, // check to remove
				AvgPrice:  entryPrice,
				MarkPrice: markPrice,
				Size:      size,
				TotalPnl:  totalPnl,
				Side:      providerPosition.PositionSide,
				Symbol:    providerPosition.Symbol,
			}
			break
		}
	}

	orderApiRequest := request.ApiGetRequest{
		Uri: "/fapi/v1/openOrders",
		ApiParams: binanceStructs.ApiParams{
			"symbol": coinPare,
		},
		Secrets: d.secrets,
	}
	var orderResponse binanceStructs.ApiResponse
	orderResponse, err = orderApiRequest.DoRequest()
	if err != nil {
		return
	}

	providerOrders := make([]OrderResponse, 0)
	err = json.Unmarshal(orderResponse, &providerOrders)
	if err != nil {
		msg := fmt.Sprintf("Can not unmarshal Binance order response data: %s", orderResponse)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	for _, providerOrder := range providerOrders {
		if providerOrder.Symbol == coinPare && providerOrder.Status == orderStatuses.New && providerOrder.Type == orderTypeTP {
			if position.TakeProfit > 0 {
				msg := fmt.Sprintf("Second take profit value got from provider order list")
				logger.Error(msg)
				err = tools.AppError{
					Message: msg,
				}
				return
			}
			var stopPrice float64
			stopPrice, err = _type.ToFloat64(providerOrder.StopPrice)
			if err != nil {
				return
			}
			position.TakeProfit = stopPrice
		}
		if providerOrder.Symbol == coinPare && providerOrder.Status == orderStatuses.New && providerOrder.Type == orderTypeSL {
			if position.StopLoss > 0 {
				msg := fmt.Sprintf("Second stop loss value got from provider order list")
				logger.Error(msg)
				err = tools.AppError{
					Message: msg,
				}
				return
			}
			var stopPrice float64
			stopPrice, err = _type.ToFloat64(providerOrder.StopPrice)
			if err != nil {
				return
			}
			position.StopLoss = stopPrice
		}
	}

	return
}
