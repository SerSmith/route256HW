package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/domain"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {

	model := domain.New()

	handCreateOrder := &createorder.Handler{Model: model}
	http.Handle("/createOrder", srvwrapper.New(handCreateOrder.Handle))

	handListOrder := &listorder.Handler{Model: model}
	http.Handle("/listOrder", srvwrapper.New(handListOrder.Handle))

	handOrderPayed := &orderpayed.Handler{Model: model}
	http.Handle("/orderPayed", srvwrapper.New(handOrderPayed.Handle))

	handCancelOrder := &cancelorder.Handler{Model: model}
	http.Handle("/cancelOrder", srvwrapper.New(handCancelOrder.Handle))

	handStocks := &stocks.Handler{Model: model}
	http.Handle("/stocks", srvwrapper.New(handStocks.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
