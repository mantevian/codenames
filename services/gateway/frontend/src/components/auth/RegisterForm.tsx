import { TargetedEvent } from "preact";
import { useContext, useRef } from "preact/hooks";
import { WSContext } from "../WebSocketProvider";
import { useLocation } from "preact-iso";

export default function RegisterForm() {
	const ws = useContext(WSContext);
	const messageRef = useRef<HTMLParagraphElement>(null);
	const { route } = useLocation();

	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await ws.request({
			action: "register",
			payload: entries
		});

		const { success, message } = res.payload;

		messageRef!.current!.innerHTML = message;

		if (success) {
			route("/login");
		}
	}

	return <>
		<section id="register">
			<h2>Register</h2>
			<form action="/api/v1/register" method="post" onSubmit={onSubmit}>
				<label>
					<p>Username: </p>
					<input type="text" name="name" required />
				</label>

				<label>
					<p>Password: </p>
					<input type="password" name="password" required />
				</label>

				<label>
					<p>Repeat password: </p>
					<input type="password" name="password_confirm" required />
				</label>
				
				<input type="submit" />
				<p ref={messageRef}></p>
			</form>

			<p>Already have an account? <a href="/login">Login</a></p>
		</section>
	</>;
}
