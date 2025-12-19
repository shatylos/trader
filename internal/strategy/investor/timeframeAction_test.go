package investor

import (
	"testing"
)

func TestModifySidewaysCoeff(t *testing.T) {

	origCoeff := 50.0
	percentage := 7.0

	i := Investor{}
	result := i.modifySidewaysCoeff(origCoeff, percentage)

	expected := 53.5
	if result != expected {
		t.Errorf("TestModifySidewaysCoeff failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestModifySidewaysCoeffMinPerc(t *testing.T) {

	origCoeff := 50.0
	percentage := -7.0

	i := Investor{}
	result := i.modifySidewaysCoeff(origCoeff, percentage)

	expected := 46.5
	if result != expected {
		t.Errorf("TestModifySidewaysCoeffMinPerc failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestModifySidewaysCoeffMinOrig(t *testing.T) {

	origCoeff := -50.0
	percentage := 7.0

	i := Investor{}
	result := i.modifySidewaysCoeff(origCoeff, percentage)

	expected := -46.5
	if result != expected {
		t.Errorf("TestModifySidewaysCoeffMinOrig failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestModifySidewaysCoeffMinOrigMinPerc(t *testing.T) {

	origCoeff := -50.0
	percentage := -7.0

	i := Investor{}
	result := i.modifySidewaysCoeff(origCoeff, percentage)

	expected := -53.5
	if result != expected {
		t.Errorf("TestModifySidewaysCoeffMinOrigMinPerc failed. Got value: %f. Expected: %f", result, expected)
	}
}
