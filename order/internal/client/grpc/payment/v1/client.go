package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// callTimeout — таймаут на один вызов PaymentService
const callTimeout = 5 * time.Second

// Client — обёртка над gRPC-клиентом PaymentService
type Client struct {
	paymentClient paymentv1.PaymentServiceClient
}

// New создаёт обёртку над gRPC-клиентом PaymentService
func New(c paymentv1.PaymentServiceClient) *Client {
	return &Client{paymentClient: c}
}

// PayOrder проводит оплату заказа и возвращает UUID транзакции.
func (c *Client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: converter.PaymentMethodToDTO(method),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return uuid.Nil, errs.ErrInvalidPaymentMethod
		}
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	txUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("разобрать transaction_uuid %q: %w", resp.GetTransactionUuid(), err)
	}

	return txUUID, nil
}
