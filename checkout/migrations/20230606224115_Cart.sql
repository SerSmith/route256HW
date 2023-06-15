-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart (
    user_id bigint NOT NULL,
    sku bigint NOT NULL,
    "count" int,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
