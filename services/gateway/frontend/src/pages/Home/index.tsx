import AuthGuard from "../../components/auth/AuthGuard";
import './style.css';

export function Home() {
	return <>
		<h1>codenames</h1>

		<a href="/login">login</a>
		<a href="/register">register</a>

		<AuthGuard>
			<a href="/lobby">lobby</a>
		</AuthGuard>
	</>;
}
