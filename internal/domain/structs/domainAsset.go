package structs

import "time"

type AssetTransaction struct {
	Id              *string   `bson:"_id,omitempty"`
	TransactionType string    `bson:"TransactionType"`
	Amount          float64   `bson:"Amount"`
	CreatedTime     time.Time `bson:"CreatedTime"`
}

const TransactionTypeDeposit = "deposit"
const TransactionTypeWithdraw = "withdraw"
