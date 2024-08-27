-- +goose Up
-- +goose StatementBegin
create table sources
(
    id         serial primary key,
    name       varchar(255) not null,
    feed_url   varchar(255) not null,
    priority   int          not null,
    created_at timestamp    not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sources;
-- +goose StatementEnd
