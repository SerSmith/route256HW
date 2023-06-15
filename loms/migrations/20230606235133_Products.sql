-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    sku bigint PRIMARY KEY,
    price int
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
