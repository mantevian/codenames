import { useContext, useEffect, useState } from "preact/hooks";
import { RouteProps, useLocation } from "preact-iso";
import { WSContext, wsReady } from "../WebSocketProvider";
import { Component, ComponentType } from "preact";

type AuthStatus = "waiting" | "success" | "fail";

interface RouteAuthGuardProps {
	component: ComponentType<any>;
	redirectTo?: string;
	[key: string]: any;
}

export default function RouteAuthGUard({
	component: Component,
	redirectTo = "/login",
	...rest
}: RouteAuthGuardProps) {
	const { route } = useLocation();
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
					route(redirectTo, true);
				}
			}).catch(() => {
				setStatus("fail");
				route(redirectTo, true);
			});
		}).catch(() => {
			setStatus("fail");
			route(redirectTo, true);
		});
	}, []);

	switch (status) {
		case "waiting":
			return <div>Checking authentication...</div>;

		case "success":
			return <Component {...rest} />;

		default:
			return null;
	}
}
