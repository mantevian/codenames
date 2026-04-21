package types

import (
	"database/sql"
	"time"

	"mantevian.xyz/codenames/shared/enums"
)

type Game struct {
	Id           Uuid             `json:"id"`
	StartingTeam enums.Team       `json:"starting_team"`
	JoinCode     JoinCode         `json:"join_code"`
	Language     enums.Language   `json:"language"`
	TeamWon      enums.Team       `json:"team_won"`
	Status       enums.GameStatus `json:"status"`
	FinishedAt   sql.NullTime     `json:"finished_at"`
	CreatedAt    time.Time        `json:"created_at"`
}

type User struct {
	Id        Uuid      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Player struct {
	Id      Uuid       `json:"id"`
	UserId  Uuid       `json:"user_id"`
	GameId  Uuid       `json:"game_id"`
	Team    enums.Team `json:"team"`
	Role    enums.Role `json:"Role"`
	IsReady bool       `json:"is_ready"`
}

type Tile struct {
	Position   int        `json:"position"`
	GameId     Uuid       `json:"game_id"`
	Type       enums.Tile `json:"type"`
	IsRevealed bool       `json:"is_revealed"`
	Word       string     `json:"word"`
}

type Turn struct {
	Id         Uuid      `json:"id"`
	PlayerId   Uuid      `json:"player_id"`
	GameId     Uuid      `json:"game_id"`
	ClueWord   string    `json:"clue_word"`
	ClueNumber int       `json:"clue_number"`
	CreatedAt  time.Time `json:"created_at"`
}

type Guess struct {
	TurnId    Uuid      `json:"turn_id"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type Word struct {
	Language enums.Language `json:"language"`
	Word     string         `json:"word"`
}
