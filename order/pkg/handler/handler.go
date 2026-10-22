package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// orderStore — хранилище заказов (in-memory)
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// handler реализует интерфейс orderv1.Handler, сгенерированный ogen
type handler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт новый обработчик заказов
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *handler {
	return &handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика
func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации)
// GET /api/v1/orders/{order_uuid}.
func (h *handler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	if req.EngineUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid обязателен",
		}, nil
	}

	if req.HullUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid обязателен",
		}, nil
	}

	var shieldPtr, weaponPtr *uuid.UUID

	uuids := make([]string, 0, 4)
	uuids = append(uuids, req.HullUUID.String(), req.EngineUUID.String())
	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		uuids = append(uuids, shieldUUID.String())
		shieldPtr = &shieldUUID
	}

	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		uuids = append(uuids, weaponUUID.String())
		weaponPtr = &weaponUUID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: st.Message()}, nil
			case codes.InvalidArgument:
				return &orderv1.CreateOrderBadRequest{Code: http.StatusBadRequest, Message: st.Message()}, nil
			}
		}
		slog.ErrorContext(ctx, "ошибка вызова InventoryService", "error", err)
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка на стороне inventory service",
		}, nil
	}

	var totalPrice int64
	for _, p := range resp.GetParts() {
		if p.GetStockQuantity() <= 0 {
			return &orderv1.CreateOrderConflict{Code: http.StatusConflict, Message: "нет детали на складе" + p.GetName()}, nil
		}
		totalPrice += p.GetPrice()
	}

	orderUUID := uuid.New()
	order := Order{
		OrderUUID:  orderUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldPtr,
		WeaponUUID: weaponPtr,
		TotalPrice: totalPrice,
		Status:     OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	h.store.mu.Lock()
	h.store.orders[orderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CreateOrderResponse{OrderUUID: orderUUID, TotalPrice: totalPrice}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.PayOrderNotFound{Code: http.StatusNotFound, Message: "заказ не найден"}, nil
	}

	if order.Status != OrderStatusPendingPayment {
		return &orderv1.PayOrderConflict{Code: http.StatusConflict, Message: "заказ уже оплачен или отменён"}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	payResp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: toProtoPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return &orderv1.PayOrderBadRequest{Code: http.StatusBadRequest, Message: st.Message()}, nil
		}
		slog.ErrorContext(ctx, "ошибка вызова PaymentService", "error", err)
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	txUUID, err := uuid.Parse(payResp.GetTransactionUuid())
	if err != nil {
		slog.ErrorContext(
			ctx, "PaymentService вернул некорректный transaction_uuid",
			"transaction_uuid", payResp.GetTransactionUuid(),
			"error", err,
		)
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	paymentMethod := PaymentMethod(req.PaymentMethod)

	h.store.mu.Lock()
	current, ok := h.store.orders[params.OrderUUID]
	if !ok || current.Status != OrderStatusPendingPayment {
		h.store.mu.Unlock()
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен или отменён",
		}, nil
	}

	current.Status = OrderStatusPaid
	current.TransactionUUID = &txUUID
	current.PaymentMethod = &paymentMethod
	h.store.orders[params.OrderUUID] = current
	h.store.mu.Unlock()

	return &orderv1.PayOrderResponse{TransactionUUID: txUUID}, nil
}

// toProtoPaymentMethod преобразует способ оплаты из HTTP API в proto-enum
func toProtoPaymentMethod(m orderv1.PaymentMethod) paymentv1.PaymentMethod {
	switch m {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *handler) CancelOrder(_ context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	order, ok := h.store.orders[params.OrderUUID]

	if !ok {
		return &orderv1.CancelOrderNotFound{Code: http.StatusNotFound, Message: "заказ не найден"}, nil
	}

	if order.Status != OrderStatusPendingPayment {
		return &orderv1.CancelOrderConflict{Code: http.StatusConflict, Message: "заказ уже оплачен или отменён"}, nil
	}

	order.Status = OrderStatusCancelled
	h.store.orders[params.OrderUUID] = order

	return &orderv1.CancelOrderResponse{}, nil
}

// func (h *handler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
//     // 1. Найти заказ в store
//     // 2. Проверить статус == PENDING_PAYMENT
//     // 3. Обновить статус на CANCELLED
//     // 4. Вернуть success
// }
