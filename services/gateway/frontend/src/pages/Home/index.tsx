import './style.css';

export function Home() {
	return (
		<>
			<p>register</p>
			<form action="/api/v1/register" method="post">
				<input type="text" name="name" required />
				<input type="password" name="password" required />
				<input type="password" name="password_confirm" required />
				<input type="submit" />
			</form>

			<p>login</p>
			<form action="/api/v1/login" method="post">
				<input type="text" name="name" required />
				<input type="password" name="password" required />
				<input type="submit" />
			</form>
		</>
	);
}
