package converter

import (
	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/payment/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// PayOrderRequestToModel извлекает из запроса UUID заказа и способ оплаты.
func PayOrderRequestToModel(req *paymentv1.PayOrderRequest) (uuid.UUID, model.PaymentMethod, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return uuid.Nil, model.PaymentMethodUnspecified, errs.ErrInvalidUUID
	}

	method := paymentMethodToModel(req.GetPaymentMethod())
	if method == model.PaymentMethodUnspecified {
		return uuid.Nil, model.PaymentMethodUnspecified, errs.ErrInvalidPaymentMethod
	}

	return orderUUID, method, nil
}

// paymentMethodToModel переводит способ оплаты из proto-enum в доменный.
func paymentMethodToModel(m paymentv1.PaymentMethod) model.PaymentMethod {
	switch m {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
