package productservice

import (
	"log"
	"route256/checkout/pkg/product_service/product_service"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	wait_time  = time.Second * 10
	Limit      = 10
	BurstLimit = 10
)

type Client struct {
	psClient  product_service.ProductServiceClient
	token     string
	wait_time time.Duration
	Limiter   *rate.Limiter
}

func New(clientUrl string, token string) (*Client, error) {
	conn, err := grpc.Dial(clientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	MyLimiter := rate.NewLimiter(rate.Limit(Limit), BurstLimit)

	c := product_service.NewProductServiceClient(conn)
	return &Client{psClient: c, token: token, wait_time: wait_time, Limiter: MyLimiter}, nil
}
