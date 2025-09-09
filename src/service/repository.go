package service

import (
	"fmt"

	"app/models"

	"gitlab.online-fx.com/go-packages/logger"
)

func InsertAlert(subscription *models.Subscription) error {
	err := prepareSubscription(subscription)

	if err != nil {
		return err
	}

	err = db.Create(&subscription).Error
	if err != nil {
		return err
	}

	logger.Infof("Insert alert. Subscription: %+v", subscription)

	err = createSubscription(subscription)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAlert(subscription *models.Subscription) error {
	err := removeSubscription(subscription)
	if err != nil {
		return err
	}

	logger.Infof("Delete alert. Subscription: %+v", subscription)

	err = db.Delete(&subscription).Error
	if err != nil {
		return err
	}

	return nil
}

func prepareSubscription(subscription *models.Subscription) error {
	quote, err := redisClient.GetQuote(subscription.Instrument)

	if err != nil {
		return fmt.Errorf("failed get qoute. Instrument: %s. Error: %s", subscription.Instrument, err)
	}

	switch subscription.PriceType {
	case models.Bid:
		subscription.InitialPrice = quote.Bid
	case models.Ask:
		subscription.InitialPrice = quote.Ask
	case models.Last:
		subscription.InitialPrice = quote.Last
	default:
		return fmt.Errorf("unknow price type for inital price. Type: %s", subscription.PriceType)
	}

	if subscription.InitialPrice < subscription.Price {
		subscription.PriceDirection = true
	} else if subscription.InitialPrice > subscription.Price {
		subscription.PriceDirection = false
	} else {
		return fmt.Errorf("the current quote price is indicated. Select the price above or below")
	}

	return nil
}
