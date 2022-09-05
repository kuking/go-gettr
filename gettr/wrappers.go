package gettr

import (
	"errors"
	"fmt"
	"net/http"
)

type resultData struct {
	Data resultDataAux `json:"result"`
}
type resultDataAux struct {
	Data interface{}     `json:"Data"`
	Aux  resultAuxiliary `json:"Aux"`
}

type resultAuxiliary struct {
	Users  map[string]User `json:"uinf"`
	Cursor interface{}     `json:"Cursor"`
}

type resultLogin struct {
	Result resultLoginPayload `json:"result"`
}

type aPIErrorWrap struct {
	Payload APIError `json:"error"`
}

// APIError is a result object when the API call fails
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"emsg"`
}

type resultLoginPayload struct {
	User   User   `json:"user"`
	Token  string `json:"token"`
	Rtoken string `json:"rtoken"`
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
