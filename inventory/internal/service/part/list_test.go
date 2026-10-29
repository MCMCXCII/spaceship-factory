package inventory

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/inventory/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/input"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/part/mocks"
)

func TestList(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		hullAl = model.Part{UUID: uuid.New(), Name: "Алюминиевый корпус", PartType: model.PartTypeHull}
		hullTi = model.Part{UUID: uuid.New(), Name: "Титановый корпус", PartType: model.PartTypeHull}
		engine = model.Part{UUID: uuid.New(), Name: "Ионный двигатель", PartType: model.PartTypeEngine}
		shield = model.Part{UUID: uuid.New(), Name: "Энергетический щит", PartType: model.PartTypeShield}

		unsorted = []model.Part{shield, hullTi, engine, hullAl}

		missingUUID = uuid.New()
		errRepo     = errors.New("ошибка хранилища")
	)

	type args struct {
		in input.PartFilter
	}

	type expected struct {
		parts []model.Part
		err   error
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.InventoryRepository)
		expected  expected
	}{
		{
			name: "по UUID: детали в порядке запроса",
			args: args{in: input.PartFilter{
				UUIDs: []uuid.UUID{shield.UUID, hullAl.UUID, engine.UUID},
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, shield.UUID).Return(shield, nil)
				repo.EXPECT().Get(ctx, hullAl.UUID).Return(hullAl, nil)
				repo.EXPECT().Get(ctx, engine.UUID).Return(engine, nil)
			},
			expected: expected{
				parts: []model.Part{shield, hullAl, engine},
			},
		},
		{
			name: "по UUID: нет детали",
			args: args{in: input.PartFilter{
				UUIDs: []uuid.UUID{shield.UUID, missingUUID, engine.UUID},
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, shield.UUID).Return(shield, nil)
				repo.EXPECT().Get(ctx, missingUUID).Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{
				parts: nil,
				err:   errs.ErrPartNotFound,
			},
		},
		{
			name: "по UUID: тип в фильтре игнорируется",
			args: args{in: input.PartFilter{
				UUIDs:    []uuid.UUID{hullAl.UUID, engine.UUID},
				PartType: model.PartTypeWeapon,
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, hullAl.UUID).Return(hullAl, nil)
				repo.EXPECT().Get(ctx, engine.UUID).Return(engine, nil)
			},
			expected: expected{
				parts: []model.Part{hullAl, engine},
			},
		},
		{
			name: "по типу: PartTypeHull",
			args: args{in: input.PartFilter{
				PartType: model.PartTypeHull,
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().List(ctx).Return(unsorted, nil)
			},
			expected: expected{
				parts: []model.Part{hullAl, hullTi},
			},
		},
		{
			name: "по типу: PartTypeUnspecified",
			args: args{in: input.PartFilter{
				PartType: model.PartTypeUnspecified,
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().List(ctx).Return(unsorted, nil)
			},
			expected: expected{
				parts: []model.Part{hullAl, engine, hullTi, shield},
			},
		},
		{
			name: "по типу: ошибка репозитория",
			args: args{in: input.PartFilter{
				PartType: model.PartTypeHull,
			}},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().List(ctx).Return(nil, errRepo)
			},
			expected: expected{
				err: errRepo,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inventoryRepo := mocks.NewInventoryRepository(t)
			tc.setupMock(inventoryRepo)

			svc := New(inventoryRepo)

			got, err := svc.List(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected.parts, got)
		})
	}
}
