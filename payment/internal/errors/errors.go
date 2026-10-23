package errs

import "errors"

var (
	ErrInvalidUUID          = errors.New("неверный формат uuid")
	ErrInvalidPaymentMethod = errors.New("неверный способ оплаты")
)
