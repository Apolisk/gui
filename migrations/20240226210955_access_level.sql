-- +goose Up


alter table  users add access_level INT;

-- +goose Down

alter table users drop column access_level;


