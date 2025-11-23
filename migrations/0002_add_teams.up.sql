create table teams (
	id varchar(36) primary key not null
	, name varchar(50) unique not null
);

create index idx_teams_name on teams (name);

