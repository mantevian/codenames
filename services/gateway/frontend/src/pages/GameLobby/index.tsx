import { useLocation, useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { Message, WSContext } from "../../components/WebSocketProvider";
import { Me, Storage } from "../../storage/user";

type JoinStatus = "waiting" | "success" | "fail";

export default function GameLobby() {
	const { params } = useRoute();
	const { route } = useLocation();
	const ws = useContext(WSContext);
	const [joined, setJoined] = useState<JoinStatus>("waiting");

	useEffect(() => {
		ws.request({
			action: "join_game",
			payload: {
				"join_code": params.code
			}
		}).then(res => {
			setJoined(res.payload.success ? "success" : "fail");
		});

		ws.on("update_player_list", (msg: Message) => {
			Storage.players.value = msg.payload.players;
		});
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

	switch (joined) {
		case "waiting":
			return <section class="game-lobby">
				<h1>joining...</h1>
			</section>;

		case "success":
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
			</section>;

		default:
			return <section class="lobby">
				<h1>can't join this game</h1>
			</section>;
	}
}
