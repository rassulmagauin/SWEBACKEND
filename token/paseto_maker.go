package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, it should be %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific user
func (p *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

// VerifyToken verifies the token and returns the username
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	if err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil); err != nil {
		return nil, errors.New("invalid token")
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
