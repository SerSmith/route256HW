-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS StocksBought (
    sku bigint,
    warehouseID bigint,
    "count" bigint,
    PRIMARY KEY (sku, warehouseID)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS StocksBought;
-- +goose StatementEnd
