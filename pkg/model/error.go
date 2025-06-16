package model

import "net/http"

type ErrorResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func StatusInternalServerError(msg string) ErrorResponse {
	return ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

func StatusBadRequestError(msg string) ErrorResponse {
	return ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}
