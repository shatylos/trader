package fibonacci

import "testing"

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
