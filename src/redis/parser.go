package redis

import (
	"fmt"
	"strconv"
	"strings"

	"app/models"
)

var (
	ErrInvalidQuoteData          = fmt.Errorf("invalid quote data")
	ErrInvalidBidQuoteData       = fmt.Errorf("invalid bid quote data")
	ErrInvalidAskQuoteData       = fmt.Errorf("invalid ask quote data")
	ErrInvalidLastQuoteData      = fmt.Errorf("invalid last quote data")
	ErrInvalidTimestampQuoteData = fmt.Errorf("invalid timestamp quote data")
	ErrNoQuoteData               = fmt.Errorf("no symbol quote data")
)

func parseQuote(symbol string, quote string) (*models.SymbolQuote, error) {
	quoteData := strings.Split(quote, ";")

	if len(quoteData) == 5 {
		bid, err := strconv.ParseFloat(quoteData[0], 64)
		if err != nil {
			return nil, ErrInvalidBidQuoteData
		}

		ask, err := strconv.ParseFloat(quoteData[1], 64)
		if err != nil {
			return nil, ErrInvalidAskQuoteData
		}

		last, err := strconv.ParseFloat(quoteData[3], 64)
		if err != nil {
			return nil, ErrInvalidLastQuoteData
		}

		timestamp, err := strconv.Atoi(quoteData[4])
		if err != nil {
			return nil, ErrInvalidTimestampQuoteData
		}

		return &models.SymbolQuote{
			Symbol:    symbol,
			Bid:       bid,
			Ask:       ask,
			Last:      last,
			Timestamp: int64(timestamp),
		}, nil
	} else {
		return nil, ErrInvalidQuoteData
	}
}
