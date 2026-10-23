package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	inventoryAPI "github.com/MCMCXCII/spaceship-factory/inventory/internal/api/inventory/v1"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/inventory"
	inventoryService "github.com/MCMCXCII/spaceship-factory/inventory/internal/service/inventory"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const grpcAddress = ":50051"

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	//nolint:noctx
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	)

	repo := inventoryRepository.New()
	service := inventoryService.New(repo)
	api := inventoryAPI.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	// Контекст, который отменяется по SIGINT/SIGTERM или при падении сервера
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("grpc server запущен", "address", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	// Ждём либо сигнал от ОС, либо падение сервера
	<-ctx.Done()
	slog.Info("остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("сервер остановлен")
}
