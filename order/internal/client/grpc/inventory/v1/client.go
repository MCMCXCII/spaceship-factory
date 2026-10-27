package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
)

const callTimeout = 5 * time.Second

type Client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func New(c inventoryv1.InventoryServiceClient) *Client {
	return &Client{inventoryClient: c}
}

func (c *Client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: converter.UUIDsToDTO(uuids),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, errs.ErrPartNotFound
			case codes.InvalidArgument:
				return nil, errs.ErrInvalidUUID
			}
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	parts, err := converter.PartsToModel(resp.GetParts())
	if err != nil {
		return nil, fmt.Errorf("разобрать список деталей: %w", err)
	}

	return parts, nil
}
