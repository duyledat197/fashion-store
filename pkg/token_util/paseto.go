package token_util

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoAuthenticator is representation of [Authenticator] engine that implement using paseto.
type PasetoAuthenticator[T Claims] struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoAuthenticator[T Claims](symmetricKey string) (Authenticator[T], error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symmetricKey must have at least 32 bytes")
	}

	return &PasetoAuthenticator[T]{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (a *PasetoAuthenticator[T]) Generate(payload T, expirationTime time.Duration) (string, error) {
	payload.AddExpired(expirationTime)
	token, err := a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %w", err)
	}

	return token, nil
}

func (a *PasetoAuthenticator[T]) Verify(token string) (T, error) {
	var payload T

	if err := a.paseto.Decrypt(token, a.symmetricKey, &payload, nil); err != nil {
		return payload, fmt.Errorf("token is not valid: %w", err)
	}

	if err := payload.Valid(); err != nil {
		return payload, fmt.Errorf("token is not valid: %w", err)
	}

	return payload, nil
}
