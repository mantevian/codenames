import { TargetedEvent } from "preact";

export default function LoginForm() {
	function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		fetch(form.action, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(entries)
		});
	}

	return <>
		<form action="/api/v1/login" method="post" onSubmit={onSubmit}>
			<input type="text" name="name" required />
			<input type="password" name="password" required />
			<input type="submit" />
		</form>
	</>;
}