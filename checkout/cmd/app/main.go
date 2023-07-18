package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/postgres"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/closer"
	"route256/libs/logger"
	"route256/libs/metrics/call_counter"
	"route256/libs/metrics/time_hist"
	"route256/libs/tracer"
	"route256/libs/tx"
	"syscall"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)



const (
	grpcPort        = 50051
	shutdownTimeout = 5 * time.Second
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func run(ctx context.Context) error {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Fatal(http.ListenAndServe(os.Getenv("PROMETHEUSADDR"), nil))
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.MiddlewareGRPC,
			tracer.MiddlewareGRPC,
			callCounter.MiddlewareGRPC,
			timeHist.MiddlewareGRPC))
	reflection.Register(s)

	err = config.Init()

	if err != nil {
		logger.Fatalf("error reading config: %v", err)
	}

	loms, err := loms.New(config.AppConfig.Services.Loms)

	if err != nil {
		logger.Fatalf("failed to create loms client: %v", err)
	}

	var closer = new(closer.Closer)

	BDPath := config.AppConfig.DSN()

	pool, err := pgxpool.Connect(ctx, BDPath)
	if err != nil {
		logger.Fatalf("connect to db: %s", err)
	}

	closer.Add(func(ctx context.Context) error {
		pool.Close()
		return nil
	})

	provider := tx.New(pool)
	repo := postgres.New(provider)

	logger.Info("config", config.AppConfig)

	productservice, err := productservice.New(config.AppConfig.Services.ProductService, config.AppConfig.Token)

	if err != nil {
		logger.Fatalf("failed to create productservice client: %v", err)
	}

	desc.RegisterCheckoutServer(s, service.NewServer(loms, productservice, repo))

	logger.Info("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "error Close"))
	}

	return nil

}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.SetLoggerByEnvironment(*environment)

	if err := tracer.InitGlobal(domain.ServiceName); err != nil {
		logger.Fatalf("error init tracer: ", err)
	}

	if err := run(ctx); err != nil {
		logger.Fatal(err)
	}
}
