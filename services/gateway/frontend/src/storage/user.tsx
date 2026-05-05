import { computed, signal } from "@preact/signals";
import { Game, Player, Team, Tile, Turn, User } from "../types/game";

export const Storage = {
	token: sessionSignal<string | undefined>("token", undefined),
	user: sessionSignal<User | undefined>("user", undefined),
	game: sessionSignal<Game | undefined>("game", undefined),
	turn: sessionSignal<Turn | undefined>("turn", undefined),
	tiles: sessionSignal<Tile[]>("tiles", []),
	players: sessionSignal<Player[]>("players", []),
	logout: () => {
		Storage.token.value = undefined;
		Storage.user.value = undefined;
		Storage.game.value = undefined;
		Storage.tiles.value = [];
		Storage.players.value = [];
	},
	clear: () => {
		Storage.logout();
	}
} as const;

export const Me = computed(() => {
	const userId = Storage.user.value?.id;
	return Storage.players.value.find(p => p.user_id == userId);
});

export const IsMyTurn = computed(() => {
	if (!Storage.game.value?.current_turn_role || !Storage.game.value?.current_turn_team) {
		return false;
	}
	
	return Storage.game.value?.current_turn_role == Me.value?.role && Storage.game.value?.current_turn_team == Me.value?.team;
});

export const SortedPlayers = computed(() => {
	const game = Storage.game.value;

	if (!game) {
		return [];
	}

	const startingTeam = game.starting_team;
	const otherTeam: Team = startingTeam == "red" ? "blue" : "red";

	const currentRole = game.current_turn_role;
	const currentTeam = game.current_turn_team;

	const players = Storage.players.value;

	const orderedPlayers = [
		players.find(p => p.team == startingTeam && p.role == "spymaster")!,
		players.find(p => p.team == startingTeam && p.role == "operative")!,
		players.find(p => p.team == otherTeam && p.role == "spymaster")!,
		players.find(p => p.team == otherTeam && p.role == "operative")!
	];

	return orderedPlayers.map(p => ({
		...p,
		is_current_turn: p.role == currentRole && p.team == currentTeam
	}));
});

export function sessionSignal<T>(key: string, initial: T, serialize = JSON.stringify, deserialize = JSON.parse) {
	try {
		const raw = sessionStorage.getItem(key);
		if (raw != null) {
			initial = deserialize(raw);
		}
	} catch (e) {

	}

	const s = signal(initial);

	const unsubscribe = s.subscribe((value) => {
		try {
			if (value === undefined) {
				sessionStorage.removeItem(key);
			} else {
				sessionStorage.setItem(key, serialize(value));
			}
		} catch (e) {

		}
	});

	return s;
}
