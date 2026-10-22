package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	paymentService "github.com/MCMCXCII/spaceship-factory/payment/pkg/service"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

const grpcAddress = ":50052"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	var ls net.ListenConfig
	lis, err := ls.Listen(ctx, "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	const (
		grpcMaxConnectionIdle     = 15 * time.Minute
		grpcMaxConnectionAge      = 30 * time.Minute
		grpcMaxConnectionAgeGrace = 5 * time.Second
		grpcKeepaliveTime         = 5 * time.Minute
		grpcKeepaliveTimeout      = 1 * time.Second
		grpcMinPingInterval       = 5 * time.Minute
	)

	grpcServer := grpc.NewServer(
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
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentService.NewServer())

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск PaymentService", "адрес", grpcAddress)

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		slog.Info("получен сигнал, останавливаем PaymentService")
		stop()
		grpcServer.GracefulStop()
		slog.Info("PaymentService останавлен")
	case err := <-errCh:
		stop()
		if err != nil {
			slog.Error("ошибка работы сервера", "error", err)
			os.Exit(1)
		}
	}
}
