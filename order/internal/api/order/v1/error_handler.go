package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
)

// errorResponse — тело ответа с ошибкой, формат из error.yaml
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler — глобальный обработчик ошибок ogen.
// Переводит доменные ошибки хендлеров в HTTP-коды. Подключается через orderv1.WithErrorHandler.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code, message := mapError(err)

	if code == http.StatusInternalServerError {
		slog.ErrorContext(ctx, "необработанная ошибка HTTP-запроса", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if encErr := json.NewEncoder(w).Encode(errorResponse{Code: code, Message: message}); encErr != nil {
		slog.ErrorContext(ctx, "ошибка кодирования ответа", "error", encErr)
	}
}

// mapError определяет HTTP-код и текст ответа по ошибке.
func mapError(err error) (int, string) {
	var (
		decodeParams  *ogenerrors.DecodeParamsError
		decodeRequest *ogenerrors.DecodeRequestError
	)

	switch {
	// Ошибки разбора запроса от ogen: кривой UUID в пути, нет обязательного поля,
	// неизвестный payment_method — всё это 400
	case errors.As(err, &decodeParams), errors.As(err, &decodeRequest):
		return http.StatusBadRequest, err.Error()

	// 404 Not Found
	case errors.Is(err, errs.ErrOrderNotFound),
		errors.Is(err, errs.ErrPartNotFound):
		return http.StatusNotFound, err.Error()

	// 409 Conflict
	case errors.Is(err, errs.ErrOrderAlreadyPaid),
		errors.Is(err, errs.ErrOrderCancelled),
		errors.Is(err, errs.ErrOutOfStock):
		return http.StatusConflict, err.Error()

	// 400 Bad Request
	case errors.Is(err, errs.ErrInvalidUUID),
		errors.Is(err, errs.ErrInvalidPaymentMethod):
		return http.StatusBadRequest, err.Error()

	// 500 Internal Server Error — детали только в лог, клиенту общий текст
	default:
		return http.StatusInternalServerError, "внутренняя ошибка"
	}
}
