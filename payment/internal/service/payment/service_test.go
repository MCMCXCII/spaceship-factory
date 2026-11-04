package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/payment/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/service/input"
)

func TestPay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orderUUID := uuid.New()

	tests := []struct {
		name        string
		in          input.PayOrderInput
		expectedErr error
	}{
		{
			name: "оплата картой",
			in:   input.PayOrderInput{OrderUUID: orderUUID, PaymentMethod: model.PaymentMethodCard},
		},
		{
			name: "оплата через СБП",
			in:   input.PayOrderInput{OrderUUID: orderUUID, PaymentMethod: model.PaymentMethodSBP},
		},
		{
			name:        "способ оплаты не указан",
			in:          input.PayOrderInput{OrderUUID: orderUUID, PaymentMethod: model.PaymentMethodUnspecified},
			expectedErr: errs.ErrInvalidPaymentMethod,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := New()

			got, err := svc.Pay(ctx, tc.in)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				assert.Equal(t, uuid.Nil, got)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, got)
		})
	}
}
