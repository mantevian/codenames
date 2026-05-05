import { signal } from "@preact/signals";
import CreateGameForm from "../../components/CreateGameForm";
import GameList from "../../components/GameList";
import { useState } from "preact/hooks";

export function Lobby() {
	const [joinCode, setJoinCode] = useState("");

	return <>
		<h1>lobby</h1>
		<CreateGameForm />
		<GameList />

		<section class="join-by-code">
			<p>connect by join code:</p>
			<input type="text" name="join-code" maxLength={4} value={joinCode} onInput={(e) => setJoinCode(e.currentTarget.value)} />
			<br />
			{joinCode.length == 4 ? <a href={`/room/${joinCode}`}>join '{joinCode}'</a> : <></>}
		</section>
	</>
}
