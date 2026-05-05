export type Team = "red" | "blue";
export type Role = "operative" | "spymaster";

export type Game = {
	id: string;
	starting_team: Team;
	current_turn_team: Team;
	current_turn_role: Role;
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
	is_current_turn: boolean;
};

export type User = {
	id: string;
	name: string;
};

export type Tile = {
	position: number;
	game_id: string;
	type: string;
	is_revealed: boolean;
	word: string;
};

export type Turn = {
	player_id: string;
	clue_word: string;
	clue_number: number;
	guesses_left: number;
	created_at: string;
}
