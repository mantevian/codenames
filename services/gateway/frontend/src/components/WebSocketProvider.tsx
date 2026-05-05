import { createContext } from "preact";
import { useEffect } from "preact/hooks";
import { signal, Signal } from "@preact/signals";
import { v4 } from "uuid";
import { Me, Storage } from "../storage/user";

type WSStatus = "connecting" | "open" | "closed" | "error";

export type Message = {
	id?: string;
	action?: string;
	payload?: any;
};

export type Listener = (msg: Message) => any;
export type ListenerMap = {
	[key: string]: Listener[]
};

type PendingMap = {
	[id: string]: {
		resolve: (v: Message) => void;
		reject: (e: any) => void;
		timer?: number
	}
};

export type WSContextValue = {
	status: Signal<WSStatus>;
	send: (msg: Message) => void;
	request: (out: Message) => Promise<Message>;
	on: (action: string, fn: Listener) => void;
	off: (action: string, fn: Listener) => void;
};

const WS_PATH = "/ws";

export const WSContext = createContext<WSContextValue>({
	status: signal("closed"),
	send: () => { },
	request: (_out: Message) => Promise.reject(new Error("no websocket")),
	on: () => { },
	off: () => { },
});

let ws: WebSocket | null = null;
const listeners: ListenerMap = {};
const pendings: PendingMap = {};
const status = signal<WSStatus>("closed");

let openResolve: () => void;
let openReject: (e: any) => void;
export const wsReady = new Promise<void>((res, rej) => {
	openResolve = res;
	openReject = rej;
});

type Props = { children: preact.ComponentChildren };

export function WebSocketProvider({ children }: Props) {
	useEffect(() => {
		status.value = "connecting";
		ws = new WebSocket(WS_PATH);

		ws.addEventListener("open", () => {
			console.log("websocket opened");
			status.value = "open";
			openResolve();
		});

		ws.addEventListener("message", (ev) => {
			const data: Message = JSON.parse(ev.data);
			console.log(`incoming:`, data);

			if (!data || typeof data.action !== "string") return;

			const id = data.id;

			if (id) {
				const p = pendings[id];
				if (p) {
					if (p.timer) clearTimeout(p.timer);
					delete pendings[id];
					try {
						p.resolve(data as Message);
					} catch (err) { }
				}
			}

			const arr = listeners[data.action];
			if (!arr || arr.length === 0) return;

			for (const fn of [...arr]) {
				try {
					fn(data as Message);
				} catch (err) {
					console.error("websocket listener error", err);
				}
			}
		});

		ws.addEventListener("close", () => {
			status.value = "closed";
			ws = null;
		});

		ws.addEventListener("error", (e) => {
			status.value = "error";
			try { ws?.close(); } catch { }
			ws = null;
			openReject(e);
		});

		return () => {
			try { ws?.close(); } catch { }
			ws = null;
			for (const k of Object.keys(listeners)) delete listeners[k];
			status.value = "closed";
		};
	}, []);

	const send = (msg: Message) => {
		if (ws && ws.readyState === WebSocket.OPEN) {
			try {
				if (!msg.payload) {
					msg.payload = {};
				}
				
				msg.payload["player_id"] = Me.value?.id;

				const finalMessage = {
					...msg,
					token: Storage.token,
				};

				console.log("sent:", finalMessage);
				ws.send(JSON.stringify(finalMessage));
			} catch (e) {
				console.warn("failed to send websocket message", e);
			}
		} else {
			console.warn("websocket is not open, can't send");
		}
	};

	const request = (out: Message): Promise<Message> => {
		return new Promise((resolve, reject) => {
			const id = v4();

			if (pendings[id]) {
				reject(new Error('duplicate request id'));
				return;
			}

			const entry = {
				resolve,
				reject,
				timer: -1
			};

			pendings[id] = entry;

			entry.timer = window.setTimeout(() => {
				delete pendings[id];
				reject(new Error('websocket request timed out'));
			}, 10000);

			out.id = id;
			
			send(out);
		});
	};


	const on = (action: string, fn: Listener) => {
		const arr = listeners[action] || [];
		arr.push(fn);
		listeners[action] = arr;
	};

	const off = (action: string, fn: Listener) => {
		const arr = listeners[action];
		if (!arr) return;
		const i = arr.indexOf(fn);
		if (i !== -1) arr.splice(i, 1);
		if (arr.length === 0) delete listeners[action];
	};

	const value: WSContextValue = {
		status,
		send,
		request,
		on,
		off
	};

	return <WSContext.Provider value={value}>{children}</WSContext.Provider>;
}
