import { TargetedEvent } from "preact";

export default function LoginForm() {
	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await fetch(form.action, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(entries)
		});

		const { token } = await res.json();
		localStorage.setItem("token", token);
	}

	return <>
		<form action="/api/v1/login" method="post" onSubmit={onSubmit}>
			<input type="text" name="name" required />
			<input type="password" name="password" required />
			<input type="submit" />
		</form>
	</>;
}