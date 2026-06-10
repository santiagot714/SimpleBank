package token

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// Error cases
var (
	ErrInvalidSecretKeyLength = errors.New("invalid secret key length")
	ErrInvalidToken           = errors.New("invalid token")
	ErrExpiredToken           = errors.New("token has expired")
	ErrInvalidTokenID         = errors.New("invalid token id")
	ErrInvalidTokenSubject    = errors.New("invalid token subject")
	ErrInvalidTokenIssuedAt   = errors.New("invalid token issued at")
	ErrInvalidTokenExpiration = errors.New("invalid token expiration")
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto paseto.V4AsymmetricSecretKey
}

// NewPasetoMaker creates a new PasetoMaker from a 32-byte seed encoded as a 64-char hex string.
func NewPasetoMaker(secretKey string) (Maker, error) {
	// secretKey must be a 32-byte seed encoded as hex (64 hex chars).
	// We derive the full ed25519 private key from the seed so that the
	// public key portion is always consistent with the seed.
	if len(secretKey) != 64 {
		return nil, ErrInvalidSecretKeyLength
	}
	seedBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	key, err := paseto.NewV4AsymmetricSecretKeyFromBytes(privateKey)
	if err != nil {
		return nil, err
	}
	return &PasetoMaker{paseto: key}, nil
}

// CreateToken generates a signed PASETO v4 token for the given username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()
	token.SetJti(payload.ID.String())
	token.SetSubject(payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	signedToken := token.V4Sign(maker.paseto, nil)
	return signedToken, payload, nil
}

// VerifyToken parses and validates a PASETO v4 token, returning its payload.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parsedToken, err := parser.ParseV4Public(maker.paseto.Public(), token, nil)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	idStr, err := parsedToken.GetJti()
	if err != nil {
		return nil, ErrInvalidTokenID
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ErrInvalidTokenSubject
	}

	username, err := parsedToken.GetSubject()
	if err != nil {
		return nil, ErrInvalidTokenIssuedAt
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidTokenIssuedAt
	}

	expiredAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, ErrInvalidTokenExpiration
	}

	return &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
