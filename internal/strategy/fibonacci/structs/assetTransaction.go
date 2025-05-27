package structs

type AssetTransaction struct {
	Id              *string `bson:"_id,omitempty"`
	TransactionType string  `bson:"TransactionType"`
	Amount          float64 `bson:"Amount"`
	CreatedTime     int64   `bson:"CreatedTime"`
}
