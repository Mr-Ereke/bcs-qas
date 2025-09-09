package bellhop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"app/metrics"
	"app/models"

	"gitlab.online-fx.com/go-packages/logger"
)

var bellhopConfig *models.Bellhop

const pushUrl = "/v1/send"

func InitConfig(config *models.Config) {
	bellhopConfig = &config.Bellhop
}

func SendPush(customerId uint, body string, symbol string, quoteTimestamp int64) error {
	requestTimeout := time.Second * time.Duration(bellhopConfig.Timeout)

	httpClient := http.Client{
		Timeout: requestTimeout,
	}

	requestData := models.BuildRequest(bellhopConfig.PushName, customerId, body, symbol)

	requestJson, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("json encoding error: %s", err)
	}

	request, newRequestErr := http.NewRequest(http.MethodPost, bellhopConfig.Host+pushUrl, bytes.NewBuffer(requestJson))

	if newRequestErr != nil {
		return fmt.Errorf("request init error: %s", newRequestErr)
	}

	request.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(request.Context(), requestTimeout)
	request = request.WithContext(ctx)
	defer cancel()

	logger.Infof("Pushing send to Bellhop. QuoteTime: %s. Request: %s", time.UnixMilli(quoteTimestamp), requestJson)

	infoMessage := metrics.SetPushMetrics(quoteTimestamp)

	if infoMessage != "" {
		logger.Infof("Request: %s. Metrics info: %s", requestJson, infoMessage)
	}

	sendMilli := time.Now().UnixMilli()

	response, requestErr := httpClient.Do(request)

	metrics.SetBellhopMetrics(sendMilli)

	if requestErr != nil {
		return fmt.Errorf("failed request to bellhop. Request: %s. Error: %s", requestJson, requestErr)
	}

	result, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("failed response read. Request: %s. Error: %s", requestJson, readErr)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			logger.Errorf("Fail close body. Error: %s", err)
		}
	}()

	logger.Infof("Success sent to Bellhop. Push: %s. Status: %s. Response: %s", requestJson, response.Status, result)

	return nil
}
