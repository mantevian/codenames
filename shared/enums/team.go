package enums

import "math/rand/v2"

type Team string

const (
	TeamRed  Team = "red"
	TeamBlue Team = "blue"
)

func RandomTeam() Team {
	if rand.Int()%2 == 0 {
		return TeamRed
	} else {
		return TeamBlue
	}
}
