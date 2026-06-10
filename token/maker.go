// Package token provides token creation and verification for authentication.
package token

import (
	"time"
)

// Maker is an interface for creating and verifying tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
