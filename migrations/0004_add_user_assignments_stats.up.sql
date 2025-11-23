create table users_stats (
	user_id varchar(36) primary key not null
	, assignments_count bigint not null default 0
	, status_changes_count bigint not null default 0
	, updated_at timestamp not null default now()
)
