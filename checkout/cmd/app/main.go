package main

import (
	"fmt"
	"log"
	"net"
	"route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/mw/mylogging"
	"route256/libs/mw/mypanic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {

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

	log.Println("config", config.AppConfig)

	productservice, err := productservice.New(config.AppConfig.Services.ProductService, config.AppConfig.Token)

	if err != nil {
		log.Fatalf("failed to create productservice client: %v", err)
	}

	desc.RegisterCheckoutServer(s, service.NewServer(loms, productservice))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
