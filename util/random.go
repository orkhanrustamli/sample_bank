package util

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomName() string {
	return gofakeit.Name()
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "AZN"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
