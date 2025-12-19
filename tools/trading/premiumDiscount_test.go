package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"testing"
)

func TestPDPremium100(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000 // max
	klines[0].Low = 4000
	klines[0].Close = 5000 // price

	klines[1].High = 4500
	klines[1].Low = 3500 // min

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (5000 - 4250) / ((5000 - 4250) / 100)

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := 100.0
	if pdKoef != expected {
		t.Errorf("TestPDPremium100 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDPremium60(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000 // max
	klines[0].Low = 4000
	klines[0].Close = 4700 // price

	klines[1].High = 4500
	klines[1].Low = 3500 // min

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (4700 - 4250) / ((5000 - 4250) / 100)

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := 60.0
	if pdKoef != expected {
		t.Errorf("TestPDPremium60 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDPremium20(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000 // max
	klines[0].Low = 4000
	klines[0].Close = 4400 // price

	klines[1].High = 4500
	klines[1].Low = 3500 // min

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (4400 - 4250) / ((5000 - 4250) / 100)

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := 20.0
	if pdKoef != expected {
		t.Errorf("TestPDPremium20 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDDiscount100(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000  // max
	klines[0].Low = 3500   // min
	klines[0].Close = 3500 // price

	klines[1].High = 4500
	klines[1].Low = 4000

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (4250 - 3500) / ((4250 - 3500) / 100) * -1

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := -100.0
	if pdKoef != expected {
		t.Errorf("TestPDDiscount100 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDDiscount60(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000  // max
	klines[0].Low = 3500   // min
	klines[0].Close = 3800 // price

	klines[1].High = 4500
	klines[1].Low = 4000

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (4250 - 3800) / ((4250 - 3500) / 100) * -1

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := -60.0
	if pdKoef != expected {
		t.Errorf("TestPDDiscount60 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDDiscount20(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000  // max
	klines[0].Low = 3500   // min
	klines[0].Close = 4100 // price

	klines[1].High = 4500
	klines[1].Low = 4000

	klines[2].High = 4500
	klines[2].Low = 4000

	// middlePrice = (5000 - 3500) / 2 + 3500 = 4250
	// pdKoef = (4250 - 4100) / ((4250 - 3500) / 100) * -1

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := -20.0
	if pdKoef != expected {
		t.Errorf("TestPDDiscount20 failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}

func TestPDMiddle(t *testing.T) {
	klines := make([]domainStructs.DomainCandle, 3)

	klines[0].High = 5000  // max
	klines[0].Low = 3000   // min
	klines[0].Close = 4000 // price

	klines[1].High = 4500
	klines[1].Low = 4000

	klines[2].High = 4500
	klines[2].Low = 4000

	pdKoef, _, _ := PremiumDiscount(klines)

	expected := 0.0
	if pdKoef != expected {
		t.Errorf("TestPDMiddle failed. Got value: %f. Expected: %f", pdKoef, expected)
	}
}
