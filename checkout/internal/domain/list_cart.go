package domain

import (
	"context"
	"fmt"
	"log"
	"route256/libs/workerpool"
	"sync"
)


const (
	WORKERS_NUM = 5
)


func (m *Model) ListCart(ctx context.Context, user int64) (uint32, []ItemCart, error) {
	OrderItems, err := m.DB.GetCartDB(ctx, user)

	if err != nil {
		return 0, nil, fmt.Errorf("err in GetCartDB: %v", err)
	}

	resChan := make(chan ItemCart, len(OrderItems))

	wp := workerpool.New(WORKERS_NUM)
	wg := sync.WaitGroup{}


	for _, item := range OrderItems {
		wg.Add(1)
		NowItem := item
		err := wp.Run(ctx,
							func (ctx context.Context){
								defer wg.Done()

								product, err := m.productServiceClient.GetProduct(ctx, NowItem.SKU)

								if err != nil {
									log.Print("err in runOmeGetProductInstance ", err)
									product = &Product{
										Name:	"Unknown",
										Price:	0,
									}
								}

								resChan <- ItemCart{
									SKU:     NowItem.SKU,
									Product: *product,
									Count: NowItem.Count,
								}})
		if err != nil{
			return 0, nil, fmt.Errorf("Error in workerpool ", err)
		}
		
	}

	wg.Wait()

	var outCart []ItemCart
	var totalPrice uint32

	for range OrderItems {

		oneOutCart := <- resChan

		outCart = append(outCart, oneOutCart)
		totalPrice += oneOutCart.Product.Price * uint32(oneOutCart.Count)
	}


	return totalPrice, outCart, nil
}
