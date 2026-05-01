create type team as enum ('red', 'blue');
create type role as enum ('operative', 'spymaster');
create type game_status as enum ('waiting', 'playing', 'finished');
create type tile as enum ('red', 'blue', 'neutral', 'assassin');
create type language as enum ('en', 'ru');

create table games (
	id uuid,
	starting_team team,
	join_code varchar(4),
	language language,
	team_won team NULL,
	status game_status,
	finished_at timestamptz NULL,
	created_at timestamptz,

	primary key (id)
);

create table users (
	id uuid,
	name varchar(64),
	password varchar(256),
	created_at timestamptz,

	primary key (id)
);

create table players (
	id uuid,
	game_id uuid,
	user_id uuid,
	team team,
	role role,
	is_ready boolean,

	primary key (id),
	foreign key (game_id) references games (id),
	foreign key (user_id) references users (id)
);

create table tiles (
	position int,
	game_id uuid,
	type tile,
	is_revealed boolean,
	word varchar(64),

	primary key (position, game_id),
	foreign key (game_id) references games (id)
);

create table turns (
	id uuid,
	player_id uuid,
	game_id uuid,
	clue_word varchar(64),
	clue_number int,
	created_at timestamptz,

	primary key (id),
	foreign key (player_id) references players (id),
	foreign key (game_id) references games (id)
);

create table guesses (
	turn_id uuid,
	position int,
	created_at timestamptz,

	primary key (turn_id, position),
	foreign key (turn_id) references turns (id)
);

create table words (
	language language,
	word varchar(64),

	primary key (language, word)
);
