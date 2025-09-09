package service

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"app/bellhop"
	"app/models"

	"gitlab.online-fx.com/go-packages/logger"
)

const (
	up   = "up +"
	down = "down -"
)

func sendBellhop(doneAlerts []*models.Subscription, direction bool) {
	if len(doneAlerts) == 0 {
		return
	}

	if len(doneAlerts) == 1 {
		push(doneAlerts[0])
		return
	}

	sort.Slice(doneAlerts, func(i, j int) bool {
		if direction {
			return doneAlerts[i].Price < doneAlerts[j].Price
		} else {
			return doneAlerts[i].Price > doneAlerts[j].Price
		}
	})

	if push(doneAlerts[len(doneAlerts)-1]) {
		for _, alert := range doneAlerts[:len(doneAlerts)-1] {
			err := DeleteAlert(alert)
			if err != nil {
				logger.Errorf("Delete alert. Error: %s", err)
			}
		}
	}
}

func push(subscription *models.Subscription) bool {
	err := bellhop.SendPush(subscription.CustomerId, generateBody(subscription), subscription.Instrument, subscription.QuoteTimestamp)

	if err != nil {
		return false
	}

	err = DeleteAlert(subscription)
	if err != nil {
		logger.Errorf("Delete alert. Error: %s", err)
	}

	return true
}

// Example GOLD is up +1,93% to R170,98
func generateBody(subscription *models.Subscription) string {
	body := strings.Builder{}
	body.WriteString(subscription.DisplayTitle + " is ")
	body.WriteString(getDirectionString(subscription.PriceDirection))
	body.WriteString(getPercentString(subscription) + "% to ")
	body.WriteString(getCurrencyPrice(subscription))

	return body.String()
}

func getCurrencyPrice(subscription *models.Subscription) string {
	return fmt.Sprintf("%s%g", subscription.Currency, subscription.Price)
}

func getDirectionString(direction bool) string {
	if direction {
		return up
	} else {
		return down
	}
}

func getPercentString(subscription *models.Subscription) string {
	var (
		different, percent float64
	)

	if subscription.PriceDirection {
		different = subscription.Price - subscription.InitialPrice
	} else {
		different = subscription.InitialPrice - subscription.Price
	}

	percent = different * 100 / subscription.InitialPrice

	percent = math.Round(percent*100) / 100

	return fmt.Sprintf("%g", percent)
}
