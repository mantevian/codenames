import LoginForm from "../../components/auth/LoginForm";
import RegisterForm from "../../components/auth/RegisterForm";
import CreateGameForm from "../../components/CreateGameForm";
import WaitinGameList from "../../components/WaitingGameList";
import './style.css';

export function Home() {
	return <>
		<RegisterForm />
		<LoginForm />
		<CreateGameForm />
		<WaitinGameList />
	</>;
}
