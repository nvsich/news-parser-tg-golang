-- +goose Up
-- +goose StatementBegin
create table articles
(
    id           serial primary key,
    source_id    int references sources (source_id),
    title        text      not null,
    link         text      not null,
    summary      text      not null,
    published_at timestamp not null,
    created_at   timestamp not null,
    posted_at    timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
