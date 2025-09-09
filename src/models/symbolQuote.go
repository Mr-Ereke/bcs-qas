package models

import (
	"gitlab.online-fx.com/go-packages/logger"
)

type SymbolQuote struct {
	Symbol    string
	Bid       float64
	Ask       float64
	Last      float64
	Timestamp int64 // milliseconds
}

func (sq *SymbolQuote) GetQuoteByType(priceType string) float64 {
	switch priceType {
	case Bid:
		return sq.Bid
	case Ask:
		return sq.Ask
	case Last:
		return sq.Last
	default:
		logger.Errorf("unknow quote type. Type: %s", priceType)
	}

	return 0
}
