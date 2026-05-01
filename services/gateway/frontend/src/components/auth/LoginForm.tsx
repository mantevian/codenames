import { TargetedEvent } from "preact";
import { useContext } from "preact/hooks";
import { WSContext } from "../WebSocketProvider";

export default function LoginForm() {
	const ws = useContext(WSContext);

	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await ws.request({
			action: "login",
			payload: entries
		});

		const { token } = res.payload;
		localStorage.setItem("token", token);
	}

	return <>
		<section id="login">
			<h2>Login</h2>
			<form action="/api/v1/login" method="post" onSubmit={onSubmit}>
				<input type="text" name="name" required />
				<input type="password" name="password" required />
				<input type="submit" />
			</form>
		</section>
	</>;
}
