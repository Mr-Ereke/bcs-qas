package service

import (
	"fmt"

	"app/models"

	"gitlab.online-fx.com/go-packages/logger"
)

func GetSubscriptionList(symbol string, priceType string, direction bool) (map[uint]*models.Subscription, error) {
	priceTypeDirection, instrumentExists := memory.Get(symbol)

	if !instrumentExists {
		return map[uint]*models.Subscription{}, fmt.Errorf("not found alert by instrument: %s", symbol)
	}

	list, err := priceTypeDirection.Get(priceType, direction)

	return list, err
}

func createSubscription(subscription *models.Subscription) error {
	priceTypeDirection, instrumentExists := memory.Get(subscription.Instrument)

	if !instrumentExists {
		priceTypeDirection = models.NewPriceTypeDirection()
		memory.Set(subscription.Instrument, priceTypeDirection)
		SubscribeInstrument(subscription.Instrument)
	}

	return priceTypeDirection.Create(subscription)
}

func removeSubscription(subscription *models.Subscription) error {
	priceTypeDirection, instrumentExists := memory.Get(subscription.Instrument)

	if !instrumentExists {
		return fmt.Errorf("not found alert by instrument: %s", subscription.Instrument)
	}

	err := priceTypeDirection.Delete(subscription)

	if err != nil {
		return err
	}

	if priceTypeDirection.TypeExists.IsEmpty() {
		logger.Infof("Deleting %s instrument from memory ...", subscription.Instrument)
		memory.Delete(subscription.Instrument)
		UnsubscribeInstrument(subscription.Instrument)
	}

	return nil
}
