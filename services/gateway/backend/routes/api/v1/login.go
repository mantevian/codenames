package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func Login(c fiber.Ctx, db *sql.DB) error {
	args := c.Request().PostArgs()
	name := string(args.Peek("name"))
	password := string(args.Peek("password"))

	if len(name) == 0 {
		return c.Status(400).SendString("no name")
	}

	if len(password) == 0 {
		return c.Status(400).SendString("no password")
	}

	rows, err := db.Query(
		"select password_hash from users where name = $1",
		name,
	)

	if err != nil {
		return c.Status(500).SendString("error")
	}

	defer rows.Close()

	var hash string
	rows.Next()
	err = rows.Scan(&hash)

	if err != nil {
		// user not found
		return c.Status(500).SendString("error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		// wrong password
		return c.Status(500).SendString("error")
	}

	return c.SendString("ok")
}
