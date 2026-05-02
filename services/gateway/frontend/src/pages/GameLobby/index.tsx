import { useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { Message, WSContext } from "../../components/WebSocketProvider";

type JoinStatus = "waiting" | "success" | "fail";

type Player = {
	id: string;
	name: string;
	team: "red" | "blue";
	role: "operative" | "spymaster";
	is_ready: boolean;
};

export default function GameLobby() {
	const { params } = useRoute();
	const ws = useContext(WSContext);
	const [joined, setJoined] = useState<JoinStatus>("waiting");
	const [players, setPlayers] = useState<Player[]>([]);

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
			setPlayers(msg.payload.players);
		});
	}, []);

	switch (joined) {
		case "waiting":
			return <>
				<h1>joining...</h1>
			</>;

		case "success":
			return <>
				<h1>game lobby</h1>
				<p>join code: {params.code}</p>

				<p>players:</p>
				<ul>
					{players.map(p => (
						<li key={p.id}>
							<p>{p.name} ({p.team}, {p.role}) {p.is_ready ? "ready" : "not ready"}</p>
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
