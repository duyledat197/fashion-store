package token_util

import "time"

// Authenticator is a representation of token generator that implement generate and verify.
type Authenticator[T Claims] interface {
	Generate(payload T, expirationTime time.Duration) (string, error)
	Verify(token string) (T, error)
}

// Claims is representation of info inject into tokens.
type Claims interface {
	Valid() error
	AddExpired(time.Duration)
}
