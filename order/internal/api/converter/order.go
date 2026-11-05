package converter

import (
	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/input"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

// CreateOrderRequestToInput переводит тело запроса на создание заказа во вход сервиса.
// Формат UUID и наличие обязательных полей уже проверил ogen.
func CreateOrderRequestToInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	return input.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: optNilUUIDToModel(req.ShieldUUID),
		WeaponUUID: optNilUUIDToModel(req.WeaponUUID),
	}
}

// PaymentMethodToModel переводит способ оплаты из HTTP API в доменный.
// Неизвестное значение — errs.ErrInvalidPaymentMethod.
func PaymentMethodToModel(m orderv1.PaymentMethod) (model.PaymentMethod, error) {
	switch m {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard, nil
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return "", errs.ErrInvalidPaymentMethod
	}
}

// OrderToCreateResponse собирает ответ 201 на создание заказа.
func OrderToCreateResponse(o model.Order) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  o.UUID,
		TotalPrice: o.TotalPrice(),
	}
}

// OrderToDTO переводит доменный заказ в DTO ответа GET.
// Детали раскладываются по слотам ответа по их типу.
func OrderToDTO(o model.Order) *orderv1.OrderDto {
	hullUUID, _ := itemUUID(o.Items, model.PartTypeHull)
	engineUUID, _ := itemUUID(o.Items, model.PartTypeEngine)

	return &orderv1.OrderDto{
		OrderUUID:       o.UUID,
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      optionalItemToDTO(o.Items, model.PartTypeShield),
		WeaponUUID:      optionalItemToDTO(o.Items, model.PartTypeWeapon),
		TotalPrice:      o.TotalPrice(),
		TransactionUUID: uuidPtrToDTO(o.TransactionUUID),
		PaymentMethod:   paymentMethodToDTO(o.PaymentMethod),
		Status:          orderStatusToDTO(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

// itemUUID ищет в заказе деталь указанного типа и возвращает её UUID.
// Если детали такого типа нет — uuid.Nil и false.
func itemUUID(items []model.OrderItem, partType model.PartType) (uuid.UUID, bool) {
	for _, item := range items {
		if item.PartType == partType {
			return item.PartUUID, true
		}
	}
	return uuid.Nil, false
}

// optionalItemToDTO возвращает UUID опциональной детали для ответа.
// Если детали такого типа в заказе нет — null.
func optionalItemToDTO(items []model.OrderItem, partType model.PartType) orderv1.OptNilUUID {
	id, ok := itemUUID(items, partType)
	if !ok {
		return nullUUID()
	}
	return orderv1.NewOptNilUUID(id)
}

// optNilUUIDToModel переводит опциональный UUID из запроса в указатель.
func optNilUUIDToModel(v orderv1.OptNilUUID) *uuid.UUID {
	id, ok := v.Get()
	if !ok {
		return nil
	}
	return &id
}

// uuidPtrToDTO переводит указатель на UUID в опциональное поле ответа.
// nil превращается в null.
func uuidPtrToDTO(v *uuid.UUID) orderv1.OptNilUUID {
	if v == nil {
		return nullUUID()
	}
	return orderv1.NewOptNilUUID(*v)
}

// nullUUID возвращает опциональный UUID со значением null.
func nullUUID() orderv1.OptNilUUID {
	var out orderv1.OptNilUUID
	out.SetToNull()
	return out
}

// paymentMethodToDTO переводит доменный способ оплаты в поле ответа.
// nil (заказ не оплачен) и неизвестные значения превращаются в null.
func paymentMethodToDTO(m *model.PaymentMethod) orderv1.OptNilPaymentMethod {
	var out orderv1.OptNilPaymentMethod
	if m == nil {
		out.SetToNull()
		return out
	}

	switch *m {
	case model.PaymentMethodCard:
		out.SetTo(orderv1.PaymentMethodCARD)
	case model.PaymentMethodSBP:
		out.SetTo(orderv1.PaymentMethodSBP)
	case model.PaymentMethodCreditCard:
		out.SetTo(orderv1.PaymentMethodCREDITCARD)
	case model.PaymentMethodInvestorMoney:
		out.SetTo(orderv1.PaymentMethodINVESTORMONEY)
	default:
		out.SetToNull()
	}

	return out
}

// orderStatusToDTO переводит доменный статус заказа в значение API.
func orderStatusToDTO(s model.OrderStatus) orderv1.OrderStatus {
	switch s {
	case model.OrderStatusPaid:
		return orderv1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderv1.OrderStatusCANCELLED
	default:
		return orderv1.OrderStatusPENDINGPAYMENT
	}
}
