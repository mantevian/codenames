import LoginForm from "../../components/auth/LoginForm";
import RegisterForm from "../../components/auth/RegisterForm";
import './style.css';

export function Home() {
	return <>
		<RegisterForm />
		<LoginForm />
	</>;
}
