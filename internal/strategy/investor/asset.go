package investor

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"time"
)

func (i *Investor) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) (err error) {
	ctx := i.getContext()
	err = i.Storage.AddAssetTransaction(ctx, structs.AssetTransaction{
		TransactionType: transactionType,
		Amount:          amount,
		CreatedTime:     dateTime,
	})
	if err != nil {
		return
	}
	err = i.updateWalletInfo()
	if err != nil {
		return
	}

	return
}
