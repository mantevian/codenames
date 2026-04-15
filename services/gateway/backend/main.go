package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
	"mantevian.xyz/codenames/gateway/routes"
	v1 "mantevian.xyz/codenames/gateway/routes/api/v1"
)

type TestRow struct {
	id   int
	text string
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(503).SendString("Database unavailable")
		}
		return c.SendString("OK")
	})

	app.Get("/test", func(c fiber.Ctx) error {
		rows, err := db.Query("select * from test")
		if err != nil {
			return c.Status(503).SendString("Database unavailable")
		}
		defer rows.Close()

		var result strings.Builder

		for rows.Next() {
			var row TestRow
			if err := rows.Scan(&row.id, &row.text); err != nil {
				return c.Status(500).SendString("Error reading data")
			}
			result.WriteString(row.text)
			result.WriteString("\n")
		}

		return c.SendString(result.String())
	})

	app.Post("/api/v1/register", func(c fiber.Ctx) {
		v1.Register(c, db)
	})

	app.Post("/api/v1/login", func(c fiber.Ctx) {
		v1.Login(c, db)
	})

	app.Get("/*", routes.Generic())

	log.Fatal(app.Listen(":8080"))
}
