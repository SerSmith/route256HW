package productservice

import (
	"log"
	"route256/checkout/pkg/product_service/product_service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const wait_time = time.Second * 10

type Client struct {
	psClient  product_service.ProductServiceClient
	token     string
	wait_time time.Duration
}

func New(clientUrl string, token string) (*Client, error) {
	conn, err := grpc.Dial(clientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	c := product_service.NewProductServiceClient(conn)
	return &Client{psClient: c, token: token, wait_time: wait_time}, nil
}
