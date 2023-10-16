package dto

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// DefaultResponse create a default response object
func DefaultResponse(codeDescription, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":    codeDescription,
		"message": message,
	}
}

// RenderResponse render http json response
func RenderResponse(ctx context.Context, writer http.ResponseWriter, httpStatusCode int, payload interface{}) {
	writer.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	writer.Header().Set("Content-Type", "application/json")

	marshal, err := json.Marshal(payload)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(httpStatusCode)
	_, _ = writer.Write(marshal)
}

// RenderErrorResponse render http error response
func RenderErrorResponse(ctx context.Context, writer http.ResponseWriter, httpStatusCode int, err error) {
	response := DefaultResponse(http.StatusText(httpStatusCode), err.Error())

	if errors, ok := err.(validator.ValidationErrors); ok {
		var details []map[string]interface{}
		for _, err := range errors {
			detail := map[string]interface{}{
				"field":       err.Field(),
				"value":       err.Value(),
				"location":    err.Namespace(),
				"issue":       err.Tag(),
				"description": err.Error(),
			}
			details = append(details, detail)
		}
		response["message"] = "Validation errors"
		response["details"] = details
	}

	RenderResponse(ctx, writer, httpStatusCode, response)
}
