-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications_sent (
    OldStatus varchar NOT NULL,
    NewStatus varchar NOT NULL,
    OrderID bigint NOT NULL,
    DT timestamp NOT NULL,
    UserID bigint NOT NULL,
    PRIMARY KEY (UserID, DT, OrderID, OldStatus)
);

INSERT INTO notifications_sent (OldStatus, NewStatus, OrderID, DT, UserID) VALUES 
    ('Null', 'Available', 1, '2023-07-10 10:10:10', 1),
    ('Available', 'Payed', 1, '2023-07-10 11:10:10', 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications_sent;
-- +goose StatementEnd
