import LoginForm from "../../components/register/LoginForm";
import RegisterForm from "../../components/register/RegisterForm";
import './style.css';

export function Home() {
	return <>
		<RegisterForm />
		<LoginForm />
	</>;
}
