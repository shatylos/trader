package math

import (
	"github.com/shopspring/decimal"
	"math"
)

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

func Round(value float64, precision int64) float64 {
	roundNum := math.Pow(10, float64(precision))
	return Div(math.Round(Mul(value, roundNum)), roundNum)
}

func RoundCell(value float64, precision int64) float64 {
	roundNum := math.Pow(10, float64(precision))
	return Div(math.Ceil(Mul(value, roundNum)), roundNum)
}

func RoundFloor(value float64, precision int64) float64 {
	roundNum := math.Pow(10, float64(precision))
	return Div(math.Floor(Mul(value, roundNum)), roundNum)
}
