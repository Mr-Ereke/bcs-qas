package models

import (
	"net/http"

	"gitlab.online-fx.com/go-packages/apiresponse"
)

type ValidateResponse struct {
	IsValid        bool
	ResponseData   apiresponse.ResponseData
	HttpStatusCode int
}

func GetValidResponse() *ValidateResponse {
	return &ValidateResponse{
		IsValid:        true,
		ResponseData:   apiresponse.ResponseData{},
		HttpStatusCode: http.StatusOK,
	}
}
