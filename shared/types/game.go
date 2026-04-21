package types

import (
	"time"

	"mantevian.xyz/codenames/shared/enums"
)

type BasicGameResponse struct {
	Id           Uuid             `json:"id"`
	StartingTeam enums.Team       `json:"starting_team"`
	JoinCode     JoinCode         `json:"join_code"`
	Language     enums.Language   `json:"language"`
	Status       enums.GameStatus `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
}
