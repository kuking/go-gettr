package gettr

import (
	"errors"
	"fmt"
	"net/http"
)

type result struct {
	Data resultData `json:"result"`
}
type resultData struct {
	Data interface{}     `json:"Data"`
	Aux  resultAuxiliary `json:"Aux"`
}

type resultAuxiliary struct {
	Users  map[string]User `json:"uinf"`
	Cursor interface{}     `json:"Cursor"`
}

type aPIErrorWrap struct {
	Payload APIError `json:"error"`
}

// APIError is a result object when the API call fails
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"emsg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("gettr: [%v] %v", e.Code, e.Message)
}

func relevantError(response *http.Response, httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if response != nil {
		if response.StatusCode == 404 {
			return errors.New("not found")
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return errors.New("Http status " + response.Status)
		}
	}
	if apiError.Code == "" {
		return nil
	}
	return apiError
}
