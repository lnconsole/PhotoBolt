package shared

import (
	"math/rand"
)

func RandInt(max int) int {
	return rand.Intn(max)
}
