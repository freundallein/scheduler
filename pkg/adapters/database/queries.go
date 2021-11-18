package database

const (
	create = `
	insert into 
		task(id, execute_at, deadline, payload, meta) 
	values 
		($1, $2, $3, $4, $5)
	returning id, claim_id, state, execute_at, deadline, payload, result, meta, task.created_at;
`
	findByID = `
	select
		id, claim_id, state, execute_at, deadline, payload, result, meta, task.created_at, task.done_at
	from 
		task 
	where id=$1;
`
	claimPending = `
	with claimed_tasks as (
		select 
			id 
		from task 
		where 
			state <> $2 
			and execute_at <= current_timestamp
			--and deadline >= current_timestamp
		order by execute_at
		limit $3
		for update skip locked
	)
	update task 
	set 
		state = $1, 
		execute_at = current_timestamp + interval '1 minute',
		claim_id = uuid_generate_v4()
	from claimed_tasks
	where task.id = claimed_tasks.id
	returning 
		task.id, 
		task.claim_id, 
		task.state, 
		task.execute_at, 
		task.deadline, 
		task.payload, 
		task.result, 
		task.meta,
		task.created_at;
`
	markAsSucceeded = `
	update task
	set 
		state = $1,
		claim_id = null,
		result = $4,
		done_at = current_timestamp
	where 
		id = $2
		and claim_id = $3;
`
	markAsFailed = `
	update task
	set 
		state = $1,
		claim_id = null,
		meta = meta::jsonb || $4 || CONCAT('{"attempts":', COALESCE(meta->>'attempts','0')::int + 1, '}')::jsonb
	where 
		id = $2
		and claim_id = $3;
`
	deleteStaleTasks = `
	delete from 
		task
	where
		state = 'succeeded' and done_at < current_timestamp - $1 * '1 hour'::interval;
`
)
