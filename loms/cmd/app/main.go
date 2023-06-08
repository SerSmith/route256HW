package main

import (
	"fmt"
	"log"
	"net"
	"route256/libs/mw/mylogging"
	"route256/libs/mw/mypanic"
	"route256/loms/internal/api/v1"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

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
	desc.RegisterLomsServer(s, service.NewServer())

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
