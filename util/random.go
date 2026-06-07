// Package util contains utility functions for random data generation.
package util

import (
	"math/rand"

	"github.com/shopspring/decimal"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomInt generates a random integer between min and max.
func RandomInt(minRange, maxRange int64) int64 {
	return minRange + rand.Int63n(maxRange-minRange+1)
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

// RandomOwner generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random money amount.
func RandomMoney() decimal.Decimal {
	return decimal.NewFromInt(RandomInt(0, 1000))
}

// RandomCurrency generates a random currency.
func RandomCurrency() string {
	return []string{"USD", "EUR", "CAD", "GBP", "COP"}[rand.Intn(5)]
}
