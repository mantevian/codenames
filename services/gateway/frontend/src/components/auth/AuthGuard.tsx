import { useContext, useEffect, useState } from "preact/hooks";
import { WSContext, wsReady } from "../WebSocketProvider";

type AuthStatus = "waiting" | "success" | "fail";

export default function AuthGuard({ children }: { children?: any }) {
	const [status, setStatus] = useState<AuthStatus>("waiting");
	const ws = useContext(WSContext);

	useEffect(() => {
		wsReady.then(() => {
			ws.request({
				action: "validate_token"
			}).then(msg => {
				const success = msg.payload.success === true;
				if (success) {
					setStatus("success");
				} else {
					setStatus("fail");
				}
			}).catch(() => {
				setStatus("fail");
			});
		}).catch(() => {
			setStatus("fail");
		});
	}, []);

	switch (status) {
		case "success":
			return children;

		default:
			return null;
	}
}
