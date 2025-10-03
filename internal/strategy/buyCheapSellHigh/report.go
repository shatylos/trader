package buyCheapSellHigh

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type StrategyReportPage struct {
	PrevPeriodLink         string
	NextPeriodLink         string
	DateFrom               time.Time
	DateTo                 time.Time
	Revenue                float64
	RevenuePercents        float64
	LocalRevenue           float64
	LocalRevenuePercents   float64
	LocalAvgPriceBuy       float64
	LocalAvgPriceSell      float64
	Commission             float64
	MainCurrency           string
	TradeCurrency          string
	MainCurrencyPrecision  int
	TradeCurrencyPrecision int
	ReportOrderItems       []ReportOrderItem
	Amounts                ReportAmounts
	Assets                 ReportAssets
}

type ReportAmounts struct {
	BeginMainCurrency  float64
	BeginTradeCurrency float64
	BeginTotal         float64
	EndMainCurrency    float64
	EndTradeCurrency   float64
	EndTotal           float64
	DiffMainCurrency   float64
	DiffTradeCurrency  float64
	DiffTotal          float64
}

type ReportAssets struct {
	Deposit  float64
	Withdraw float64
}

type ReportOrderItem struct {
	DateTime                       time.Time
	OrderId                        string
	Direction                      string
	Price                          float64
	MainCurrencyAmount             float64
	TradeCurrencyAmount            float64
	AvaragePrice                   float64
	Revenue                        float64
	Commission                     float64
	TotalRevenue                   float64
	TotalMainCurrencyAmountBefore  float64
	TotalTradeCurrencyAmountBefore float64
}

