import { useContext, useEffect } from "preact/hooks";
import { Storage } from "../../storage/user";
import { Message, WSContext } from "../WebSocketProvider";
import "./style.css";

export default function Game() {
	const ws = useContext(WSContext);

	useEffect(() => {
		
	}, []);

	return <section class="game">
		<h1>game</h1>

		<ul class="board">
			{Storage.tiles.value.map(tile =>
				<li class="tile" key={tile.position} data-type={tile.type ? tile.type : undefined} data-revealed={tile.is_revealed ? "" : undefined}>
					{tile.word}
				</li>
			)}
		</ul>
	</section>;
}
