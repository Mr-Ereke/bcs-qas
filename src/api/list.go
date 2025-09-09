package api

import (
	"net/http"

	"app/models"
	"app/service"

	"gitlab.online-fx.com/go-packages/apiresponse"
	"gitlab.online-fx.com/go-packages/gormdb"
	"gitlab.online-fx.com/go-packages/logger"
)

func List(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r) {
		return
	}

	var (
		subscription     models.Subscription
		subscriptionList []models.Subscription
	)

	rp := NewRequestParams(r)

	customerId := rp.GetUint("customerId")
	if customerId > 0 {
		subscription.CustomerId = customerId
	}

	instrument := rp.GetString("instrument", false)
	if instrument != "" {
		subscription.Instrument = instrument
	}

	if rp.Err() != nil {
		w.WriteHeader(http.StatusBadRequest)

		data := apiresponse.ResponseData{
			"code":   1,
			"result": rp.Err().Error(),
		}

		logHandlerError("alerts/list", "", &subscription, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/list")

		return
	}

	db := gormdb.GetClient(models.ServiceDB)
	err := db.Model(&models.Subscription{}).Where(&subscription).Order("id asc").Find(&subscriptionList).Error
	if err != nil {
		logger.Errorf("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validate := service.MaxAlertCountValid(&subscription)

	response := apiresponse.ResponseData{
		"code":             0,
		"result":           "Success",
		"data":             subscriptionList,
		"priceAlertsLimit": !validate.IsValid,
	}

	apiresponse.SendResponse(w, response, "alerts/list")
}
