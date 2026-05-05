import { useContext, useEffect } from "preact/hooks";
import { IsMyTurn, Me, Storage } from "../../storage/user";
import { Message, WSContext } from "../WebSocketProvider";
import "./style.css";

export default function Game() {
	const ws = useContext(WSContext);

	useEffect(() => {
		
	}, []);

	return <section class="game">
		<ul class="board">
			{Storage.tiles.value.map(tile =>
				<li class="tile" key={tile.position} data-type={tile.type ? tile.type : undefined} data-revealed={tile.is_revealed ? "" : undefined}>
					{tile.word}
				</li>
			)}
		</ul>

		<div>
			<p>your team: {Me.value?.team}</p>
			<p>your role: {Me.value?.role}</p>

			{IsMyTurn.value ?
				<>
					<p>it's your turn</p>
				</>
				: ""}
		</div>
	</section>;
}
