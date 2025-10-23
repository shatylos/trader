package structs

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DomainWallet struct {
	Available   []DomainWalletCoinItem
	Reserved    []DomainWalletCoinItem
	UpdatedTime time.Time
}

type DomainWalletCoinItem struct {
	Coin   string
	Amount float64
}

type DomainWalletApiResult struct {
	DomainCode   string
	DomainWallet *DomainWallet
	Error        error
}

func (w *DomainWallet) IsEqual(w2 *DomainWallet) (result bool) {
	w1str := w.toString()
	w2str := w2.toString()

	result = w1str == w2str
	return
}

func (w *DomainWallet) toString() (result string) {
	values := make([]string, 0)

	for _, coin := range w.Reserved {
		if coin.Amount > 0.0 {
			values = append(values, fmt.Sprintf("R%s=%g", coin.Coin, coin.Amount))
		}
	}
	for _, coin := range w.Available {
		if coin.Amount > 0.0 {
			values = append(values, fmt.Sprintf("A%s=%g", coin.Coin, coin.Amount))
		}
	}

	sort.Strings(values)
	result = strings.Join(values, "|")

	return
}
