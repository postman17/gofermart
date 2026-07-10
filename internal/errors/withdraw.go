package errors

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidAmount = errors.New("amount must be greater than zero")
