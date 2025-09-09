package api

import (
	"fmt"
	"net/http"
	"strconv"
)

type RequestParams struct {
	r   *http.Request
	err error
}

func NewRequestParams(r *http.Request) *RequestParams {
	return &RequestParams{
		r:   r,
		err: nil,
	}
}

func (rp *RequestParams) setErr(err error) {
	if rp.err == nil {
		rp.err = err
	}
}

func (rp *RequestParams) Err() error {
	return rp.err
}

func (rp *RequestParams) GetValueFromRequestData(name string) (string, error) {
	values := rp.r.URL.Query()

	if values.Has(name) {
		if values.Get(name) != "" {
			return values.Get(name), nil
		} else {
			return "", fmt.Errorf("'%s' parameter is empty", name)
		}
	} else {
		return "", fmt.Errorf("'%s' parameter is required", name)
	}
}

func (rp *RequestParams) GetString(name string, required bool) string {
	valueString, err := rp.GetValueFromRequestData(name)

	if required && err != nil {
		rp.setErr(err)
	}

	return valueString
}

func (rp *RequestParams) GetUint(name string) uint {
	valueString, err := rp.GetValueFromRequestData(name)

	if err != nil {
		rp.setErr(err)
	}

	value, convertErr := strconv.Atoi(valueString)

	if convertErr != nil {
		rp.setErr(err)
		return 0
	}

	return uint(value)
}
