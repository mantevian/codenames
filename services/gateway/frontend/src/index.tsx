import { render } from "preact";
import { LocationProvider, Router, Route } from "preact-iso";

import { Header } from "./components/Header.js";
import { Home } from "./pages/Home/index.js";
import { NotFound } from "./pages/_404.js";
import "./style.css";
import RouteAuthGuard from "./components/auth/RouteAuthGuard.js";
import { WebSocketProvider } from "./components/WebSocketProvider.js";
import LoginForm from "./components/auth/LoginForm.js";
import RegisterForm from "./components/auth/RegisterForm.js";
import { Lobby } from "./pages/Lobby/index.js";
import GameLobby from "./pages/GameLobby/index.js";

export function App() {
	return (
		<LocationProvider>
			<WebSocketProvider>
				<Header />
				<main>
					<Router>
						<Route path="/" component={Home} />
						<Route path="/login" component={LoginForm} />
						<Route path="/register" component={RegisterForm} />
						
						<RouteAuthGuard path="/lobby" component={Lobby} />
						<RouteAuthGuard path="/game/:code" component={GameLobby} />
						
						<Route default component={NotFound} />
					</Router>
				</main>
			</WebSocketProvider>
		</LocationProvider>
	);
}

render(<App />, document.getElementById("app")!);
