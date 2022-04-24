package util

import "math/rand"

// RandIntn returns a random number in [0, n).
func RandIntn(n int) int {
	return rand.Intn(n) // nolint:gosec
}
