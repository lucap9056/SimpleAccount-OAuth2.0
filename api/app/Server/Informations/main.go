package Informations

import (
	"encoding/json"
	"net/http"
	"simple_account_oauth/app/Error"
	"simple_account_oauth/app/Server/Message"
)

func Handler(context *Message.Context) {
	var result string
	var errCode int
	var err error

	switch context.Request.Method {
	case http.MethodGet:
		result, errCode, err = GetAppInformation(context)
	case http.MethodPost:
		result, errCode, err = SetAppInformation(context)
	case http.MethodDelete:
		result, errCode, err = DeleteApp(context)
	default:
		errCode = Error.CLIENT_INVALID_REQUEST
	}

	if err != nil {
		context.Logger.Error.Write(err)
	}

	response := Message.Response{
		Success: errCode == Error.NULL,
		Result:  result,
		Error:   errCode,
	}

	responseBytes, _ := json.Marshal(response)
	writer := context.Writer

	if response.Success {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}

	writer.Header().Set("Content-Type", "text/json")
	writer.Write(responseBytes)
}
