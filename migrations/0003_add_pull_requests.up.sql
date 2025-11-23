create table pull_requests (
	id varchar(36) primary key not null
	, name varchar(50) not null
	, author_id varchar(36) not null
	, reviewers_ids varchar(36)[] not null
	, status smallint not null
	, created_at timestamp default now()
	, merged_at timestamp
);

create index idx_pull_requests_reviewers_ids_gin on pull_requests using gin (reviewers_ids);

