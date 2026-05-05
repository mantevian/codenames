import { Router, Route } from "preact-iso";

import "./style.css";

import { Header } from "./components/Header.js";
import { Home } from "./pages/Home/index.js";
import { NotFound } from "./pages/_404.js";
import RouteAuthGuard from "./components/auth/RouteAuthGuard.js";
import LoginForm from "./components/auth/LoginForm.js";
import RegisterForm from "./components/auth/RegisterForm.js";
import { Lobby } from "./pages/Lobby/index.js";
import Room from "./pages/Room/index.js";

import { useContext, useEffect } from "preact/hooks";
import { WSContext, wsReady } from "./components/WebSocketProvider";
import { Storage } from "./storage/user";

export default function App() {
	const ws = useContext(WSContext);

	useEffect(() => {
		function checkAuth() {
			wsReady.then(() => {
				ws.request({
					action: "validate_token"
				}).then(msg => {
					const success = msg.payload.success === true;
					if (!success) {
						Storage.logout();
					}
				}).catch(() => {
					Storage.logout();
				});
			}).catch(() => {
				Storage.logout();
			});
		}

		checkAuth();
		setInterval(checkAuth, 30 * 1000);
	}, []);

	return <>
		<Header />
		<main>
			<Router>
				<Route path="/" component={Home} />
				<Route path="/login" component={LoginForm} />
				<Route path="/register" component={RegisterForm} />

				<RouteAuthGuard path="/lobby" component={Lobby} />
				<RouteAuthGuard path="/room/:code" component={Room} />

				<Route default component={NotFound} />
			</Router>
		</main></>;
}
