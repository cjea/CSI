package coinflip

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

func Flip() bool {
	return rand.Intn(2) == 0
}
