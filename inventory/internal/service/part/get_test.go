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
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		partUUID    = uuid.New()
		part        = model.Part{UUID: partUUID, Name: "Ионный двигатель класса D", PartType: model.PartTypeEngine, Price: 64000}
		missingUUID = uuid.New()
		errRepo     = errors.New("ошибка хранилища")
	)

	type expected struct {
		part model.Part
		err  error
	}

	type args struct {
		partUUID uuid.UUID
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.InventoryRepository)
		expected  expected
	}{
		{
			name: "деталь вернулась успешно",
			args: args{
				partUUID: partUUID,
			},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(part, nil)
			},
			expected: expected{
				part: part,
				err:  nil,
			},
		},
		{
			name: "деталь не найдена",
			args: args{
				partUUID: missingUUID,
			},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, missingUUID).Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{
				part: model.Part{},
				err:  errs.ErrPartNotFound,
			},
		},
		{
			name: "ошибка репозитория",
			args: args{
				partUUID: partUUID,
			},
			setupMock: func(repo *mocks.InventoryRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(model.Part{}, errRepo)
			},
			expected: expected{
				part: model.Part{},
				err:  errRepo,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inventoryRepo := mocks.NewInventoryRepository(t)

			tc.setupMock(inventoryRepo)

			svc := New(inventoryRepo)

			got, err := svc.Get(ctx, tc.args.partUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected.part, got)
		})
	}
}
