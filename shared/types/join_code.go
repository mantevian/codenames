package types

import (
	"math/rand/v2"
)

type JoinCode string

func NewJoinCode() JoinCode {
	a := byte(rand.Int()%26 + 'a')
	b := byte(rand.Int()%26 + 'a')
	c := byte(rand.Int()%26 + 'a')
	d := byte(rand.Int()%26 + 'a')

	return JoinCode(string([]byte{a, b, c, d}))
}
