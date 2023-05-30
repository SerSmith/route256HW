package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {

	err := config.Init()
	if err != nil {
		log.Fatalln("error reading config: ", err)
	}

	loms := loms.New(config.AppConfig.Services.Loms)
	log.Println("config", config.AppConfig)

	productservice := productservice.New(config.AppConfig.Token, config.AppConfig.Services.ProductService)

	model := domain.New(loms, productservice)

	handAddToCart := addtocart.Handler{Model: model}
	http.Handle("/addToCart", srvwrapper.New(handAddToCart.Handle))

	handDeleteFromCart := deletefromcart.Handler{Model: model}
	http.Handle("/deleteFromCart", srvwrapper.New(handDeleteFromCart.Handle))

	handListCart := listcart.Handler{Model: model}
	http.Handle("/listCart", srvwrapper.New(handListCart.Handle))

	handPurchase := purchase.Handler{Model: model}
	http.Handle("/purchase", srvwrapper.New(handPurchase.Handle))

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
