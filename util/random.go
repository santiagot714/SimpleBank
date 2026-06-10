package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"

	"github.com/shopspring/decimal"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomInt generates a random integer between min and max.
func RandomInt(minRange, maxRange int64) int64 {
	return minRange + mathrand.Int63n(maxRange-minRange+1)
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[mathrand.Intn(len(alphabet))]
	}
	return string(b)
}

// RandomSecretKey generates a cryptographically random 32-byte seed encoded as hex (64 chars).
// NewPasetoMaker uses this seed to derive a valid ed25519 keypair.
func RandomSecretKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
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
	return []string{USD, EUR, GBP, CAD, CHF, NZD, AUD, COP}[mathrand.Intn(8)]
}

// RandomEmail generates a random email.
func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}
