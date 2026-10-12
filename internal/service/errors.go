package service

import "errors"

var (
	ErrProductIsEmpty       = errors.New("product name is empty")
	ErrProductPriceNegative = errors.New("product price cannot be negative")
	ErrStockIsNegative      = errors.New("stock quantity cannot be negative")
	ErrCurrencyInvalid      = errors.New("currency is invalid")
)
