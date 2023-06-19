-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Orders (
    orderID serial NOT NULL,
    user_id bigint NOT NULL,
    sku bigint NOT NULL,
    "count" int,
    PRIMARY KEY (orderID, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Orders;
-- +goose StatementEnd
