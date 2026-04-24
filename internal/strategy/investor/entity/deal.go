package entity

import "time"

const DealStatusNew = "NEW"
const DealStatusActive = "ACTIVE"
const DealStatusClosed = "CLOSED"

type Deal struct {
	Id          *string   `bson:"_id,omitempty"`
	Timeframe   string    `bson:"Timeframe"`
	Status      string    `bson:"Status"`
	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`
	ClosedTime  time.Time `bson:"ClosedTime"`
}

func (d *Deal) SetClose() {
	d.Status = DealStatusClosed
	d.ClosedTime = time.Now()
}
