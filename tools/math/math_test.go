package math

import "testing"

func TestMapRangeFromZero(t *testing.T) {
	rangeA1 := 0.0
	rangeA2 := 10.0

	rangeB1 := 0.0
	rangeB2 := 1000.0

	point := 4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeFromZero2(t *testing.T) {
	rangeA1 := 0.0
	rangeA2 := 10.0

	rangeB1 := 1000.0
	rangeB2 := 0.0

	point := 4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeToZero(t *testing.T) {
	rangeA1 := 10.0
	rangeA2 := 0.0

	rangeB1 := 1000.0
	rangeB2 := 0.0

	point := 4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeToZero2(t *testing.T) {
	rangeA1 := 10.0
	rangeA2 := 0.0

	rangeB1 := 0.0
	rangeB2 := 1000.0

	point := 4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeFromNonZero(t *testing.T) {
	rangeA1 := 100.0
	rangeA2 := 200.0

	rangeB1 := 1000.0
	rangeB2 := 2000.0

	point := 140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 1400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeFromNonZero2(t *testing.T) {
	rangeA1 := 100.0
	rangeA2 := 200.0

	rangeB1 := 2000.0
	rangeB2 := 1000.0

	point := 140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 1600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeToNonZero(t *testing.T) {
	rangeA1 := 200.0
	rangeA2 := 100.0

	rangeB1 := 2000.0
	rangeB2 := 1000.0

	point := 140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 1400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeToNonZero2(t *testing.T) {
	rangeA1 := 200.0
	rangeA2 := 100.0

	rangeB1 := 1000.0
	rangeB2 := 2000.0

	point := 140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 1600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegFromZero(t *testing.T) {
	rangeA1 := 0.0
	rangeA2 := -10.0

	rangeB1 := 0.0
	rangeB2 := 1000.0

	point := -4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegFromZero2(t *testing.T) {
	rangeA1 := 0.0
	rangeA2 := -10.0

	rangeB1 := -1000.0
	rangeB2 := 0.0

	point := -4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegToZero(t *testing.T) {
	rangeA1 := -10.0
	rangeA2 := 0.0

	rangeB1 := -1000.0
	rangeB2 := 0.0

	point := -4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegToZero2(t *testing.T) {
	rangeA1 := -10.0
	rangeA2 := 0.0

	rangeB1 := 0.0
	rangeB2 := -1000.0

	point := -4.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegFromNonZero(t *testing.T) {
	rangeA1 := -100.0
	rangeA2 := -200.0

	rangeB1 := -1000.0
	rangeB2 := -2000.0

	point := -140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -1400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegFromNonZero2(t *testing.T) {
	rangeA1 := -100.0
	rangeA2 := -200.0

	rangeB1 := -2000.0
	rangeB2 := -1000.0

	point := -140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -1600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegToNonZero(t *testing.T) {
	rangeA1 := -200.0
	rangeA2 := -100.0

	rangeB1 := -2000.0
	rangeB2 := -1000.0

	point := -140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -1400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeNegToNonZero2(t *testing.T) {
	rangeA1 := -200.0
	rangeA2 := -100.0

	rangeB1 := -1000.0
	rangeB2 := -2000.0

	point := -140.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := -1600.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestMapRangeOut(t *testing.T) {
	rangeA1 := 10.0
	rangeA2 := 20.0

	rangeB1 := 1000.0
	rangeB2 := 2000.0

	point := 24.0

	result := MapRange(rangeA1, rangeA2, rangeB1, rangeB2, point)

	expected := 2400.0
	if result != expected {
		t.Errorf("TestMapRange failed. Got value: %f. Expected: %f", result, expected)
	}
}
