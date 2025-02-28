package structs

type AssetTransaction struct {
	Id              *string `bson:"_id,omitempty"`
	TransactionType string  `bson:"type"`
	Amount          float64 `bson:"amount"`
	CreatedTime     int64   `bson:"created_time"`
}
