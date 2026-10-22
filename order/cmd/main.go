package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/MCMCXCII/spaceship-factory/order/pkg/handler"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

const (
	httpAddress             = ":8080"
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 30 * time.Second

	grpcKeepaliveTime    = 5 * time.Minute
	grpcKeepaliveTimeout = 1 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("OrderService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}

func run() error {
	inventoryConn, err := newGRPCConn(inventoryServiceAddress)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к InventoryService: %w", err)
	}
	defer closeConn(inventoryConn, "InventoryService")

	paymentConn, err := newGRPCConn(paymentServiceAddress)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PaymentService: %w", err)
	}
	defer closeConn(paymentConn, "PaymentService")

	h := handler.NewHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		handler.NewOrderStore(),
	)

	orderServer, err := handler.SetupServer(h)
	if err != nil {
		return fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("запуск OrderService", "address", httpAddress)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		slog.Info("получен сигнал, останавливаем OrderService")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("ошибка HTTP-сервера: %w", err)
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка остановки HTTP-сервера: %w", err)
	}

	slog.Info("OrderService остановлен")
	return nil
}

// newGRPCConn создаёт gRPC-соединение с keepalive-параметрами
func newGRPCConn(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
}

// closeConn закрывает gRPC-соединение и логирует ошибку закрытия
func closeConn(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		slog.Error("ошибка закрытия gRPC-соединения", "service", name, "error", err)
	}
}
