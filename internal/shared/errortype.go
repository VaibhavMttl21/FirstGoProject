package shared

import "errors"

var (
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrRequestTooLarge      = errors.New("request body too large")
	ErrEmptyRequestBody     = errors.New("empty request body")
	ErrInvalidJSON          = errors.New("invalid JSON")
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyexists = errors.New("product already exists")
	ErrProductIsEmpty       = errors.New("product name is empty")
	ErrProductPriceNegative = errors.New("product price cannot be negative")
	ErrStockIsNegative      = errors.New("stock quantity cannot be negative")
	ErrCurrencyInvalid      = errors.New("currency is invalid")
)
