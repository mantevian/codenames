import { useLocation, useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { Me, Storage } from "../../storage/user";
import { Message, WSContext } from "../../components/WebSocketProvider";
import Game from "../../components/game/Game";

type JoinStatus = "waiting" | "success" | "fail" | "already_started";

export default function Room() {
	const { params } = useRoute();
	const ws = useContext(WSContext);
	const [joinStatus, setJoinStatus] = useState<JoinStatus>("waiting");

	useEffect(() => {
		ws.request({
			action: "join_game",
			payload: {
				"join_code": params.code
			}
		}).then(res => {
			if (res.payload.message == "game already started") {
				setJoinStatus("already_started");
			} else {
				setJoinStatus(res.payload.success ? "success" : "fail");
			}
		});

		ws.on("update_player_list", (msg: Message) => {
			Storage.players.value = msg.payload.players;
		});

		ws.on("game_started", () => {
			setJoinStatus("already_started");
		});
	}, []);

	switch (joinStatus) {
		case "already_started":
			return <Game />;

		case "waiting":
			return <section class="game-lobby">
				<h1>joining...</h1>
			</section>;

		default:
			return <section class="lobby">
				<h1>can't join this game</h1>
			</section>;
	}
}
