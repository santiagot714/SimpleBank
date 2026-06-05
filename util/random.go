package util

import (
	"math/rand"
	"time"
	"github.com/shopspring/decimal"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() decimal.Decimal {
	return decimal.NewFromInt(RandomInt(0, 1000))
}

func RandomCurrency() string {
	return []string{"USD", "EUR", "CAD", "GBP", "COP"}[rand.Intn(5)]
}