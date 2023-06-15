package main

import (
	"fmt"
	"log"
	"net"
	"os/signal"
	"context"
	"route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/mw/mylogging"
	"route256/libs/mw/mypanic"
	"route256/libs/tx"
	"syscall"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/closer"
	"github.com/jackc/pgx/v4/pgxpool"
	"route256/checkout/internal/repository/postgres"
)

const grpcPort = 50051



func run(ctx context.Context) error {


	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mylogging.Interceptor),
		grpc.ChainUnaryInterceptor(mypanic.Interceptor),
	)
	reflection.Register(s)

	err = config.Init()

	if err != nil {
		log.Fatalln("error reading config: ", err)
	}

	loms, err := loms.New(config.AppConfig.Services.Loms)

	if err != nil {
		log.Fatalf("failed to create loms client: %v", err)
	}


	var closer = new(closer.Closer)

	// "postgres://user:password@postgres_checkout:5433/checkout?sslmode=disable"

	BDPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
							config.AppConfig.DB.User,
							config.AppConfig.DB.Password,
							config.AppConfig.DB.Server,
							config.AppConfig.DB.Name,
						)

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

	log.Println("config", config.AppConfig)

	productservice, err := productservice.New(config.AppConfig.Services.ProductService, config.AppConfig.Token)

	if err != nil {
		log.Fatalf("failed to create productservice client: %v", err)
	}


	desc.RegisterCheckoutServer(s, service.NewServer(loms, productservice, repo))

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
