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

func MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point float64) (qtyPercent float64) {
	var priceDiff float64
	if rangeA2 > rangeA1 {
		priceDiff = rangeA2 - rangeA1
	} else {
		priceDiff = rangeA1 - rangeA2
	}

	onePercent := Div(priceDiff, 100)
	percents := 0.0
	if onePercent > 0.0 {
		if rangeA2 > rangeA1 {
			percents = Div(point-rangeA1, onePercent)
		} else {
			percents = Div(rangeA1-rangeA2-(point-rangeA2), onePercent)
		}
	}

	qtyPercentRange := 0.0
	if rangeB2 > rangeB1 {
		qtyPercentRange = rangeB2 - rangeB1
	} else {
		qtyPercentRange = rangeB1 - rangeB2
	}
	qtyPercentDiff := Mul(Div(qtyPercentRange, 100), percents)

	if rangeB2 > rangeB1 {
		qtyPercent = rangeB1 + qtyPercentDiff
	} else {
		qtyPercent = rangeB1 - qtyPercentDiff
	}

	return
}
