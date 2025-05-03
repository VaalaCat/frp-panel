package utils

import "github.com/lucasepe/codename"

func NewCodeName(tokenLength int) string {
	rng, _ := codename.DefaultRNG()
	return codename.Generate(rng, tokenLength)
}
