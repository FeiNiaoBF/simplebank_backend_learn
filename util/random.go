package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

var (
	rng      *rand.Rand
	seed     = time.Now().UnixNano()
	currency = []string{"USD", "EUR", "CNY", "CAD"}
)

// init random seed
func init() {
	// After go 1.20, use New(NewSource(seed))
	// but unsafe for goroutine
	rng = rand.New(rand.NewSource(seed))
}

// RandomInt generates a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	size := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(size)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomAmount generates a random amount of negative or positive money
func RandomAmount() int64 {
	return RandomInt(-1000, 1000)
}
// RandomCurrency generates a random currency
func RandomCurrency() string {
	n := len(currency)
	return currency[rng.Intn(n)]
}
