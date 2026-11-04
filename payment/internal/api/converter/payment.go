package converter

import (
	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/payment/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/service/input"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToInput(req *paymentv1.PayOrderRequest) (input.PayOrderInput, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return input.PayOrderInput{}, errs.ErrInvalidOrderUUID
	}

	return input.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: paymentMethodToModel(req.GetPaymentMethod()),
	}, nil
}

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
