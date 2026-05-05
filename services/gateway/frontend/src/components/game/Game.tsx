import { useContext, useState } from "preact/hooks";
import { IsMyTurn, Me, SortedPlayers, Storage } from "../../storage/user";
import { WSContext } from "../WebSocketProvider";
import "./style.css";
import { TargetedEvent } from "preact";
import { Tile } from "../../types/game";

export default function Game() {
	const ws = useContext(WSContext);
	const [submitResponse, setSubmitResponse] = useState("");

	async function submitClue(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);

		const res = await ws.request({
			action: "submit_clue",
			payload: {
				'number': parseInt(formData.get("number")?.toString() || "1"),
				'word': formData.get("word")?.toString()
			}
		});

		const { success, message } = res.payload;

		if (message) {
			setSubmitResponse(message);
		}
	}

	async function onClickTile(tile: Tile) {
		const res = await ws.request({
			action: "submit_guess",
			payload: {
				position: tile.position
			}
		});

		const { success, message } = res.payload;

		if (message) {
			setSubmitResponse(message);
		}
	}

	async function endTurn() {
		const res = await ws.request({
			action: "end_turn",
		});
	}

	return <section class="game" data-finished={Storage.game.value?.status == "finished" ? "" : undefined}>
		{Storage.game.value?.status == "finished" ? <h1>game finished!<br />winner: {Storage.game.value?.team_won}</h1> : <></>}

		<ul class="board">
			{Storage.tiles.value.map(tile =>
				<li
					class="tile"
					key={tile.position}
					data-type={tile.type ? tile.type : undefined}
					data-revealed={tile.is_revealed ? "" : undefined}
					onClick={() => onClickTile(tile)}
					data-enabled={(IsMyTurn.value && Me.value?.role == "operative" && Storage.turn.value!.guesses_left > 0) ? "" : undefined}>
					{tile.word}
				</li>
			)}
		</ul>

		<div>
			<ul class="players">
				{SortedPlayers.value.filter(p => p != undefined).map(p =>
					<li class="player" data-team={p.team} data-current={p.is_current_turn ? "" : undefined}>
						<p>
							<strong>{p.name}</strong>
							{" "} / {p.role} {Me.value?.id == p.id ? <strong>(you)</strong> : <></>}
						</p>
					</li>
				)}
			</ul>

			{Storage.turn.value?.clue_word ? <>
				<p class="clue">clue: <span>{Storage.turn.value?.clue_word} {Storage.turn.value?.clue_number}</span></p>
			</> : <></>
			}

			{IsMyTurn.value ?
				<>
					{Me.value?.role == "operative" ?
						<>
							<p>Click tiles to guess. Guesses left: <strong>{Storage.turn.value?.guesses_left}</strong></p>
							<button onClick={endTurn}>End turn</button>
						</>
						:
						<>
							<form onSubmit={submitClue}>
								<p><strong>Submit your clue:</strong></p>
								<label>
									<p>Word:</p>
									<input type="text" name="word" />
								</label>

								<label>
									<p>Number:</p>
									<input type="number" name="number" min={1} max={4} />
								</label>

								<input type="submit" />
							</form>
						</>
					}
				</>
				: ""}

				<p>{submitResponse}</p>
		</div>
	</section>;
}
