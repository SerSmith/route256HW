-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS StocksAvalible (
    sku bigint,
    warehouseID bigint,
    "count" bigint,
    PRIMARY KEY (sku, warehouseID)
);

INSERT INTO StocksAvalible (sku, warehouseID, "count") VALUES 
    (1076963, 1, 10),
    (1148162, 1, 20),
    (1625903, 1, 15),
    (1625903, 2, 15);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS StocksAvalible;
-- +goose StatementEnd
