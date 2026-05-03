import { useLocation } from 'preact-iso';
import AuthGuard from "./auth/AuthGuard";

export function Header() {
	const { url } = useLocation();

	return (
		<header>
			<nav>
				<a href="/" class={url == '/' ? 'active' : ''}>
					Home
				</a>
				
				<AuthGuard>
					<p>Welcome, {sessionStorage.getItem("username")}</p>
				</AuthGuard>
			</nav>
		</header>
	);
}
