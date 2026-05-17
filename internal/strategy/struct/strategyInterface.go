package _struct

import (
	"html/template"
	"net/http"
	"time"
)

type Report struct {
	InnerHtml template.HTML
	SetupId   string
}

type Stats struct {
	SetupId               string
	PNLTotal              Pnl
	PNL12Months           Pnl
	PNLLastMonth          Pnl
	WithdrawablePrevMonth float64
}

type Pnl struct {
	Amount         float64
	Percent        float64
	AvPercPerMonth float64
	Currency       string
}

type StrategyInterface interface {
	GetId() string
	IsEnabled() bool
	GetTitle() string
	SetConfig(interface{}, map[string]interface{}) error
	DoAction() error
	GetReport(from time.Time, to time.Time) (Report, error)
	GetStats() (Stats, error)
	WaitDuration() time.Duration
	AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) error
	Init(mux *http.ServeMux)
}
