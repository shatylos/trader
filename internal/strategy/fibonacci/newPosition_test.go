package fibonacci

import (
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"testing"
)

func TestGetMinMaxPriceLong(t *testing.T) {
}

func TestGetMinMaxPriceShort(t *testing.T) {
}

func TestCalculateLongFibonacciChart(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibEntryPoint1: 0.618,
		FibEntryPoint2: 0.5,
		FibEntryPoint3: 0.382,
		FibStopLoss:    0,
		FibTakeProfit1: 1,
		FibTakeProfit2: 1.272,
		FibTakeProfit3: 1.618,
		PricePrecision: 4,
	}
	minPrice := 3184.87
	maxPrice := 3623.58
	chart := f.calculateFibonacciChart(minPrice, maxPrice, true)

	expected := 3455.9928
	if chart.EntryPoint1 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint1 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint1, expected)
	}
	expected = 3404.225
	if chart.EntryPoint2 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint2 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint2, expected)
	}
	expected = 3352.4572
	if chart.EntryPoint3 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint3 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint3, expected)
	}
	expected = 3184.87
	if chart.StopLoss != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).StopLoss = %f; expects %f", minPrice, maxPrice, chart.StopLoss, expected)
	}
	expected = 3623.58
	if chart.TakeProfit1 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit1 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit1, expected)
	}
	expected = 3742.9091
	if chart.TakeProfit2 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit2 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit2, expected)
	}
	expected = 3894.7028
	if chart.TakeProfit3 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit3 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit3, expected)
	}
}

func TestCalculateShortFibonacciChart(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibEntryPoint1: 0.618,
		FibEntryPoint2: 0.5,
		FibEntryPoint3: 0.382,
		FibStopLoss:    0,
		FibTakeProfit1: 1,
		FibTakeProfit2: 1.272,
		FibTakeProfit3: 1.618,
		PricePrecision: 4,
	}
	minPrice := 16280.32
	maxPrice := 23654.88
	chart := f.calculateFibonacciChart(minPrice, maxPrice, false)

	expected := 19097.4019
	if chart.EntryPoint1 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint1 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint1, expected)
	}
	expected = 19967.6
	if chart.EntryPoint2 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint2 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint2, expected)
	}
	expected = 20837.7981
	if chart.EntryPoint3 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).EntryPoint3 = %f; expects %f", minPrice, maxPrice, chart.EntryPoint3, expected)
	}
	expected = 23654.88
	if chart.StopLoss != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).StopLoss = %f; expects %f", minPrice, maxPrice, chart.StopLoss, expected)
	}
	expected = 16280.32
	if chart.TakeProfit1 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit1 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit1, expected)
	}
	expected = 14274.4397
	if chart.TakeProfit2 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit2 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit2, expected)
	}
	expected = 11722.8419
	if chart.TakeProfit3 != expected {
		t.Errorf("calculateFibonacciChart(%f, %f, true).TakeProfit3 = %f; expects %f", minPrice, maxPrice, chart.TakeProfit3, expected)
	}
}

func TestCalculateFullQtyBull(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibStopLoss:         0,
		QtyPrecision:        4,
		RiskPercent:         10,
		EP1ToFullQtyPercent: 20,
		EP2ToFullQtyPercent: 30,
		EP3ToFullQtyPercent: 50,
	}
	fibChart := structs.FibonacciChart{
		EntryPoint1:    80000,
		EntryPoint2:    76000,
		EntryPoint3:    73000,
		SourceMaxPrice: 0,
		SourceMinPrice: 0,
		StopLoss:       70000,
		TakeProfit1:    0,
		TakeProfit2:    0,
		TakeProfit3:    0,
		FullQty:        0,
	}

	fullQty, err := f.calculateFullQty(500.0, fibChart)
	if err != nil {
		t.Errorf("Error calculateFullQty: %s", err.Error())
	}

	expected := 0.0094
	if fullQty != expected {
		t.Errorf("calculateFullQty Bull. FullQty %f; expects %f", fullQty, expected)
	}
}

func TestCalculateFullQtyBear(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibStopLoss:         0,
		QtyPrecision:        4,
		RiskPercent:         10,
		EP1ToFullQtyPercent: 20,
		EP2ToFullQtyPercent: 30,
		EP3ToFullQtyPercent: 50,
	}
	fibChart := structs.FibonacciChart{
		EntryPoint1:    70000,
		EntryPoint2:    74000,
		EntryPoint3:    77000,
		SourceMaxPrice: 0,
		SourceMinPrice: 0,
		StopLoss:       80000,
		TakeProfit1:    0,
		TakeProfit2:    0,
		TakeProfit3:    0,
		FullQty:        0,
	}

	fullQty, err := f.calculateFullQty(500.0, fibChart)
	if err != nil {
		t.Errorf("Error calculateFullQty: %s", err.Error())
	}

	expected := 0.0094
	if fullQty != expected {
		t.Errorf("calculateFullQty Bull. FullQty %f; expects %f", fullQty, expected)
	}
}

func TestCalculateFullQtyOneOrder1(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibStopLoss:         0,
		QtyPrecision:        4,
		RiskPercent:         10,
		EP1ToFullQtyPercent: 20,
		EP2ToFullQtyPercent: 30,
		EP3ToFullQtyPercent: 50,
	}
	fibChart := structs.FibonacciChart{
		EntryPoint1: 80000,
		EntryPoint2: 0,
		EntryPoint3: 0,
		StopLoss:    70000,
		FullQty:     0,
	}

	fullQty, err := f.calculateFullQty(500.0, fibChart)
	if err != nil {
		t.Errorf("Error calculateFullQty: %s", err.Error())
	}

	expected := 0.025
	if fullQty != expected {
		t.Errorf("calculateFullQty Bull. FullQty %f; expects %f", fullQty, expected)
	}
}

func TestCalculateFullQtyOneOrder2(t *testing.T) {
	f := Fibonacci{}
	f.config = Config{
		FibStopLoss:         0,
		QtyPrecision:        4,
		RiskPercent:         10,
		EP1ToFullQtyPercent: 20,
		EP2ToFullQtyPercent: 0,
		EP3ToFullQtyPercent: 0,
	}
	fibChart := structs.FibonacciChart{
		EntryPoint1:    80000,
		EntryPoint2:    76000,
		EntryPoint3:    73000,
		SourceMaxPrice: 0,
		SourceMinPrice: 0,
		StopLoss:       70000,
		FullQty:        0,
	}

	fullQty, err := f.calculateFullQty(500.0, fibChart)
	if err != nil {
		t.Errorf("Error calculateFullQty: %s", err.Error())
	}

	expected := 0.025
	if fullQty != expected {
		t.Errorf("calculateFullQty Bull. FullQty %f; expects %f", fullQty, expected)
	}
}
