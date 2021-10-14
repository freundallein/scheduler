-- Initial migration

drop type if exists state cascade;
create type task_state as enum (
	'pending',
	'processing',
	'succeeded',
	'failed'
);

drop table if exists task;
create table task (
	id uuid not null,
	claim_id uuid,
	state task_state not null default 'pending',
    executed_at timestamp with time zone not null,
    deadline timestamp with time zone not null,
    payload JSONB not null,
    result JSONB not null default '{}',
    meta JSONB not null default '{}',
	primary key(id)
);

create index task_state on task (executed_at, id) where state <> 'succeeded';