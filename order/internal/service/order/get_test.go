package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		missingUUID = uuid.New()
		errRepo     = errors.New("ошибка хранилища")

		order = model.Order{
			UUID:       orderUUID,
			HullUUID:   uuid.New(),
			EngineUUID: uuid.New(),
			TotalPrice: 800000,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC),
		}
	)

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		order model.Order
		err   error
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "заказ найден",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
			},
			expected: expected{
				order: order,
			},
		},
		{
			name: "заказ не найден",
			args: args{orderUUID: missingUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, missingUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка хранилища",
			args: args{orderUUID: orderUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, errRepo)
			},
			expected: expected{
				err: errRepo,
			},
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

			got, err := svc.Get(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected.order, got)
		})
	}
}
