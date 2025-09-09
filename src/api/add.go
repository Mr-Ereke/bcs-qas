package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"app/models"
	"app/service"

	"gitlab.online-fx.com/go-packages/apiresponse"
)

func Add(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription

	if !validateHttpMethod(w, r) {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&subscription)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		data := apiresponse.ResponseData{
			"code":   1,
			"result": "Decode error! please check your JSON formating!",
		}

		logHandlerError("alerts/add", "", &subscription, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/add")
		return
	}

	invalidParams := service.FieldValidate(&subscription)

	if len(invalidParams) > 0 {
		w.WriteHeader(http.StatusBadRequest)

		data := apiresponse.ResponseData{
			"code":   1,
			"result": fmt.Sprintf("Invalid params: %s", strings.Join(invalidParams, ", ")),
		}

		logHandlerError("alerts/add", "", &subscription, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/add")
		return
	}

	validate := service.DigitsValid(&subscription)
	if !validate.IsValid {
		w.WriteHeader(validate.HttpStatusCode)

		logHandlerError("alerts/add", "service.DigitsValid", &subscription, validate.ResponseData["code"], validate.ResponseData["result"])

		apiresponse.SendResponse(w, validate.ResponseData, "alerts/add")
		return
	}

	validate = service.MaxAlertInstrumentCountValid(&subscription)
	if !validate.IsValid {
		w.WriteHeader(validate.HttpStatusCode)

		logHandlerError("alerts/add", "service.MaxAlertInstrumentCountValid", &subscription, validate.ResponseData["code"], validate.ResponseData["result"])

		apiresponse.SendResponse(w, validate.ResponseData, "alerts/add")
		return
	}

	validate = service.MaxAlertCountValid(&subscription)
	if !validate.IsValid {
		w.WriteHeader(validate.HttpStatusCode)

		logHandlerError("alerts/add", "service.MaxAlertCountValid", &subscription, validate.ResponseData["code"], validate.ResponseData["result"])

		apiresponse.SendResponse(w, validate.ResponseData, "alerts/add")
		return
	}

	validate = service.DuplicateValid(&subscription)
	if !validate.IsValid {
		w.WriteHeader(validate.HttpStatusCode)

		logHandlerError("alerts/add", "service.DuplicateValid", &subscription, validate.ResponseData["code"], validate.ResponseData["result"])

		apiresponse.SendResponse(w, validate.ResponseData, "alerts/add")
		return
	}

	err = service.InsertAlert(&subscription)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data := apiresponse.ResponseData{
			"code":   1,
			"result": fmt.Sprintf("Can not add subscription. Error: %s", err),
		}

		logHandlerError("alerts/add", "service.InsertAlert", &subscription, data["code"], data["result"])

		apiresponse.SendResponse(w, data, "alerts/add")
		return
	}

	data := apiresponse.ResponseData{
		"code":   0,
		"result": "Success",
	}

	apiresponse.SendResponse(w, data, "alerts/add")
}
