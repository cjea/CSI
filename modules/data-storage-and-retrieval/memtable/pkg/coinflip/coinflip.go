package coinflip

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Flip() bool {
	return rand.Intn(4) == 0
}
