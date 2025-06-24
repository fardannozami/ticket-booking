package exception

import (
	"encoding/json"
	"net/http"

	"github.com/fardannozami/ticket-booking/helper"
	"github.com/fardannozami/ticket-booking/response"
)

func PanicHandler(writer http.ResponseWriter, req *http.Request, err interface{}) {
	internalServerError(writer, req, err)
}

func internalServerError(writer http.ResponseWriter, req *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	var message string

	switch e := err.(type) {
	case string:
		message = e
	case error:
		message = e.Error()
	default:
		message = "Unknown error"
	}

	apiResponse := response.ApiResponse{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Data:    message,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	errorEncode := json.NewEncoder(writer).Encode(&apiResponse)
	helper.PanicIfError(errorEncode)
}
