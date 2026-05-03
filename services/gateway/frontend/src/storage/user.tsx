import { computed, signal } from "@preact/signals";
import { Game, Player, User } from "../types/game";

export const Storage = {
	token: sessionSignal<string | undefined>("token", undefined),
	user: sessionSignal<User | undefined>("user", undefined),
	game: sessionSignal<Game | undefined>("game", undefined),
	players: sessionSignal<Player[]>("players", []),
} as const;

export const Me = computed(() => {
	const userId = Storage.user.value?.id;
	return Storage.players.value.find(p => p.user_id == userId);
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
