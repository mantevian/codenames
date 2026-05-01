import { TargetedEvent } from "preact";
import { useContext } from "preact/hooks";
import { WSContext } from "../WebSocketProvider";

export default function RegisterForm() {
	const ws = useContext(WSContext);

	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await ws.request({
			action: "register",
			payload: entries
		});
	}

	return <>
		<section id="register">
			<h2>Register</h2>
			<form action="/api/v1/register" method="post" onSubmit={onSubmit}>
				<input type="text" name="name" required />
				<input type="password" name="password" required />
				<input type="password" name="password_confirm" required />
				<input type="submit" />
			</form>
		</section>
	</>;
}
