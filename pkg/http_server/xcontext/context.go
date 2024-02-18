package xcontext

import (
	"context"
	"fmt"
	"time"
)

type UserInfo struct {
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Session struct {
	IP          string `json:"ip,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

func (p *UserInfo) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return fmt.Errorf("token has been expired")
	}

	return nil
}

func (p *UserInfo) AddExpired(expirationTime time.Duration) {
	p.ExpiredAt = time.Now().Add(expirationTime)
}

// ImportUserInfoToContext inject the user info which retrieved from token
// into the given context.
func ImportUserInfoToContext(ctx context.Context, info *UserInfo) context.Context {
	return context.WithValue(ctx, &userInfoKey{}, info)
}

// ExtractUserInfoFromContext returns an user info which was injected from [ImportUserInfoToContext].
func ExtractUserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	info, ok := ctx.Value(&userInfoKey{}).(*UserInfo)

	if !ok || info == nil {
		return nil, fmt.Errorf("authorization is not valid")
	}

	return info, nil
}

// ImportSessionToContext .
func ImportSessionToContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, &sessionKey{}, session)
}

// ExtractSessionFromContext returns an.
func ExtractSessionFromContext(ctx context.Context) (*Session, error) {
	info, ok := ctx.Value(&sessionKey{}).(*Session)

	if !ok || info == nil {
		return nil, fmt.Errorf("session is not valid")
	}

	return info, nil
}
