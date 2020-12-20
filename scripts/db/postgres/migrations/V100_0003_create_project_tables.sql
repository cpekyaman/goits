-- +migrate Up

-- project status
create sequence config.project_status_seq;
create table config.project_status (
    id              integer not null default nextval('config.project_status_seq'),
    name            varchar(25) not null,
    description     varchar(250) not null,
    version         integer not null default 1
);
create unique index project_status_name_unq on config.project_status(name);
create unique index project_status_id_pk on config.project_status(id);

ALTER TABLE config.project_status
    add constraint project_status_pk primary key using INDEX project_status_id_pk
;

-- project type
create sequence config.project_type_seq;
create table config.project_type (
    id              integer not null default nextval('config.project_type_seq'),
    name            varchar(25) not null,
    description     varchar(250) not null,
    version         integer not null default 1
);
create unique index project_type_name_unq on config.project_type(name);
create unique index project_type_id_pk on config.project_type(id);

ALTER TABLE config.project_type
    add constraint project_type_pk primary key using INDEX project_type_id_pk
;

-- project
create sequence data.project_seq;
create table data.project (
    id                  bigint not null default nextval('data.project_seq'),
    name                varchar(50) not null,
    description         varchar(250) not null,
    type                integer not null,
    status              integer not null,
    version             integer not null default 1,
    create_time         timestamp with time zone not null default now(),
    last_modified_time  timestamp with time zone not null default now()
);
create unique index project_name_unq on data.project(name);
create unique index project_id_pk on data.project(id);
ALTER TABLE data.project
    add constraint project_pk primary key using INDEX project_id_pk,
    add constraint project_status_fk foreign key (status) references config.project_status(id)
    add constraint project_type_fk foreign key (type) references config.project_type(id)
;

-- +migrate Down
drop table data.project;
drop sequence data.project_seq;

drop table config.project_type;
drop sequence config.project_type_seq;

drop table config.project_status;
drop sequence config.project_status_seq;