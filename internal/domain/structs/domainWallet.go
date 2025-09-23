package structs

type DomainWallet struct {
	Available []DomainWalletCoinItem
	Reserved  []DomainWalletCoinItem
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
