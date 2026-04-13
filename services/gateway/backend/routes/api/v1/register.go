package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func Register(c fiber.Ctx, db *sql.DB) error {
	args := c.Request().PostArgs()
	name := string(args.Peek("name"))
	password := string(args.Peek("password"))
	passwordConfirm := string(args.Peek("password_confirm"))

	if len(name) == 0 {
		return c.Status(400).SendString("no name")
	}

	if len(password) < 6 {
		return c.Status(400).SendString("password too short")
	}

	if passwordConfirm != password {
		return c.Status(400).SendString("passwords don't match")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).SendString("error")
	}

	_, err = db.Exec(
		"insert into users (id, name, password_hash, created_at) values (gen_random_uuid(), $1, $2, now())",
		name,
		string(passwordHash),
	)

	if err != nil {
		return c.Status(500).SendString("error")
	}

	return c.SendString("ok")
}
