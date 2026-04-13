import express from "express";

export const PORT = process.env.PORT || 9001;
export const app = express();

export function startWeb() {
	app.use(express.json());

	app.listen(PORT, () => {
		console.log(`server started at ${PORT}`);
	});

	app.use("/", express.static("../client/dist/"));

	// app.get("/", (req, res) => {
	// 	res.sendFile("./client/dist/index.html");
	// });
}
