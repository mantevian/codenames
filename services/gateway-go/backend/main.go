package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()

	app.Get("/*", static.New("../frontend/dist", static.Config{
		Compress:      true,
		MaxAge:        0,
		CacheDuration: -1,
	}))

	log.Fatal(app.Listen(":8080"))
}
