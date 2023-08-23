package dto

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
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
