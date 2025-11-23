create table pull_requests_stats (
	pull_request_id varchar(36) primary key not null
	, assignments_count bigint not null default 0
	, updated_at timestamp not null default now()
)
