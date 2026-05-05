import { render } from "preact";
import { LocationProvider } from "preact-iso";

import "./style.css";

import { WebSocketProvider } from "./components/WebSocketProvider.js";
import App from "./App";

export function Root() {
	return (
		<LocationProvider>
			<WebSocketProvider>
				<App />
			</WebSocketProvider>
		</LocationProvider>
	);
}

render(<Root />, document.getElementById("app")!);
