-- +goose Up
-- +goose StatementBegin
alter table sources
    alter column created_at set default current_timestamp,
    alter column updated_at set default current_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table sources
    alter column created_at drop default,
    alter column updated_at drop default;
-- +goose StatementEnd
