import { useLocation, useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { Message, WSContext } from "../WebSocketProvider";
import { Me, Storage } from "../../storage/user";

type JoinStatus = "waiting" | "success" | "fail";

export default function GameLobby() {
	const { params } = useRoute();
	const ws = useContext(WSContext);
	const { route } = useLocation();
	const [startResponse, setStartResponse] = useState("");

	useEffect(() => {

	}, []);

	function setReady(value: boolean) {
		ws.request({
			action: "set_ready",
			payload: {
				"player_id": Me.value?.id,
				"is_ready": value,
			}
		});
	}

	function startGame() {
		ws.request({
			action: "start_game",
			payload: {
				"player_id": Me.value?.id,
			}
		}).then(res => {
			setStartResponse(res.payload.message);
		});
	}

	function quit() {
		ws.request({
			action: "quit_game",
			payload: {
				"player_id": Me.value?.id
			}
		}).then(() => {
			route("/lobby");
		});
	}

	return <section class="game-lobby">
		<h1>game lobby</h1>
		<p>join code: {params.code}</p>

		<button onClick={quit}>quit</button>

		<label>
			{"I am "}
			<button onClick={() => setReady(!Me.value?.is_ready)}>
				{Me.value?.is_ready ? "ready" : "not ready"}
			</button>
		</label>

		<p>players:</p>
		<ul>
			{Storage.players.value.map(p => (
				<li key={p.id}>
					<p>{p.name} ({p.team}) {p.is_ready ? "ready" : "not ready"}</p>
				</li>
			))}
		</ul>

		<button onClick={startGame}>start</button>

		<p>{startResponse}</p>
	</section>;
}
