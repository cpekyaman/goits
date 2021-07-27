-- +migrate Up

-- audit event type
create sequence config.audit_event_type_seq;
create table config.audit_event_type (
    id integer not null default nextval('config.audit_event_type_seq'),
    name varchar(25) not null,
)
create unique index audit_event_type_name_unq on config.audit_event_type(name);
create unique index audit_event_type_id_pk on config.audit_event_type(id);
ALTER TABLE config.audit_event_type
    add constraint audit_event_type_pk primary key using INDEX audit_event_type_id_pk
;

-- audit event
create sequence event.audit_event_seq;
create table event.audit_event (
    id bigint not null default nextval('event.audit_event_seq'),
    reference uuid not null,
    event_type integer not null,
    target_type varchar(50) not null,
    target_id bigint not null,
    event_data jsonb not null,
    event_time timestamp with time zone not null default now(),
    description varchar(200) null
);
create unique index audit_event_id_pk on event.audit_event(id);
create index audit_event_reference_idx on event.audit_event(reference);
create index audit_event_target_idx on event.audit_event(target_type, target_id);
ALTER TABLE event.audit_event
    add constraint audit_event_pk primary key using INDEX audit_event_id_pk,
    add constraint audit_event_type_fk foreign key (event_type) references config.audit_event_type(id)
;

-- domain event type
create sequence config.domain_event_type_seq;
create table config.domain_event_type (
    id integer not null default nextval('config.domain_event_type_seq'),
    name varchar(25) not null,
)
create unique index domain_event_type_name_unq on config.domain_event_type(name);
create unique index domain_event_type_id_pk on config.domain_event_type(id);
ALTER TABLE config.domain_event_type
    add constraint domain_event_type_pk primary key using INDEX domain_event_type_id_pk
;

-- domain event
create sequence event.domain_event_seq;
create table event.domain_event (
    id bigint not null default nextval('event.domain_event_seq'),
    reference uuid not null,
    event_type integer not null,
    target_type varchar(50) not null,
    target_id bigint not null,
    event_data jsonb not null,
    event_time timestamp with time zone not null default now()
);
create unique index domain_event_id_pk on event.domain_event(id);
create index domain_event_target_idx on event.domain_event(target_type, target_id);
create index domain_event_reference_idx on event.domain_event(reference);
ALTER TABLE event.domain_event
    add constraint domain_event_pk primary key using INDEX domain_event_id_pk,
    add constraint domain_event_type_fk foreign key (event_type) references event.domain_event_type(id)
;

-- +migrate Down
drop table event.audit_event;
drop sequence event.audit_event_seq;
drop table config.audit_event_type;
drop sequence config.audit_event_type_seq;

drop table event.domain_event;
drop sequence event.domain_event_seq;
drop table config.domain_event_type;
drop sequence config.domain_event_type_seq;