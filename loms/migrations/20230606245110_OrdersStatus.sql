-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS OrdersStatus (
    orderID bigint Primary Key,
    status character varying 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS OrdersStatus;
-- +goose StatementEnd
