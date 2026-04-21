import { signal, useSignal } from "@preact/signals";
import { useEffect } from "preact/hooks";
import Game from "../types/game";

export default function WaitinGameList() {
	const games = useSignal<Game[]>([]);

	useEffect(() => {
		fetch("/api/v1/get_waiting_game_list", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${localStorage.getItem("token")}`,
			}
		})
			.then(res => {
				res.json().then(json => {
					if (json.success) {
						games.value = json.games;
					}
				});
			});
	}, []);

	return <>
		<section id="waiting-game-list">
			<h2>Games</h2>

			<ul>
				{games.value ? games.value.map(game => (
					<li>
						{game.id} {game.join_code} {game.language} {game.starting_team}
					</li>
				)) : <li>
					<p>No games</p>
				</li>}
			</ul>
		</section>
	</>;
}