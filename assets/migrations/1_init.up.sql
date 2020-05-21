create table person(
    id uuid not null ,
    login varchar(128) not null,
    password varchar(16) not null,
    full_name text not null,
    email varchar(1024) not null,
    role varchar(256) not null default 'employee',
    manager_id uuid,
    primary key (id)
);

alter table person add constraint person_manager foreign key (manager_id) references person(id) on update cascade on delete set null;

create table support(
    id uuid not null,
    person_id uuid not null,
    is_manager bool not null default false,
    primary key (id)
);

alter table support add constraint support_person foreign key (person_id) references person(id) on update cascade on delete cascade ;

CREATE SEQUENCE ordinal_seq;

create table incident (
    id uuid not null,
    ordinal int not null default nextval('ordinal_seq'),
    description text not null,
    created_at timestamp(0) not null,
    resolved_at timestamp(0),
    deadline timestamp(0) not null ,
    assignee_id uuid,
    creator_id uuid not null,
    status varchar(256) not null default 'created',
    comment text,
    type varchar(256) not null default 'maintenance',
    priority int not null default 0,
    primary key (id)
);

ALTER SEQUENCE ordinal_seq OWNED BY incident.ordinal;

alter table incident add constraint incident_assignee foreign key (assignee_id) references support(id) on update cascade on delete cascade;
alter table incident add constraint incident_creator foreign key (creator_id) references person(id) on update cascade on delete cascade;

create table equipment(
    id uuid not null,
    name varchar(2048) not null,
    description text not null,
    price int not null,
    primary key (id)
);

create table equipment_incident
(
    id           uuid not null,
    equipment_id uuid not null,
    incident_id  uuid not null,
    deadline     timestamp(0) not null,
    need_approval bool not null default false,
    approved bool not null default false,
    primary key (id)
);

alter table equipment_incident add constraint equipment_incident_to_incident foreign key (incident_id) references incident(id) on update cascade on delete cascade;
alter table equipment_incident add constraint equipment_incident_to_equipment foreign key (equipment_id) references equipment(id) on update cascade on delete cascade;

create table equipment_assignment(
    id uuid not null,
    person_id uuid not null,
    equipment_id uuid not null,
    deadline timestamp(0) not null,
    created_at timestamp(0) not null,
    primary key (id)
);

alter table equipment_assignment add constraint equipment_assignment_to_equipment foreign key (equipment_id) references equipment(id) on update cascade on delete cascade;
alter table equipment_assignment add constraint equipment_assignment_to_person  foreign key (person_id) references person(id) on update cascade on delete cascade;

CREATE TABLE "message"
(
    "id"        uuid          NOT NULL,
    "person_id" uuid          NOT NULL,
    "event_id"  uuid          not null,
    "login"     uuid          not null,
    "full_name" varchar(1024) not null,
    "time"      timestamp     not null,
    "message"   text          not null,
    PRIMARY KEY ("id")
);