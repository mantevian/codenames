export type Game = {
	id: string;
	starting_team: string;
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
	team: "red" | "blue";
	role: "operative" | "spymaster";
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
