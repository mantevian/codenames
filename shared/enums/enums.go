package enums

type team string

const (
	TeamRed  team = "red"
	TeamBlue team = "blue"
)

type Team enumBase[team]

type gameStatus string

const (
	Waiting  gameStatus = "waiting"
	Playing  gameStatus = "playing"
	Finished gameStatus = "finished"
)

type GameStatus enumBase[gameStatus]
