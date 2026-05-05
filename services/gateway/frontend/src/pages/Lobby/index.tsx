import { signal } from "@preact/signals";
import CreateGameForm from "../../components/CreateGameForm";
import GameList from "../../components/GameList";

export function Lobby() {
	const joinCode = signal("");

	return <>
		<h1>lobby</h1>
		<CreateGameForm />
		<GameList />

		<div>
			<p>connect by join code:</p>
			<input type="text" onInput={(e) => joinCode.value = e.currentTarget.value} />
			<a href={`/room/${joinCode}`}>join '{joinCode}'</a>
		</div>
	</>
}
