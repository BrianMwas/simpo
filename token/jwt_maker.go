package token

import (
	"errors"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt"
	"time"
)

const minSecretKeyLen = 32

// JWTMaker is a json web token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWT maker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLen)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwt := jwt2.NewWithClaims(jwt2.SigningMethodHS256, payload)
	return jwt.SignedString([]byte(jwtMaker.secretKey))
}

// VerifyToken verifies whether the token is valid or not
func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt2.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt2.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(jwtMaker.secretKey), nil
	}

	jwtToken, err := jwt2.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(jwt2.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
