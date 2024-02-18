package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/token_util"
)

type AuthService interface {
}

type authService struct {
	userRepo interface {
		RetrieveByUserName(context.Context, string) (*entity.User, error)
		Create(context.Context, database.Executor, *entity.User) (int64, error)
	}

	loginHistoryRepo interface {
		Create(context.Context, database.Executor, *entity.LoginHistory) error
		UpdateLogout(ctx context.Context, db database.Executor, accessToken string) error
	}

	userCacheRepo interface {
		RetrieveByUserName(context.Context, string) (*entity.User, error)
		Store(context.Context, string, *entity.User) error
	}

	tknGenerator token_util.Authenticator[*xcontext.UserInfo]

	db database.Database
}

func NewAuthService() AuthService {
	return authService{}
}

func (s *authService) Login(ctx context.Context, req *entity.User) (*entity.User, string, error) {
	user, _ := s.userCacheRepo.RetrieveByUserName(ctx, req.UserName)
	if user == nil {
		var err error
		user, err = s.userRepo.RetrieveByUserName(ctx, req.UserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, "", fmt.Errorf("username or password is not correctly")
			}
			return nil, "", err
		}
	}

	if err := crypto_util.CheckPassword(req.Password, user.Password); err != nil {
		return nil, "", err
	}

	tkn, err := s.tknGenerator.Generate(&xcontext.UserInfo{
		UserID: user.ID,
		Role:   string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	session, err := xcontext.ExtractSessionFromContext(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("unable to extract session: %w", err)

	}

	if err := s.loginHistoryRepo.Create(ctx, s.db, &entity.LoginHistory{
		UserID:      user.ID,
		IP:          session.IP,
		AccessToken: tkn,
		UserAgent:   session.UserAgent,
		LoginAt:     time.Now(),
	}); err != nil {
		return nil, "", fmt.Errorf("unable to create login history: %w", err)
	}

	if err := s.userCacheRepo.Store(ctx, user.UserName, user); err != nil {
		slog.Error("unable to store user cache: %w", err)
	}

	return user, tkn, nil
}

func (s *authService) Logout(ctx context.Context, req *entity.User) (*entity.User, error) {
	session, err := xcontext.ExtractSessionFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to logout history: %w", err)
	}

	if err := s.loginHistoryRepo.UpdateLogout(ctx, s.db, session.AccessToken); err != nil {
		return nil, fmt.Errorf("unable to update logout: %w", err)
	}

}
