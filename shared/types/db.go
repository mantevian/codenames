package types

import (
	"database/sql"
	"time"

	"mantevian.xyz/codenames/shared/enums"
)

type Game struct {
	Id           string           `json:"id"`
	StartingTeam enums.Team       `json:"starting_team"`
	JoinCode     string           `json:"join_code"`
	Language     string           `json:"language"`
	TeamWon      enums.Team       `json:"team_won"`
	Status       enums.GameStatus `json:"status"`
	FinishedAt   sql.NullTime     `json:"finished_at"`
	CreatedAt    time.Time        `json:"created_at"`
}

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
