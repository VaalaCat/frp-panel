package utils

import (
	"math/rand"
)

func RandomInt(a, b int) int {
	if a > b {
		a, b = b, a
	}
	return rand.Intn(b-a+1) + a
}
