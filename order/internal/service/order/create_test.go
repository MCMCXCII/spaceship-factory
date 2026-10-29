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
		missingUUID = uuid.New()
		shieldUUID  = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		partsInStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", Price: 300000, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", Price: 300000, StockQuantity: 0},
		}

		partsInStockWithOptional = []model.Part{
			{UUID: hullUUID, Name: "Hull", Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", Price: 300000, StockQuantity: 3},
			{UUID: shieldUUID, Name: "Shield", Price: 400000, StockQuantity: 5},
		}
	)

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
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsInStock, nil)

				repo.EXPECT().Create(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.UUID != uuid.Nil &&
						o.HullUUID == hullUUID &&
						o.EngineUUID == engineUUID &&
						!o.ShieldUUID.Valid &&
						!o.WeaponUUID.Valid &&
						o.TotalPrice == 800000 && // 500000 + 300000
						o.Status == model.OrderStatusPendingPayment &&
						!o.CreatedAt.IsZero()
				})).Return(nil)
			},
			expected: expected{
				totalPrice: 800000,
			},
		},
		{
			name: "деталь не найдена",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: missingUUID,
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().ListParts(ctx, []uuid.UUID{hullUUID, missingUUID}).Return([]model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{
				err: errs.ErrPartNotFound,
			},
		},
		{
			name: "деталь закончилась на складе",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).Return(partsOutOfStock, nil)
			},
			expected: expected{
				err: errs.ErrOutOfStock,
			},
		},
		{
			name: "успешное создание заказа с опциональными полями",
			args: args{in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
				ShieldUUID: uuid.NullUUID{UUID: shieldUUID, Valid: true},
			}},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().ListParts(ctx, []uuid.UUID{hullUUID, engineUUID, shieldUUID}).Return(partsInStockWithOptional, nil)
				repo.EXPECT().Create(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.UUID != uuid.Nil &&
						o.HullUUID == hullUUID &&
						o.EngineUUID == engineUUID &&
						o.ShieldUUID.Valid &&
						o.ShieldUUID.UUID == shieldUUID &&
						!o.WeaponUUID.Valid &&
						o.TotalPrice == 1200000 && // 500000 + 300000 + 400000
						o.Status == model.OrderStatusPendingPayment &&
						!o.CreatedAt.IsZero()
				})).Return(nil)
			},
			expected: expected{
				totalPrice: 1200000,
				err:        nil,
			},
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
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(ctx, mock.Anything).
					Return(errRepo)
			},
			expected: expected{
				err: errRepo,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inventoryClient := mocks.NewInventoryClient(t)
			orderRepo := mocks.NewOrderRepository(t)

			tc.setupMock(orderRepo, inventoryClient)

			svc := New(orderRepo, inventoryClient, nil)

			got, err := svc.Create(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected.totalPrice, got.TotalPrice)
			assert.Equal(t, model.OrderStatusPendingPayment, got.Status)
		})
	}
}
