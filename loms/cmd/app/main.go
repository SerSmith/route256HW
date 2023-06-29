package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"route256/libs/closer"
	"route256/libs/mw/mylogging"
	"route256/libs/mw/mypanic"
	"route256/libs/tx"
	"route256/loms/internal/api/v1"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/postgres"
	desc "route256/loms/pkg/loms_v1"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

func run(ctx context.Context) error {

	err := config.Init()

	if err != nil {
		log.Fatalln("error reading config: ", err)
	}

	var closer = new(closer.Closer)

	BDPath := config.AppConfig.DSN()

	pool, err := pgxpool.Connect(ctx, BDPath)
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	closer.Add(func(ctx context.Context) error {
		pool.Close()
		return nil
	})

	provider := tx.New(pool)
	repo := postgres.New(provider)

	m := domain.New(repo)

	serv := service.NewServer(m)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mylogging.Interceptor),
		grpc.ChainUnaryInterceptor(mypanic.Interceptor),
	)
	reflection.Register(s)
	desc.RegisterLomsServer(s, serv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil

}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
