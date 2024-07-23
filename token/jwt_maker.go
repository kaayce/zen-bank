package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// adds methods for golang-jwt v5
type JWTPayloadClaims struct {
	Payload
	jwt.RegisteredClaims
}

// JSON Web TOken Maker implements Maker interface
type JWTMaker struct {
	secretKey string
}

// Creates a new JWT Maker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// factory method to create NewJWTPayloadClaims
func NewJWTPayloadClaims(payload *Payload) *JWTPayloadClaims {
	return &JWTPayloadClaims{
		Payload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "zenbank",
			Subject:   payload.UserName,
			Audience:  []string{"clients"},
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(payload.IssuedAt),
			ID:        payload.ID.String(),
		},
	}
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewJWTPayloadClaims(payload))
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtClaims := &JWTPayloadClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, ErrInvalidToken) {
			return nil, ErrInvalidToken
		} else {
			return nil, err
		}
	}
	payloadClaims, ok := jwtToken.Claims.(*JWTPayloadClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return &payloadClaims.Payload, nil
}
