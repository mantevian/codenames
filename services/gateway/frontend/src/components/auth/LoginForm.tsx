import { TargetedEvent } from "preact";
import { useContext, useRef } from "preact/hooks";
import { WSContext } from "../WebSocketProvider";
import { Storage } from "../../storage/user";
import { useLocation } from "preact-iso";

export default function LoginForm() {
	const ws = useContext(WSContext);
	const messageRef = useRef<HTMLParagraphElement>(null);
	const { route } = useLocation();

	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await ws.request({
			action: "login",
			payload: entries
		});

		const { user_id, token, message, success } = res.payload;

		messageRef!.current!.innerHTML = message;

		if (!success || !token) {
			return;
		}

		Storage.token.value = token;
		Storage.user.value = {
			id: user_id,
			name: formData.get("name")!.toString()
		};
		route("/lobby");
	}

	return <>
		<section id="login">
			<h2>Login</h2>
			<form action="/api/v1/login" method="post" onSubmit={onSubmit}>
				<label>
					<p>Username: </p>
					<input type="text" name="name" required />
				</label>

				<label>
					<p>Password:</p>
					<input type="password" name="password" required />
				</label>

				<input type="submit" />

				<p ref={messageRef}></p>
			</form>

			<p>Don't have an account? <a href="/register">Register</a></p>
		</section>
	</>;
}
