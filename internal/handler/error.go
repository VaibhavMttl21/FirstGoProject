package handler

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"

	"firstgoproject/internal/shared"
)

type APIError struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// The response may already be committed.
		// Log the error through your server's logging infrastructure.
		return
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	writeJSON(w, status, APIError{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrProductIsEmpty):
		writeError(w, http.StatusBadRequest,
			"INVALID_PRODUCT_NAME", "Product name is required")

	case errors.Is(err, shared.ErrProductPriceNegative):
		writeError(w, http.StatusBadRequest,
			"INVALID_PRICE", "Product price cannot be negative")

	case errors.Is(err, shared.ErrStockIsNegative):
		writeError(w, http.StatusBadRequest,
			"INVALID_STOCK", "Stock quantity cannot be negative")

	case errors.Is(err, shared.ErrCurrencyInvalid):
		writeError(w, http.StatusBadRequest,
			"INVALID_CURRENCY", "Currency must contain three uppercase letters")

	case errors.Is(err, shared.ErrProductNotFound):
		writeError(w, http.StatusNotFound,
			"PRODUCT_NOT_FOUND", "Product not found")

	case errors.Is(err, shared.ErrProductAlreadyexists):
		writeError(w, http.StatusConflict,
			"PRODUCT_ALREADY_EXISTS", "Product already exists")

	case errors.Is(err, context.Canceled):
		writeError(w, 499,
			"REQUEST_CANCELED", "Request was canceled")

	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout,
			"REQUEST_TIMEOUT", "Request deadline exceeded")

	default:
		writeError(w, http.StatusInternalServerError,
			"INTERNAL_ERROR", "An unexpected error occurred")
	}
}
