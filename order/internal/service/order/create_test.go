package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/input"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		hullUUID    = uuid.New()
		engineUUID  = uuid.New()
		shieldUUID  = uuid.New()
		missingUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		hull   = model.Part{UUID: hullUUID, Name: "Корпус", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10}
		engine = model.Part{UUID: engineUUID, Name: "Двигатель", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 5}
		shield = model.Part{UUID: shieldUUID, Name: "Щит", PartType: model.PartTypeShield, Price: 400000, StockQuantity: 3}

		engineOutOfStock = model.Part{UUID: engineUUID, Name: "Двигатель", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 0}
	)

	// itemsMatch проверяет, что строки заказа совпадают с деталями по UUID, типу и цене, в том же порядке
	itemsMatch := func(items []model.OrderItem, parts ...model.Part) bool {
		if len(items) != len(parts) {
			return false
		}
		for i, part := range parts {
			if items[i].PartUUID != part.UUID ||
				items[i].PartType != part.PartType ||
				items[i].Price != part.Price {
				return false
			}
		}
		return true
	}

	type args struct {
		in input.CreateOrderInput
	}

	type expected struct {
		totalPrice int64
		err        error
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{hull, engine}, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID != uuid.Nil &&
							itemsMatch(o.Items, hull, engine) &&
							o.TotalPrice() == 800000 && // 500000 + 300000
							o.Status == model.OrderStatusPendingPayment &&
							o.TransactionUUID == nil &&
							o.PaymentMethod == nil &&
							!o.CreatedAt.IsZero()
					})).
					Return(nil)
			},
			expected: expected{totalPrice: 800000},
		},
		{
			name: "успешное создание заказа с опциональным щитом",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
				ShieldUUID: &shieldUUID,
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID, shieldUUID}).
					Return([]model.Part{hull, engine, shield}, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID != uuid.Nil &&
							itemsMatch(o.Items, hull, engine, shield) &&
							o.TotalPrice() == 1200000 && // 500000 + 300000 + 400000
							o.Status == model.OrderStatusPendingPayment &&
							!o.CreatedAt.IsZero()
					})).
					Return(nil)
			},
			expected: expected{totalPrice: 1200000},
		},
		{
			name: "деталь не найдена",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: missingUUID,
			}},
			setupMock: func(_ *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, missingUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "деталь закончилась на складе",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			}},
			setupMock: func(_ *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{hull, engineOutOfStock}, nil)
			},
			expected: expected{err: errs.ErrOutOfStock},
		},
		{
			name: "ошибка сохранения заказа",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{hull, engine}, nil)

				repo.EXPECT().
					Create(ctx, mock.Anything).
					Return(errRepo)
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

			tc.setupMock(orderRepo, inventoryClient)

			svc := New(orderRepo, inventoryClient, paymentClient)

			got, err := svc.Create(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, got.UUID)
			assert.Equal(t, tc.expected.totalPrice, got.TotalPrice())
			assert.Equal(t, model.OrderStatusPendingPayment, got.Status)
		})
	}
}
