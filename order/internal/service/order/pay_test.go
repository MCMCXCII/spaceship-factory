package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/order/mocks"
)

func TestPay(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		missingUUID = uuid.New()
		txUUID      = uuid.New()
		method      = model.PaymentMethodCard

		errRepo    = errors.New("ошибка хранилища")
		errPayment = errors.New("payment недоступен")

		pendingOrder = model.Order{
			UUID:       orderUUID,
			HullUUID:   uuid.New(),
			EngineUUID: uuid.New(),
			TotalPrice: 800000,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC),
		}
	)

	// orderWithStatus возвращает копию ожидающего заказа с другим статусом
	orderWithStatus := func(status model.OrderStatus) model.Order {
		o := pendingOrder
		o.Status = status
		return o
	}

	type args struct {
		orderUUID uuid.UUID
		method    model.PaymentMethod
	}

	type expected struct {
		txUUID uuid.UUID
		err    error
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, payment *mocks.PaymentClient)
		expected  expected
	}{
		{
			name: "успешная оплата",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
				payment.EXPECT().PayOrder(ctx, orderUUID, method).Return(txUUID, nil)
				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID == orderUUID &&
							o.Status == model.OrderStatusPaid &&
							o.TransactionUUID.Valid &&
							o.TransactionUUID.UUID == txUUID &&
							o.PaymentMethod == method &&
							o.TotalPrice == pendingOrder.TotalPrice
					})).
					Return(nil)
			},
			expected: expected{txUUID: txUUID},
		},
		{
			name: "заказ не найден",
			args: args{orderUUID: missingUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, missingUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "заказ уже оплачен",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(orderWithStatus(model.OrderStatusPaid), nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "заказ отменён",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(orderWithStatus(model.OrderStatusCancelled), nil)
			},
			expected: expected{err: errs.ErrOrderCancelled},
		},
		{
			name: "payment отклонил способ оплаты",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
				payment.EXPECT().PayOrder(ctx, orderUUID, method).Return(uuid.Nil, errs.ErrInvalidPaymentMethod)
			},
			expected: expected{err: errs.ErrInvalidPaymentMethod},
		},
		{
			name: "payment недоступен",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
				payment.EXPECT().PayOrder(ctx, orderUUID, method).Return(uuid.Nil, errPayment)
			},
			expected: expected{err: errPayment},
		},
		{
			name: "ошибка сохранения после оплаты",
			args: args{orderUUID: orderUUID, method: method},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
				payment.EXPECT().PayOrder(ctx, orderUUID, method).Return(txUUID, nil)
				repo.EXPECT().Update(ctx, mock.Anything).Return(errRepo)
			},
			expected: expected{err: errRepo},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(orderRepo, paymentClient)

			svc := New(orderRepo, inventoryClient, paymentClient)

			got, err := svc.Pay(ctx, tc.args.orderUUID, tc.args.method)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected.txUUID, got)
		})
	}
}
