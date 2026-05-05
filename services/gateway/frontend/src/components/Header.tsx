import { useLocation } from 'preact-iso';
import AuthGuard from "./auth/AuthGuard";
import { Storage } from "../storage/user";

export function Header() {
	const { url } = useLocation();

	return (
		<header>
			<nav>
				<a href="/" class={url == '/' ? 'active' : ''}>
					Home
				</a>
				
				<AuthGuard>
					<p>{Storage.user.value?.name}</p>
				</AuthGuard>
			</nav>
		</header>
	);
}
