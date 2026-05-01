import { render } from "preact";
import { LocationProvider, Router, Route } from "preact-iso";

import { Header } from "./components/Header.js";
import { Home } from "./pages/Home/index.js";
import { NotFound } from "./pages/_404.js";
import "./style.css";
import AuthGuard from "./components/auth/AuthGuard.js";
import { WebSocketProvider } from "./components/WebSocketProvider.js";

export function App() {
	return (
		<LocationProvider>
			<WebSocketProvider>
				<Header />
				<main>
					<Router>
						<Route path="/" component={Home} />
						
						<AuthGuard path="/secret"><div>secret</div></AuthGuard>
						
						<Route default component={NotFound} />
					</Router>
				</main>
			</WebSocketProvider>
		</LocationProvider>
	);
}

render(<App />, document.getElementById("app")!);
