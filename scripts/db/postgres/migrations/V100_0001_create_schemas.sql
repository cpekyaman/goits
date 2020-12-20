-- +migrate Up
create schema event;
create schema config;
create schema data;

-- +migrate Down
drop schema event;
drop schema config;
drop schema data;
