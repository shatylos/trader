package utils

import "github.com/shopspring/decimal"

func Mul(a, b float64) float64 {
	num1 := decimal.NewFromFloat(a)
	num2 := decimal.NewFromFloat(b)
	result, _ := num1.Mul(num2).Float64()
	return result
}

func Div(a, b float64) float64 {
	num1 := decimal.NewFromFloat(a)
	num2 := decimal.NewFromFloat(b)
	result, _ := num1.Div(num2).Float64()
	return result
}
