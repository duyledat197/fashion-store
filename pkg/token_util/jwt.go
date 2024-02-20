package token_util

import (
	"errors"
	"fmt"
	"time"

	"github.com/reddit/jwt-go"

	"trintech/review/pkg/http_server/xcontext"
)

// JWTAuthenticator is representation of [Authenticator] engine that implement using JWT.
type JWTAuthenticator struct {
	secretKey string
}

func NewJWTAuthenticator(secretKey string) (JWTAuthenticator, error) {
	return JWTAuthenticator{
		secretKey,
	}, nil
}

func (a *JWTAuthenticator) Generate(payload *xcontext.UserInfo, expirationTime time.Duration) (string, error) {
	payload.AddExpired(expirationTime)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %w", err)
	}

	return token, nil

}

func (a *JWTAuthenticator) Verify(token string) (*xcontext.UserInfo, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token is not valid")
		}
		return []byte(a.secretKey), nil
	}
	var claims xcontext.UserInfo
	jwtToken, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, fmt.Errorf("")) {
			return &claims, fmt.Errorf("token is not valid")
		}

		return &claims, fmt.Errorf("token is not valid: %w", err)
	}

	payload, ok := jwtToken.Claims.(*xcontext.UserInfo)
	if !ok {
		return &claims, fmt.Errorf("token is not valid")
	}

	return payload, nil
}
