-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS StocksReserved (
    sku bigint,
    orderID bigint,
    warehouseID bigint,
    "count" bigint,
    PRIMARY KEY (orderID, warehouseID, sku)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS StocksReserved;
-- +goose StatementEnd
