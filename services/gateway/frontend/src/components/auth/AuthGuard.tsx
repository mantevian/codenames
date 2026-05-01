import { useContext, useEffect, useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import { WSContext, wsReady } from "../WebSocketProvider";

type AuthStatus = "waiting" | "success" | "fail";

export default function AuthGuard({ path, children }: { path?: string, children?: any }) {
	const { path: locationPath, route } = useLocation();
	const [status, setStatus] = useState<AuthStatus>("waiting");
	const ws = useContext(WSContext);

	useEffect(() => {
		wsReady.then(() => {
			ws.request({
				action: "validate_token"
			}).then(msg => {
				const success = msg.payload.success === true;
				if (success) {
					console.log(msg);
					setStatus("success");
				} else {
					setStatus("fail");
					route("/");
				}
			}).catch(() => {
				setStatus("fail");
				route("/");
			});
		}).catch(() => {
			setStatus("fail");
			route("/");
		});
	}, [locationPath]);

	switch (status) {
		case "waiting":
			return <div>Checking authentication...</div>;

		case "success":
			return children;

		default:
			return null;
	}
}
