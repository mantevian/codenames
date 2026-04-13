import amqplib from "amqplib";

const connection = await amqplib.connect("amqp://localhost:5672");

const channel = await connection.createChannel();

await channel.assertQueue("test");

channel.consume("test", (msg) => console.log(msg?.content.toString()));
