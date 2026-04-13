package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Generic() fiber.Handler {
	return static.New("../frontend/dist", static.Config{
		Compress:      true,
		MaxAge:        0,
		CacheDuration: -1,
	})
}
