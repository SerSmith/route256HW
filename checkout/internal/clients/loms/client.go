package loms

import (
	"fmt"
	"log"
	"route256/checkout/pkg/loms/loms"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const wait_time = time.Second * 10

type Client struct {
	loms      loms_v1.LomsClient
	wait_time time.Duration
}

func New(clientUrl string) (*Client, error) {
	fmt.Printf("TRY TO connect to %s", clientUrl)
	conn, err := grpc.Dial(clientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	c := loms_v1.NewLomsClient(conn)
	return &Client{loms: c, wait_time: wait_time}, nil
}
