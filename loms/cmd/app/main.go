package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/libs/closer"
	"route256/libs/logger"
	"route256/libs/metrics/call_counter"
	"route256/libs/metrics/time_hist"
	"route256/libs/tracer"
	"route256/libs/tx"
	"route256/loms/internal/api/v1"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/infrastructure/kafka"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/sender"
	desc "route256/loms/pkg/loms_v1"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

const grpcPort = 50052

func run(ctx context.Context) error {

	err := config.Init()

	if err != nil {
		logger.Fatalf("error reading config: %s", err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Fatal(http.ListenAndServe(os.Getenv("PROMETHEUSADDR"), nil))
	}()

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

	kafkaProducer, err := kafka.NewProducer(config.AppConfig.Kafka.Brokers)
	if err != nil {
		logger.Fatalf("failed to create kafkaProducer: %v", err)
	}

	ks := sender.NewKafkaSender(kafkaProducer, config.AppConfig.Kafka.TopicStatus)

	m := domain.New(repo, ks)

	serv := service.NewServer(m)

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
	desc.RegisterLomsServer(s, serv)

	logger.Info("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
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
