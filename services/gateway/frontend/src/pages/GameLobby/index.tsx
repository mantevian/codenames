import { useRoute } from "preact-iso";
import { useContext, useEffect, useState } from "preact/hooks";
import { WSContext } from "../../components/WebSocketProvider";

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
			</>;

		default:
			return <>
				<h1>can't join this game</h1>
			</>;
	}
}
