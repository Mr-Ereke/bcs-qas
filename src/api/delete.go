package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"app/models"
	"app/service"

	"gitlab.online-fx.com/go-packages/apiresponse"
	"gitlab.online-fx.com/go-packages/gormdb"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	var (
		subscription     *models.Subscription
		subscriptionList []models.Subscription
		db               = gormdb.GetClient(models.ServiceDB)
	)

	if !validateHttpMethod(w, r) {
		return
	}

	requestData := struct {
		ID         int    `json:"id,omitempty"`
		CustomerId int    `json:"customerId,omitempty"`
		Instrument string `json:"instrument,omitempty"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		data := apiresponse.ResponseData{
			"code":   1,
			"result": "Decode error! please check your JSON formating!",
		}

		logHandlerError("alerts/delete", "", &requestData, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/delete")
		return
	}

	if requestData.ID > 0 {
		err = db.First(&subscription, requestData.ID).Error
		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			data := apiresponse.ResponseData{
				"code":   1,
				"result": "Not found subscription by ID: " + strconv.Itoa(requestData.ID),
			}

			logHandlerError("alerts/delete", "", &requestData, data["code"], data["result"])

			apiresponse.SendResponse(w, data, "alerts/delete")
			return
		}

		err = service.DeleteAlert(subscription)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			data := apiresponse.ResponseData{
				"code":   1,
				"result": fmt.Sprintf("Can not delete subscription. Error: %s", err),
			}

			logHandlerError("alerts/delete", "", &requestData, data["code"], data["result"])

			apiresponse.SendResponse(w, data, "alerts/delete")
			return
		}
	} else if requestData.Instrument != "" && requestData.CustomerId > 0 {
		err = db.Model(&models.Subscription{}).Where(&requestData).Find(&subscriptionList).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			logHandlerError("alerts/delete", "", &requestData, 1, "Internal Server Error")

			return
		}

		for _, alert := range subscriptionList {
			err = service.DeleteAlert(&alert)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				data := apiresponse.ResponseData{
					"code":   1,
					"result": fmt.Sprintf("Can not delete subscriptions. Error: %s", err),
				}

				logHandlerError("alerts/delete", "", &requestData, data["code"], data["result"])

				apiresponse.SendResponse(w, data, "alerts/delete")
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		data := apiresponse.ResponseData{
			"code":   1,
			"result": "Invalid request",
		}

		logHandlerError("alerts/delete", "", &requestData, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/delete")
		return
	}

	data := apiresponse.ResponseData{
		"code":   0,
		"result": "Success",
	}

	apiresponse.SendResponse(w, data, "alerts/delete")
}
