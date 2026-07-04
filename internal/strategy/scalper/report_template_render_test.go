package scalper

import (
	"os"
	"strings"
	"testing"
	"time"

	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/scalper/structs"
	"github.com/shatylos/trader/web/helper"
)

func TestReportTemplateRenders(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Chdir("../../..")
	defer os.Chdir(wd)

	tmpl, err := helper.GetTemplate("web/template/scalper/report.html")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	position := structs.Position{
		Chart:       structs.ScalpChart{Rsi: 55, Atr: 12, EntryPrice: 100, Qty: 0.5},
		Bias:        "LONG",
		Status:      structs.StatusActive,
		CreatedTime: time.Now(),
		ProviderPosition: domainStructs.DomainPosition{
			AvgPrice: 100, MarkPrice: 101, Size: 0.5, Leverage: 10, Side: "Buy",
			TotalPnl: 0.5, TakeProfit: 110, StopLoss: 95,
		},
	}

	data := ReportTemplateData{
		WsLink:          "/scalper-test/ws-report",
		DateFrom:        time.Now(),
		DateTo:          time.Now(),
		Positions:       []structs.Position{position},
		MainCurrency:    "USDT",
		TradeCurrency:   "BTC",
		PricePrecision:  2,
		State:           State{Bias: "LONG", Rsi: 55, CurrentPrice: 101, SkippedMessage: "waiting"},
		IsCurrentPeriod: true,
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, data)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
	rendered := sb.String()
	for _, id := range []string{`id="pnl-ws"`, `id="current-state-ws"`, `id="positions-ws"`} {
		if !strings.Contains(rendered, id) {
			t.Fatalf("expected rendered output to contain %s", id)
		}
	}

	// each block must render standalone since it is broadcast over the web socket
	for _, blockName := range []string{"pnl", "current-state", "positions"} {
		sb.Reset()
		if err = tmpl.ExecuteTemplate(&sb, blockName, data); err != nil {
			t.Fatalf("exec error for block %s: %v", blockName, err)
		}
	}

	// past period: blocks keep their ids without the -ws suffix
	data.IsCurrentPeriod = false
	sb.Reset()
	if err = tmpl.Execute(&sb, data); err != nil {
		t.Fatalf("exec error (past period): %v", err)
	}
	if !strings.Contains(sb.String(), `id="positions"`) {
		t.Fatalf("expected past period output to contain plain positions id")
	}
}
