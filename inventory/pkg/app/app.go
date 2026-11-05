package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryAPI "github.com/MCMCXCII/spaceship-factory/inventory/internal/api/inventory/v1"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/inventory"
	inventoryService "github.com/MCMCXCII/spaceship-factory/inventory/internal/service/part"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// RegisterServices собирает слои InventoryService и регистрирует его в gRPC-сервере
func RegisterServices(s grpc.ServiceRegistrar, pool *pgxpool.Pool) {
	repo := inventoryRepository.New(pool)
	service := inventoryService.New(repo)
	api := inventoryAPI.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(s, api)
}
