package stocks

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	Stocks []StockItem `json:"stocks"`
}

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Request struct {
	SKU uint32 `json:"sku"`
}

func (r Request) Validate() error {

	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	logger.Info("stocks, request: %+v", req)
	stocks, err := h.Model.Stocks(ctx, req.SKU)
	if err != nil {
		logger.Info("list stocks: %s", err)
		return Response{}, err
	}
	respStocks := make([]StockItem, 0, len(stocks))
	for _, stock := range stocks {
		respStocks = append(respStocks, StockItem{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}
	return Response{
		Stocks: respStocks,
	}, nil
}
