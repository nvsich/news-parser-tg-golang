-- +goose Up
-- +goose StatementBegin
create table sources
(
    source_id   serial primary key,
    source_name text      not null,
    feed_url    text      not null,
    created_at  timestamp not null,
    updated_at  timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sources;
-- +goose StatementEnd
