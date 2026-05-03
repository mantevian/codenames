package util

import (
	"math/rand"

	"mantevian.xyz/codenames/shared/enums"
)

func MakeShuffledTileList(red, blue, neutral, assassin int) []enums.Tile {
	list := make([]enums.Tile, 0, 25)

	for range red {
		list = append(list, enums.TileRed)
	}

	for range blue {
		list = append(list, enums.TileBlue)
	}

	for range neutral {
		list = append(list, enums.TileNeutral)
	}

	for range assassin {
		list = append(list, enums.TileAssassin)
	}

	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
