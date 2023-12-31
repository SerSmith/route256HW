package loms

import (
	"net/url"
	"route256/libs/logger"
	"time"
)

const (
	stockPath       = "stocks"
	createOrderPath = "createOrder"
	waittime        = 5 * time.Second
)

type Client struct {
	stockURL       string
	createOrderURL string
}

func New(clientUrl string) *Client {

	stockUrl, _ := url.JoinPath(clientUrl, stockPath)

	createOrderUrl, _ := url.JoinPath(clientUrl, createOrderPath)
	logger.Info("Write", clientUrl, "|", stockUrl, "|", createOrderUrl)
	return &Client{stockURL: stockUrl,
		createOrderURL: createOrderUrl}
}
