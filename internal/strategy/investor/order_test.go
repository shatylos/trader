package investor

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
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

func TestCalculateQtyToSellModify(t *testing.T) {
	wallet := domainStructs.DomainWallet{}
	wallet.Available = make([]domainStructs.DomainWalletCoinItem, 1)
	wallet.Available[0] = domainStructs.DomainWalletCoinItem{
		Coin:   "BTC",
		Amount: 0.00005,
	}

	i := Investor{
		Wallet: &wallet,
		config: Config{
			TradeCurrency:         "BTC",
			MinCoinReservePercent: 1,
			CommissionBuy:         0.18,
			QtyPrecision:          6,
		},
	}

	qty := 0.000051

	result := i.calculateQtyToSell(qty)
	expected := 0.00005

	if result != expected {
		t.Errorf("TestCalculateQtyToSell failed. Got value: %f. Expected: %f", result, expected)
	}
}

func TestCalculateQtyToSellNotModify(t *testing.T) {
	wallet := domainStructs.DomainWallet{}
	wallet.Available = make([]domainStructs.DomainWalletCoinItem, 1)
	wallet.Available[0] = domainStructs.DomainWalletCoinItem{
		Coin:   "BTC",
		Amount: 0.000051,
	}

	i := Investor{
		Wallet: &wallet,
		config: Config{
			TradeCurrency:         "BTC",
			MinCoinReservePercent: 1,
			CommissionBuy:         0.18,
			QtyPrecision:          6,
		},
	}

	qty := 0.00005

	result := i.calculateQtyToSell(qty)
	expected := 0.00005

	if result != expected {
		t.Errorf("TestCalculateQtyToSell failed. Got value: %f. Expected: %f", result, expected)
	}
}
