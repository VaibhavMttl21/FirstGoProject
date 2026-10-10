package repository

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductAlreadyexists = errors.New("product already exists")