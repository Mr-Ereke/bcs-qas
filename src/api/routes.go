package api

import (
	"net/http"

	"gitlab.online-fx.com/go-packages/apiresponse"
)

var routeMethods = map[string]string{
	"/alerts/list":   http.MethodGet,
	"/alerts/add":    http.MethodPost,
	"/alerts/delete": http.MethodPost,
}

var Routes = map[string]func(http.ResponseWriter, *http.Request){
	"/alerts/list":   List,
	"/alerts/add":    Add,
	"/alerts/delete": Delete,
}

func validateHttpMethod(w http.ResponseWriter, r *http.Request) bool {
	valid := true

	if method, exists := routeMethods[r.URL.Path]; exists {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			data := apiresponse.ResponseData{
				"code":   1,
				"result": "Method " + r.Method + " not allowed",
			}
			apiresponse.SendResponse(w, data, r.URL.Path)
			valid = false
		}
	}

	return valid
}
