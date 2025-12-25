package api_service

import "errors"

var (
	ErrDeletedConflict = errors.New("order already deleted")
	ErrNotFound        = errors.New("order not found")
	ErrResultNotReady  = errors.New("order result not found in redis")
)
