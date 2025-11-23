create table users (
	id varchar(36) primary key not null
	, name varchar(50) not null
	, is_active bool not null
	, team_id varchar(36) not null
);

create index idx_users_team_id_active on users (team_id) where is_active = true;

