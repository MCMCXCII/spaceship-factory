package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
)

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		// 400: ошибки разбора запроса от ogen
		{
			name:         "ogen: некорректный параметр пути",
			err:          &ogenerrors.DecodeParamsError{Err: errors.New("invalid uuid")},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ogen: некорректное тело запроса",
			err:          &ogenerrors.DecodeRequestError{Err: errors.New("unexpected EOF")},
			expectedCode: http.StatusBadRequest,
		},

		// 404
		{
			name:         "заказ не найден",
			err:          errs.ErrOrderNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "деталь не найдена",
			err:          errs.ErrPartNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "обёрнутая ошибка: заказ не найден",
			err:          fmt.Errorf("получить заказ: %w", errs.ErrOrderNotFound),
			expectedCode: http.StatusNotFound,
		},

		// 409
		{
			name:         "заказ уже оплачен",
			err:          errs.ErrOrderAlreadyPaid,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "заказ отменён",
			err:          errs.ErrOrderCancelled,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "обёрнутая ошибка: деталь отсутствует на складе",
			err:          fmt.Errorf("деталь Плазменный корпус: %w", errs.ErrOutOfStock),
			expectedCode: http.StatusConflict,
		},

		// 400: доменная валидация
		{
			name:         "неверный формат UUID",
			err:          errs.ErrInvalidUUID,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "неверный способ оплаты",
			err:          errs.ErrInvalidPaymentMethod,
			expectedCode: http.StatusBadRequest,
		},

		// 500
		{
			name:         "неизвестная ошибка",
			err:          errors.New("connection refused"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "обёрнутая неизвестная ошибка",
			err:          fmt.Errorf("оплатить заказ: %w", errors.New("context deadline exceeded")),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)

			ErrorHandler(context.Background(), rec, req, tc.err)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body errorResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

			assert.Equal(t, tc.expectedCode, body.Code)

			if tc.expectedCode == http.StatusInternalServerError {
				// детали внутренней ошибки не должны утекать клиенту
				assert.Equal(t, "внутренняя ошибка", body.Message)
				assert.NotContains(t, body.Message, tc.err.Error())
			} else {
				assert.Equal(t, tc.err.Error(), body.Message)
			}
		})
	}
}
