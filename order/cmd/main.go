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

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/MCMCXCII/spaceship-factory/order/pkg/app"
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

	pingTimeout = 5 * time.Second

	//nolint:gosec // учебный пример - скоро уберу
	OrderDSN = "postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable"
)

func main() {
	if err := run(); err != nil {
		slog.Error("OrderService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

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

	pool, err := pgxpool.New(ctx, OrderDSN)
	if err != nil {
		return fmt.Errorf("создать пул соединений: %w", err)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err = pool.Ping(pingCtx)
	if err != nil {
		return fmt.Errorf("проверить соединение с БД: %w", err)
	}

	slog.Info("Соединение к БД установлено")

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return fmt.Errorf("создать transaction manager: %w", err)
	}

	handler, err := app.NewHTTPHandler(
		pool,
		txManager,
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания HTTP-обработчика: %w", err)
	}

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка остановки HTTP-сервера: %w", err)
	}

	slog.Info("OrderService остановлен")
	return nil
}

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

func closeConn(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		slog.Error("ошибка закрытия gRPC-соединения", "service", name, "error", err)
	}
}
