package common

import (
	"strconv"
	"time"
)

// get a random int
func Rand(n int) int {
	seed := time.Now().Unix()
	rand64 := seed % int64(n)
	strInt64 := strconv.FormatInt(rand64, 10)
	rand, _ := strconv.Atoi(strInt64)
	return rand
}
