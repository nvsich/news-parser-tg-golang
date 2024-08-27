-- +goose Up
-- +goose StatementBegin
create table articles
(
    id           bigserial primary key,
    source_id    bigint       not null,
    title        varchar(255) not null,
    link         text         not null unique,
    published_at timestamp    not null,
    created_at   timestamp    not null default now(),
    posted_at    timestamp,
    constraint fk_articles_source_id
        foreign key (source_id)
            references sources (id)
            on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists articles;
-- +goose StatementEnd
