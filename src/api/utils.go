package api

import (
	"encoding/json"
	"fmt"

	"gitlab.online-fx.com/go-packages/logger"
)

// logHandlerError логгирует ошибки, возникшие при обращении к ручкам.
func logHandlerError(handler string, function string, request any, code any, result any) {
	functionLog := ""
	if function != "" {
		functionLog = fmt.Sprintf("func: %s\n", function)
	}

	requestData, _ := json.Marshal(request)
	logger.Warningf("Failure of '%s' handler\nrequest: %s\n%scode: %d\nresult: %s",
		handler,
		string(requestData),
		functionLog,
		code,
		result)
}
