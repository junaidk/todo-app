package api

import (
	"encoding/json"
	"net/http"
)

func NewError(w http.ResponseWriter, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	msg, _ := json.Marshal(er)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(msg)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
