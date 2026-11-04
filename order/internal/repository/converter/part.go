package converter

import (
	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	var paymentMethod *string
	if o.PaymentMethod != nil {
		m := string(*o.PaymentMethod)
		paymentMethod = &m
	}

	return record.Order{
		UUID:            o.UUID,
		Status:          string(o.Status),
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   paymentMethod,
		CreatedAt:       o.CreatedAt,
	}
}

func OrderItemsToRecord(orderUUID uuid.UUID, items []model.OrderItem) []record.OrderItem {
	result := make([]record.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, record.OrderItem{
			OrderUUID: orderUUID,
			PartUUID:  item.PartUUID,
			PartType:  string(item.PartType),
			Price:     item.Price,
		})
	}
	return result
}

func OrderToModel(rec record.Order, items []record.OrderItem) model.Order {
	var paymentMethod *model.PaymentMethod
	if rec.PaymentMethod != nil {
		m := model.PaymentMethod(*rec.PaymentMethod)
		paymentMethod = &m
	}

	return model.Order{
		UUID:            rec.UUID,
		Items:           orderItemsToModel(items),
		TransactionUUID: rec.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(rec.Status),
		CreatedAt:       rec.CreatedAt,
	}
}

func orderItemsToModel(items []record.OrderItem) []model.OrderItem {
	result := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, model.OrderItem{
			PartUUID: item.PartUUID,
			PartType: model.PartType(item.PartType),
			Price:    item.Price,
		})
	}
	return result
}
