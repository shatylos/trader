package utils

import "github.com/shopspring/decimal"

func Mul(a, b float64) float64 {
	num1 := decimal.NewFromFloat(a)
	num2 := decimal.NewFromFloat(b)
	result := num1.Mul(num2)
	resultFloat, _ := result.Float64()
	return resultFloat
}

func Div(a, b float64) float64 {
	num1 := decimal.NewFromFloat(a)
	num2 := decimal.NewFromFloat(b)
	result := num1.Div(num2)
	resultFloat, _ := result.Float64()
	return resultFloat
}
