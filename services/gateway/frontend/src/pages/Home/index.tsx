import AuthGuard from "../../components/auth/AuthGuard";
import './style.css';

export function Home() {
	return <section class="home">
		<h1>CODENAMES</h1>
		
		<nav>
			<a href="/login">Login</a>
			<a href="/register">Register</a>

			<AuthGuard>
				<a href="/lobby">Lobby</a>
			</AuthGuard>
		</nav>
	</section>;
}
