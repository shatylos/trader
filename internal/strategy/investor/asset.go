package investor

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"time"
)

func (i *Investor) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) (err error) {
	ctx := context.Background()
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
