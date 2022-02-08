package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomName() string {
	return gofakeit.Name()
}

func UsernameFromFullname(fullname string) string {
	nameSurname := strings.Split(strings.ToLower(fullname), " ")
	return strings.Join(nameSurname, ".")
}

func EmailFromUsername(username string) string {
	return fmt.Sprintf("%s@%s", username, gofakeit.DomainName())
}

func RandomPassword() string {
	return gofakeit.Password(true, true, true, true, true, 10)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{AZN, USD, EUR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomString(n uint) string {
	return gofakeit.LetterN(n)
}
