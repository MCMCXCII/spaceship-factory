package errs

import "errors"

var (
	ErrInvalidOrderUUID     = errors.New("неверный формат uuid")
	ErrInvalidPaymentMethod = errors.New("неверный способ оплаты")
)
