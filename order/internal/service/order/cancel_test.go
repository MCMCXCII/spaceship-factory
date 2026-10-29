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

func TestCancel(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		missingUUID = uuid.New()
		errRepo     = errors.New("ошибка хранилища")

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
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешная отмена",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID == orderUUID &&
							o.Status == model.OrderStatusCancelled &&
							!o.TransactionUUID.Valid &&
							o.PaymentMethod == model.PaymentMethodUnspecified &&
							o.TotalPrice == pendingOrder.TotalPrice
					})).
					Return(nil)
			},
		},
		{
			name: "заказ не найден",
			args: args{orderUUID: missingUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, missingUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "заказ уже оплачен",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(orderWithStatus(model.OrderStatusPaid), nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "заказ уже отменён",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(orderWithStatus(model.OrderStatusCancelled), nil)
			},
			expected: expected{err: errs.ErrOrderCancelled},
		},
		{
			name: "ошибка сохранения",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(pendingOrder, nil)
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

			tc.setupMock(orderRepo)

			svc := New(orderRepo, inventoryClient, paymentClient)

			err := svc.Cancel(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				return
			}

			require.NoError(t, err)
		})
	}
}
