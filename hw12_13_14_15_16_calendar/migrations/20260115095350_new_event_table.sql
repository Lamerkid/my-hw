-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id          UUID PRIMARY KEY,
    title       VARCHAR (50),
    start_time  TIMESTAMP,
    end_time    TIMESTAMP,
    description VARCHAR (255),
    user_id     UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
