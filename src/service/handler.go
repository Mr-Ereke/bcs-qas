package service

import (
	"app/models"

	"gitlab.online-fx.com/go-packages/logger"
)

func upHandler(quote *models.SymbolQuote, priceType string) {
	subscriptionList, err := GetSubscriptionList(quote.Symbol, priceType, true)

	if err != nil {
		logger.Errorf("Failed handle %s up quote. Error: %s", priceType, err)
		return
	}

	doneAlerts := make(map[uint][]*models.Subscription)

	for _, alert := range subscriptionList {
		if alert.Price <= quote.GetQuoteByType(priceType) {
			alert.QuoteTimestamp = quote.Timestamp
			doneAlerts[alert.CustomerId] = append(doneAlerts[alert.CustomerId], alert)
		}
	}

	for _, alertList := range doneAlerts {
		// send alert list group by customer
		go sendBellhop(alertList, true)
	}
}

func downHandler(quote *models.SymbolQuote, priceType string) {
	subscriptionList, err := GetSubscriptionList(quote.Symbol, priceType, false)

	if err != nil {
		logger.Errorf("Failed handle %s down quote. Error: %s", priceType, err)
		return
	}

	doneAlerts := make(map[uint][]*models.Subscription)

	for _, alert := range subscriptionList {
		if alert.Price >= quote.GetQuoteByType(priceType) {
			alert.QuoteTimestamp = quote.Timestamp
			doneAlerts[alert.CustomerId] = append(doneAlerts[alert.CustomerId], alert)
		}
	}

	for _, alertList := range doneAlerts {
		// send alert list group by customer
		go sendBellhop(alertList, false)
	}
}