func (s *BuyCheapSellHigh) GetReport(from time.Time, to time.Time) (report _struct.Report, err error) {

	reportItems, err := s.getReportOrderItems(from, to)
	if err != nil {
		return
	}

	deposit, withdraw, err := s.getAssets(from, to)
	if err != nil {
		return
	}

	tmpl, err := helper.GetTemplate("web/template/buyCheapSellHigh/report.html")
	if err != nil {
		return
	}

	revenueAmount, revenuePercent, err := s.getRevenue(reportItems)
	if err != nil {
		return
	}

	revenueLocal, revenueLocalPercents, AvgPriceBuy, AvgPriceSell, commission, err := s.getRevenueLocal(reportItems)

	BeginMainCurrency := float64(0)
	BeginTradeCurrency := float64(0)
	BeginTotal := float64(0)
	EndMainCurrency := float64(0)
	EndTradeCurrency := float64(0)
	EndTotal := float64(0)

	if len(reportItems) > 0 {
		BeginMainCurrency = reportItems[len(reportItems)-1].TotalMainCurrencyAmountBefore
		BeginTradeCurrency = reportItems[len(reportItems)-1].TotalTradeCurrencyAmountBefore
		BeginTotal = BeginMainCurrency + math.Mul(BeginTradeCurrency, reportItems[len(reportItems)-1].Price)
		EndMainCurrency = reportItems[0].TotalMainCurrencyAmountBefore
		EndTradeCurrency = reportItems[0].TotalTradeCurrencyAmountBefore
		EndTotal = EndMainCurrency + math.Mul(EndTradeCurrency, reportItems[0].Price)
	}

	data := StrategyReportPage{
		PrevPeriodLink:         fmt.Sprintf("/report/%s/%s/", s.Id, from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:         fmt.Sprintf("/report/%s/%s/", s.Id, from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:               from,
		DateTo:                 to,
		Revenue:                revenueAmount,
		RevenuePercents:        revenuePercent,
		LocalRevenue:           revenueLocal,
		LocalRevenuePercents:   revenueLocalPercents,
		LocalAvgPriceBuy:       AvgPriceBuy,
		LocalAvgPriceSell:      AvgPriceSell,
		Commission:             commission,
		MainCurrency:           s.MainCurrency,
		TradeCurrency:          s.TradeCurrency,
		MainCurrencyPrecision:  int(s.MainCurrencyPrecision),
		TradeCurrencyPrecision: int(s.PurchaseVolumePrecision),
		ReportOrderItems:       reportItems,
		Amounts: ReportAmounts{
			BeginMainCurrency:  BeginMainCurrency,
			BeginTradeCurrency: BeginTradeCurrency,
			BeginTotal:         BeginTotal,
			EndMainCurrency:    EndMainCurrency,
			EndTradeCurrency:   EndTradeCurrency,
			EndTotal:           EndTotal,
			DiffMainCurrency:   EndMainCurrency - BeginMainCurrency,
			DiffTradeCurrency:  EndTradeCurrency - BeginTradeCurrency,
			DiffTotal:          EndTotal - BeginTotal,
		},
		Assets: ReportAssets{
			Deposit:  deposit,
			Withdraw: withdraw,
		},
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		return
	}

	htmlStr := resultBuilder.String()

	report = _struct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   s.Id,
	}

	return
}

func (s *BuyCheapSellHigh) getReportOrderItems(from time.Time, to time.Time) ([]ReportOrderItem, error) {

	var reportOrderItems []ReportOrderItem

	storage, err := strategyStorage.GetStorage(s.Id)
	if err != nil {
		return reportOrderItems, err
	}

	historyOrders, err := storage.GetCalculatedHistoryOrders(from, to)
	if err != nil {
		return reportOrderItems, err
	}

	for _, historyOrder := range historyOrders {
		reportOrderItems = append(reportOrderItems, ReportOrderItem{
			DateTime:                       time.Unix(historyOrder.UpdatedTime, 0),
			OrderId:                        historyOrder.DomainOrderId,
			Direction:                      historyOrder.Side,
			Price:                          historyOrder.FilledPrice,
			MainCurrencyAmount:             math.Mul(historyOrder.FilledPrice, historyOrder.FilledQty),
			TradeCurrencyAmount:            historyOrder.FilledQty,
			AvaragePrice:                   historyOrder.AveragePrice,
			Revenue:                        historyOrder.Revenue,
			Commission:                     historyOrder.Comission,
			TotalRevenue:                   historyOrder.Revenue - historyOrder.Comission,
			TotalMainCurrencyAmountBefore:  historyOrder.MainCurrencyAmountBefore,
			TotalTradeCurrencyAmountBefore: historyOrder.TradeCurrencyAmountBefore,
		})
	}

	return reportOrderItems, nil

}

func (s *BuyCheapSellHigh) getAssets(from time.Time, to time.Time) (float64, float64, error) {
	storage, err := strategyStorage.GetStorage(s.Id)
	assets, err := storage.GetAssetTransactions(from, to)
	if err != nil {
		return 0, 0, err
	}
	var deposit, withdraw float64

	for _, asset := range assets {
		if asset.TransactionType == structs.TransactionTypeDeposit {
			deposit += asset.Amount
		}
		if asset.TransactionType == structs.TransactionTypeWithdraw {
			withdraw += asset.Amount
		}
	}

	return deposit, withdraw, nil
}

func (s *BuyCheapSellHigh) getRevenue(reportOrderItems []ReportOrderItem) (float64, float64, error) {

	if len(reportOrderItems) == 0 {
		return 0, 0, nil
	}

	var revenue float64
	var revenuePercent float64

	for _, reportOrderItem := range reportOrderItems {
		revenue += reportOrderItem.TotalRevenue
	}

	firstOrder := reportOrderItems[len(reportOrderItems)-1]
	startSum := firstOrder.TotalMainCurrencyAmountBefore + math.Mul(firstOrder.TotalTradeCurrencyAmountBefore, firstOrder.Price)
	onePercentStart := math.Div(startSum, 100)
	revenuePercent = math.Div(revenue, onePercentStart)

	return revenue, revenuePercent, nil
}

func (s *BuyCheapSellHigh) getRevenueLocal(reportOrderItems []ReportOrderItem) (float64, float64, float64, float64, float64, error) {

	if len(reportOrderItems) == 0 {
		return 0, 0, 0, 0, 0, nil
	}

	var buyTradeCurrencyAmount float64
	var buyMainCurrencyAmount float64
	var sellTradeCurrencyAmount float64
	var sellMainCurrencyAmount float64
	var comission float64

	for _, reportOrderItem := range reportOrderItems {
		if reportOrderItem.Direction == structs.OrderSideBuy {
			buyTradeCurrencyAmount += reportOrderItem.TradeCurrencyAmount
			buyMainCurrencyAmount += reportOrderItem.MainCurrencyAmount
		}
		if reportOrderItem.Direction == structs.OrderSideSell {
			sellTradeCurrencyAmount += reportOrderItem.TradeCurrencyAmount
			sellMainCurrencyAmount += reportOrderItem.MainCurrencyAmount
		}
		comission += reportOrderItem.Commission
	}

	if buyTradeCurrencyAmount == 0 || sellTradeCurrencyAmount == 0 {
		return 0, 0, 0, 0, 0, nil
	}

	AvgPriceBuy := math.Div(buyMainCurrencyAmount, buyTradeCurrencyAmount)
	AvgPriceSell := math.Div(sellMainCurrencyAmount, sellTradeCurrencyAmount)

	var minValue float64
	if buyTradeCurrencyAmount > sellTradeCurrencyAmount {
		minValue = sellTradeCurrencyAmount
	} else {
		minValue = buyTradeCurrencyAmount
	}

	revenue := math.Mul(minValue, AvgPriceSell) - math.Mul(minValue, AvgPriceBuy) - comission

	firstOrder := reportOrderItems[len(reportOrderItems)-1]
	totalAmount := firstOrder.TotalMainCurrencyAmountBefore + firstOrder.TotalTradeCurrencyAmountBefore
	revenuePercent := math.Div(revenue, math.Div(totalAmount, 100))

	return revenue, revenuePercent, AvgPriceBuy, AvgPriceSell, comission, nil
}
