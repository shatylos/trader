package structs

import "time"

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
