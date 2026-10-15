package handler

import (
	"encoding/json"
	"errors"
	"firstgoproject/internal/model"
	"firstgoproject/internal/shared"
	"fmt"
	"io"
	"mime"
	"net/http"
)

// decodeProduct decodes exactly one JSON object into a Product.
// It returns errors instead of writing HTTP responses itself.
func decodeProduct(
	w http.ResponseWriter,
	r *http.Request,
) (model.Product, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return model.Product{}, shared.ErrUnsupportedMediaType
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)

	var product model.Product
	if err := decoder.Decode(&product); err != nil {
		var maxBytesError *http.MaxBytesError

		if errors.As(err, &maxBytesError) {
			return model.Product{}, shared.ErrRequestTooLarge
		}

		if errors.Is(err, io.EOF) {
			return model.Product{}, shared.ErrEmptyRequestBody
		}

		return model.Product{}, fmt.Errorf("%w: %v", shared.ErrInvalidJSON, err)
	}

	// Reject a second JSON value or invalid trailing content.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		var maxBytesError *http.MaxBytesError

		if errors.As(err, &maxBytesError) {
			return model.Product{}, shared.ErrRequestTooLarge
		}

		return model.Product{}, shared.ErrInvalidJSON
	}

	return product, nil
}



func writeInputError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrUnsupportedMediaType):
		writeError(w, http.StatusUnsupportedMediaType,
			"UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")

	case errors.Is(err, shared.ErrRequestTooLarge):
		writeError(w, http.StatusRequestEntityTooLarge,
			"REQUEST_TOO_LARGE", "Request body exceeds the 1 MiB limit")

	case errors.Is(err, shared.ErrEmptyRequestBody):
		writeError(w, http.StatusBadRequest,
			"EMPTY_REQUEST_BODY", "Request body must not be empty")

	case errors.Is(err, shared.ErrInvalidJSON):
		writeError(w, http.StatusBadRequest,
			"INVALID_JSON", "Request body must contain exactly one valid JSON object")

	default:
		writeError(w, http.StatusInternalServerError,
			"INTERNAL_ERROR", "An unexpected error occurred")
	}
}
