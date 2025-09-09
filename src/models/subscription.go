package models

import (
	"time"
)

type Subscription struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	CustomerId     uint    `gorm:"not null; index" json:"customerId"`
	Instrument     string  `gorm:"type: varchar(255); not null" json:"instrument"`
	DisplayTitle   string  `gorm:"type: varchar(255); not null" json:"displayTitle"`
	Price          float64 `gorm:"type: numeric(16,4); not null" json:"price"`
	InitialPrice   float64 `gorm:"type: numeric(16,4); not null" json:"initialPrice"`
	PriceDirection bool    `gorm:"not null" json:"priceDirection"`
	PriceType      string  `gorm:"type: varchar(10); not null" json:"priceType"`
	Currency       string  `gorm:"type: varchar(10); not null" json:"currency"`
	Digits         int8    `gorm:"type: smallint; not null" json:"digits"`

	QuoteTimestamp int64 `gorm:"-"` // quote time in milliseconds

	CreatedAt time.Time
}

func (m Subscription) TableName() string {
	return "subscription"
}
