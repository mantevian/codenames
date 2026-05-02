import { signal, useSignal } from "@preact/signals";
import { useContext, useEffect } from "preact/hooks";
import Game from "../types/game";
import { WSContext, wsReady } from "./WebSocketProvider";

export default function GameList() {
	const games = useSignal<Game[]>([]);
	const ws = useContext(WSContext);

	useEffect(() => {
		wsReady.then(() => {
			ws.request({
				action: "get_game_list"
			}).then((res) => {
				if (res.payload && res.payload.success) {
					games.value = res.payload.games;
				}
			});
		});
	}, []);

	return <>
		<section id="game-list">
			<h2>Games</h2>

			<ul>
				{games.value ? games.value.map(game => (
					<li>
						{game.join_code} ({game.language}) <a href={`/game/${game.join_code}`}>join</a>
					</li>
				)) : <li>
					<p>No games</p>
				</li>}
			</ul>
		</section>
	</>;
}
