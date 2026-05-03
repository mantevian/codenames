import { useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { Message, WSContext } from "../../components/WebSocketProvider";
import { Me, Storage } from "../../storage/user";

type JoinStatus = "waiting" | "success" | "fail";

export default function GameLobby() {
	const { params } = useRoute();
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

	switch (joined) {
		case "waiting":
			return <>
				<h1>joining...</h1>
			</>;

		case "success":
			return <>
				<h1>game lobby</h1>
				<p>join code: {params.code}</p>

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
			</>;

		default:
			return <>
				<h1>can't join this game</h1>
			</>;
	}
}
