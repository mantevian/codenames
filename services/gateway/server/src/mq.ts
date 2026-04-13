import amqplib from "amqplib";

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost:5672";

const connection = await amqplib.connect(RABBITMQ_URL);

const channel = await connection.createChannel();

await channel.assertQueue("test");

channel.sendToQueue(
	"test",
	Buffer.from(
		JSON.stringify({
			hello: "world",
		}),
	),
);

await channel.close();

await connection.close();
