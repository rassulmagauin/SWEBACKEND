package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific user
	CreateToken(username string, role string, duration time.Duration) (string, error)
	// VerifyToken verifies the token and returns the username
	VerifyToken(token string) (*Payload, error)
}
