package contract

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/usagifm/product-service/utils/logger"
)

type Meta struct {
}
type Response struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Error   *Error      `json:"error"`
	Message string      `json:"message"`
}

func SuccessResponse(ctx context.Context, w http.ResponseWriter, message string, data interface{}, statusCode int) {
	resbody := Response{}
	resbody.Code = statusCode
	resbody.Data = data
	resbody.Message = message
	resbody.Success = true
	resp, err := json.MarshalIndent(resbody, "", "  ")

	if err != nil {
		logger.GetLogger(ctx).Error("json marshall ErrorResponse,  err: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)

}

func ErrorResponse(ctx context.Context, w http.ResponseWriter, message string, errorMessage string, errorTitle string, statusCode int) {
	resbody := Response{}
	resbody.Code = statusCode
	resbody.Message = message
	errors := Error{}
	errors.Message = errorMessage
	errors.Title = errorTitle
	resbody.Error = &errors
	resbody.Success = false

	resp, err := json.MarshalIndent(resbody, "", "  ")

	if err != nil {
		logger.GetLogger(ctx).Error("json marshall ErrorResponse,  err: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

type Error struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type EmptyObj struct{}
