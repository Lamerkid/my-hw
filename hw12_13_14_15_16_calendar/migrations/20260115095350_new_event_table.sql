-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id            UUID PRIMARY KEY,
    title         VARCHAR (50),
    description   VARCHAR (255),
    start_time    TIMESTAMP,
    end_time      TIMESTAMP,
    user_id       UUID,
    notify_before INTERVAL,
    notified      BOOL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
