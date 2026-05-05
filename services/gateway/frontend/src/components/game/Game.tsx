import { useContext, useEffect } from "preact/hooks";
import { IsMyTurn, Me, Storage } from "../../storage/user";
import { Message, WSContext } from "../WebSocketProvider";
import "./style.css";
import { TargetedEvent } from "preact";
import { signal } from "@preact/signals";
import { Tile } from "../../types/game";

export default function Game() {
	const ws = useContext(WSContext);
	const submitResponse = signal("");

	useEffect(() => {

	}, []);

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
			submitResponse.value = message;
		}
	}

	async function onClickTile(tile: Tile) {
		const guessesLeft = Storage.turn.value?.guesses_left || 0;
		const canSubmitGuess = IsMyTurn.value && Me.value?.role == "operative" && guessesLeft > 0;

		if (canSubmitGuess) {
			submitResponse.value = "not your turn";
		}

		const res = await ws.request({
			action: "submit_guess",
			payload: {
				position: tile.position
			}
		});

		const { success, message } = res.payload;

		if (message) {
			submitResponse.value = message;
		}
	}

	async function endTurn() {
		const res = await ws.request({
			action: "end_turn",
		});
	}

	return <section class="game">
		<ul class="board">
			{Storage.tiles.value.map(tile =>
				<li
					class="tile"
					key={tile.position}
					data-type={tile.type ? tile.type : undefined}
					data-revealed={tile.is_revealed ? "" : undefined}
					onClick={() => onClickTile(tile)}
					data-enabled={(IsMyTurn.value && Me.value?.role == "operative") ? "" : undefined}>
					{tile.word}
				</li>
			)}
		</ul>

		<div>
			<p>you: {Me.value?.team} {Me.value?.role}</p>

			{Storage.turn.value?.clue_word ? <>
				<p>clue: {Storage.turn.value?.clue_word} {Storage.turn.value?.clue_number}</p>
			</> : <></>
			}

			{IsMyTurn.value ?
				<>
					<p>it's your turn</p>

					{Me.value?.role == "operative" ?
						<>
							<p>click tiles to guess</p>
							<p>guesses left: {Storage.turn.value?.guesses_left}</p>
							<button onClick={endTurn}>end turn</button>
						</>
						:
						<>
							<p>submit your clue:</p>
							<form onSubmit={submitClue}>
								<label>
									<p>word:</p>
									<input type="text" name="word" />
								</label>

								<label>
									<p>number:</p>
									<input type="number" name="number" min={1} max={4} />
								</label>

								<input type="submit" />
							</form>
						</>
					}
					
					<p class="error">{submitResponse}</p>

				</>
				: ""}
		</div>
	</section>;
}
