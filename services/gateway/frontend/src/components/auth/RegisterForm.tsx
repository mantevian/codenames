import { TargetedEvent } from "preact";

export default function RegisterForm() {
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