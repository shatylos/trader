package investor

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/web/helper"
)

func TestReportTemplateRenders(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Chdir("../../..")
	defer os.Chdir(wd)

	tmpl, err := helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	order := &entity.Order{Timeframe: "1h", ConfigKey: 0, AverageBuyPrice: 100, QtyInTrade: 0.5, RealizedPNL: 1.2, CreatedTime: time.Now(), UpdatedTime: time.Now()}
	order.Side = "BUY"
	order.OrderStatus = "FILLED"
	order.Qty = 0.5
	order.Price = 100

	tf := _struct.TimeframeItem{Config: _struct.TimeframeItemConfig{Resolution: "1h"}}
	state := &entity.TimeframeState{Timeframe: "1h", AverageBuyPrice: 100, QtyInTrade: 0.5, RealizedPNL: 1, UnrealizedPNL: -2}

	data := ReportTemplateData{
		DateFrom:        time.Now(),
		DateTo:          time.Now(),
		Orders:          map[string][]*entity.Order{"1h": {order}},
		TimeframeStates: map[string]*entity.TimeframeState{"1h": state},
		MainCurrency:    "USDT",
		TradeCurrency:   "BTC",
		PricePrecision:  2,
		QtyPrecision:    5,
		Timeframes:      []_struct.TimeframeItem{tf},
		IsCurrentPeriod: true,
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, data)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
	if !strings.Contains(sb.String(), "Qty In Trade") {
		t.Fatalf("expected rendered output to contain orders table")
	}

	// past period: no states block, no current-price block
	data.IsCurrentPeriod = false
	sb.Reset()
	if err = tmpl.Execute(&sb, data); err != nil {
		t.Fatalf("exec error (past period): %v", err)
	}
}
