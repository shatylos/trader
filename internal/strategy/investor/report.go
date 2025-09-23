package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportTemplateData struct {
	PrevPeriodLink  string
	NextPeriodLink  string
	DateFrom        time.Time
	DateTo          time.Time
	DealsRelations  []*storage.DealRelation
	MainCurrency    string
	TradeCurrency   string
	PricePrecision  int
	TotalPnl        float64
	TotalPnlPercent float64
	BalanceBefore   float64
	BalanceAfter    float64
	DepoAmount      float64
	WithdrawAmount  float64
}

func (i *Investor) GetReport(from time.Time, to time.Time) (report _struct.Report, err error) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxSetupKey, i)

	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		return
	}

	var deals []*storage.Deal
	deals, err = i.Storage.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		return
	}

	var dealRelations []*storage.DealRelation
	for _, deal := range deals {
		var dealRelation *storage.DealRelation
		dealRelation, err = i.Storage.GetDealRelation(ctx, deal)
		if err != nil {
			return
		}
		dealRelations = append(dealRelations, dealRelation)
	}

	//var assets []structs.AssetTransaction
	//assets, err = storage.GetAssetTransactions(from, to)
	//totalPnl, totalPnlPercent, balanceBefore, balanceAfter, depoAmount, withdrawAmount := f.reportAmounts(positions, assets)

	data := ReportTemplateData{
		PrevPeriodLink:  fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:  fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:        from,
		DateTo:          to,
		DealsRelations:  dealRelations,
		MainCurrency:    i.config.MainCurrency,
		TradeCurrency:   i.config.TradeCurrency,
		PricePrecision:  int(i.config.PricePrecision),
		TotalPnl:        50,
		TotalPnlPercent: 50,
		BalanceBefore:   50,
		BalanceAfter:    50,
		DepoAmount:      50,
		WithdrawAmount:  50,
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		return
	}

	htmlStr := resultBuilder.String()

	report = _struct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   i.GetId(),
	}
	return
}
