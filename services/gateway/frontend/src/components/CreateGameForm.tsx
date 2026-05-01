import { TargetedEvent } from "preact";
import { useContext } from "preact/hooks";
import { WSContext } from "./WebSocketProvider";

export default function CreateGameForm() {
	const ws = useContext(WSContext);

	async function onSubmit(e: TargetedEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);
		const entries = Object.fromEntries(formData.entries());

		const res = await ws.request({
			action: "create_game",
			payload: entries
		});
	}

	return <>
		<section id="create-game">
			<h2>Create game</h2>
			<form action="/api/v1/create_game" method="post" onSubmit={onSubmit}>
				<select name="language">
					<option value="en">English</option>
					<option value="ru">Русский</option>
				</select>
				<input type="submit" />
			</form>
		</section>
	</>;
}
