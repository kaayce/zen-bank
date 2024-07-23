package token

import "time"

// Maker interface to manage tokens
type Maker interface {
	// Creates and Signs token valid for specific duration for user
	CreateToken(username string, duration time.Duration) (string, error)
	// Checks if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
