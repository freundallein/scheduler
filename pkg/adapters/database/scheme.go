package database

const initialScheme = `
create extension if not exists "uuid-ossp";

do $$
begin
	if not exists (select 1 from pg_type where typname = 'task_state') then
		create type task_state as enum (
			'pending',
			'processing',
			'succeeded',
			'failed'
		);
    end if;
end
$$;

create table if not exists task (
	id uuid not null,
	claim_id uuid,
	state task_state not null default 'pending',
    execute_at timestamp with time zone not null,
    deadline timestamp with time zone not null,
    payload JSONB not null,
    result JSONB not null default '{}',
    meta JSONB not null default '{}',

	create_at timestamp with time zone not null default current_timestamp,
	primary key(id)
) with (
	autovacuum_vacuum_threshold = 100,
	autovacuum_vacuum_scale_factor = 0.2,
	autovacuum_vacuum_cost_delay = 20,
	autovacuum_vacuum_cost_limit = 200
);

create index if not exists task_state on task (execute_at, id) where state <> 'succeeded';

-- TODO: add failure table
`
