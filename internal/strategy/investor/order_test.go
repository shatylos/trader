package investor

import (
	"testing"
)

func TestAddCommission(t *testing.T) {

	qty := 0.00005
	commission := 0.18

	i := Investor{}
	result := i.addCommission(qty, commission)

	expected := 0.00005009
	if result != expected {
		t.Errorf("TestAddCommission failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestRemoveCommission(t *testing.T) {

	qty := 0.00005009
	commission := 0.18

	i := Investor{}
	result := i.removeCommission(qty, commission)

	expected := 0.00005
	if result != expected {
		t.Errorf("TestAddCommission failed. Got value: %f. Expected: %f", result, expected)
	}
}
