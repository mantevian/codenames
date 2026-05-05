export type Team = "red" | "blue";
export type Role = "operative" | "spymaster";

export type Game = {
	id: string;
	starting_team: Team;
	current_team: Team;
	current_role: Role;
	join_code: string;
	language: string;
	team_won: string;
	status: string;
	finished_at: string;
	created_at: string;
};

export type Player = {
	id: string;
	game_id: string;
	user_id: string;
	team: Team;
	role: Role;
	is_ready: boolean;
	name: string;
};

export type User = {
	id: string;
	name: string;
}

export type Tile = {
	position: number;
	game_id: string;
	type: string;
	is_revealed: boolean;
	word: string;
}
