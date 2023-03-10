package utils

import (
	"math/rand"
	"time"
)

const (
	max = 9999999
	min = 1000000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInteger() int {
	return rand.Intn(max-min+1) + min
}
