package productservice

import (
	"net/url"
	"time"
)

const (
	getProductPath = "get_product"
	waittime       = 5 * time.Second
)

type Client struct {
	token       string
	productPath string
}

func New(token, clientUrl string) *Client {
	productUrl, _ := url.JoinPath(clientUrl, getProductPath)
	return &Client{
		token:       token,
		productPath: productUrl,
	}
}
