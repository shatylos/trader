package structs

type HistoryOrder struct {
	Id                        *string  `bson:"_id,omitempty"`
	DomainOrderId             string   `bson:"domain_order_id"`
	FilledPrice               *float64 `bson:"filled_price"`
	FilledQty                 *float64 `bson:"filled_qty"`
	Side                      *string  `bson:"side"`
	CreatedTime               int64    `bson:"created_time"`
	UpdatedTime               *int64   `bson:"updated_time"`
	MainCurrencyAmountBefore  float64  `bson:"main_currency_amount_before"`
	TradeCurrencyAmountBefore float64  `bson:"trade_currency_amount_before"`
	AveragePrice              *float64 `bson:"average_price"`
}
