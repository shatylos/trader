package _struct

import (
	"html/template"
	"time"
)

type Report struct {
	DateFrom        time.Time
	DateTo          time.Time
	RevenuePercents float64
	Revenue         float64
	Currency        string
	InnerHtml       template.HTML
}

type StrategyInterface interface {
	SetConfig(interface{}, map[interface{}]interface{}) error
	IsInit() bool
	Initialise() error
	DoAction() error
	GetReport(from time.Time, to time.Time) (Report, error)
	Wait()
}
