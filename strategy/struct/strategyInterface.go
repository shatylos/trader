package _struct

import (
	"html/template"
	"time"
)

type Report struct {
	InnerHtml template.HTML
}

type StrategyInterface interface {
	SetConfig(interface{}, map[interface{}]interface{}) error
	IsInit() bool
	Initialise() error
	DoAction() error
	GetReport(from time.Time, to time.Time) (*Report, error)
	Wait()
	ResetOrderData() error
}
