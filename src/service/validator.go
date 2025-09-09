package service

import (
	"fmt"
	"net/http"
	"strings"

	"app/models"

	"gitlab.online-fx.com/go-packages/apiresponse"
)

func DuplicateValid(subscription *models.Subscription) *models.ValidateResponse {
	priceTypeDirection, instrumentExists := memory.Get(subscription.Instrument)

	if !instrumentExists {
		return models.GetValidResponse()
	}

	alertList := priceTypeDirection.GetSubscriptionListByCustomerId(subscription.CustomerId)

	for _, alert := range alertList {
		if alert.PriceType == subscription.PriceType && alert.Price == subscription.Price {
			return &models.ValidateResponse{
				ResponseData: apiresponse.ResponseData{
					"code":   3,
					"result": "This price is already set",
				},
				HttpStatusCode: http.StatusBadRequest,
			}
		}
	}

	return models.GetValidResponse()
}

func MaxAlertCountValid(subscription *models.Subscription) *models.ValidateResponse {
	customerCount := 0

	for _, priceTypeDirection := range memory.Alerts {
		customerCount += priceTypeDirection.GetSubscriptionCountByCustomerId(subscription.CustomerId)
	}

	if customerCount >= settings.MaxAlertsCustomer {
		return &models.ValidateResponse{
			ResponseData: apiresponse.ResponseData{
				"code":   2,
				"result": fmt.Sprintf("Maximum number of alerts for all characters: %d", settings.MaxAlertsCustomer),
			},
			HttpStatusCode: http.StatusTooManyRequests,
		}
	}

	return models.GetValidResponse()
}

func MaxAlertInstrumentCountValid(subscription *models.Subscription) *models.ValidateResponse {
	priceTypeDirection, instrumentExists := memory.Get(subscription.Instrument)

	if !instrumentExists {
		return models.GetValidResponse()
	}

	alertCount := priceTypeDirection.GetSubscriptionCountByCustomerId(subscription.CustomerId)

	if alertCount >= settings.MaxAlertsCustomerInstrument {
		return &models.ValidateResponse{
			ResponseData: apiresponse.ResponseData{
				"code":   1,
				"result": fmt.Sprintf("Maximum number of alerts: %d", settings.MaxAlertsCustomerInstrument),
			},
			HttpStatusCode: http.StatusTooManyRequests,
		}
	}

	return models.GetValidResponse()
}

func DigitsValid(subscription *models.Subscription) *models.ValidateResponse {
	priceStr := fmt.Sprintf("%g", subscription.Price)

	if !strings.Contains(priceStr, ".") {
		return models.GetValidResponse()
	}

	priceSlice := strings.Split(priceStr, ".")

	if len(priceSlice) == 2 {
		if len(priceSlice[1]) <= int(subscription.Digits) {
			return models.GetValidResponse()
		}
	}

	return &models.ValidateResponse{
		ResponseData: apiresponse.ResponseData{
			"code":   4,
			"result": "Wrong digits",
		},
		HttpStatusCode: http.StatusBadRequest,
	}
}

func FieldValidate(subscription *models.Subscription) []string {
	var invalidParams []string

	if subscription.CustomerId <= 0 {
		invalidParams = append(invalidParams, "customerId")
	}

	if subscription.Instrument == "" {
		invalidParams = append(invalidParams, "instrument")
	}

	if subscription.DisplayTitle == "" {
		invalidParams = append(invalidParams, "displayTitle")
	}

	if subscription.Price <= 0 {
		invalidParams = append(invalidParams, "price")
	}

	if subscription.PriceType == "" {
		invalidParams = append(invalidParams, "priceType")
	}

	if subscription.Currency == "" {
		invalidParams = append(invalidParams, "currency")
	}

	if subscription.Digits < 0 {
		invalidParams = append(invalidParams, "digits")
	}

	return invalidParams
}
