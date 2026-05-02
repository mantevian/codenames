import CreateGameForm from "../../components/CreateGameForm";
import GameList from "../../components/GameList";

export function Lobby() {
	return <>
		<h1>lobby</h1>
		<CreateGameForm />
		<GameList />
	</>
}
