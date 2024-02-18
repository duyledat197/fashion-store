// Package service is representation of
package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	msgpb "trintech/review/dto/msg/common"
	pb "trintech/review/dto/user-management/auth"
	"trintech/review/internal/user-management/entity"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server/xcontext"
	postgresclient "trintech/review/pkg/postgres_client"
	"trintech/review/pkg/pubsub"
	"trintech/review/pkg/token_util"
)

// AuthService is representation of
type AuthService interface {
}

type authService struct {
	userRepo interface {
		RetrieveByEmail(context.Context, string) (*entity.User, error)
		RetrieveByUserName(context.Context, string) (*entity.User, error)
		Create(context.Context, database.Executor, *entity.User) (int64, error)
		UpdatePassword(ctx context.Context, email, password string) error
	}

	loginHistoryRepo interface {
		Create(context.Context, database.Executor, *entity.LoginHistory) error
		UpdateLogout(ctx context.Context, db database.Executor, accessToken string) error
	}

	userCacheRepo interface {
		RetrieveByUserName(context.Context, string) (*entity.User, error)
		StoreByUserName(context.Context, string, *entity.User) error
		RemoveByUserName(context.Context, string) error

		RetrieveByEmail(context.Context, string) (*entity.User, error)
		StoreByEmail(context.Context, string, *entity.User) error
		RemoveByEmail(context.Context, string) error

		IncrementForgotPassword(ctx context.Context, email string, duration time.Duration) (int64, error)

		StoreResetToken(ctx context.Context, email string, resetToken string, duration time.Duration) error
		RetrieveResetToken(ctx context.Context, email string, resetToken string) error
		RemoveByResetToken(context.Context, string, string) error
	}

	tknGenerator token_util.Authenticator[*xcontext.UserInfo]
	db           database.Database

	publisher pubsub.Publisher

	pb.UnimplementedAuthServiceServer
}

// NewAuthService is representation of
func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) retrieveUserByUserName(ctx context.Context, userName string) (*entity.User, error) {
	user, err := s.userCacheRepo.RetrieveByUserName(ctx, userName)
	if user == nil {
		user, err = s.userRepo.RetrieveByUserName(ctx, userName)
	}

	return user, err
}

func (s *authService) retrieveUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userCacheRepo.RetrieveByEmail(ctx, email)
	if user == nil {
		user, err = s.userRepo.RetrieveByEmail(ctx, email)
	}

	return user, err
}

// Register implements
func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.retrieveUserByUserName(ctx, req.GetUserName())
	switch {
	case user != nil:
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	case !errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.Internal, "unable to retrieve user by username: %v", err.Error())
	}

	user, err = s.retrieveUserByEmail(ctx, req.GetEmail())
	switch {
	case user != nil:
		return nil, status.Errorf(codes.AlreadyExists, "email already register")
	case !errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.Internal, "unable to retrieve user by email: %v", err.Error())
	}

	pwd, err := crypto_util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to hash password")
	}

	id, err := s.userRepo.Create(ctx, s.db, &entity.User{
		UserName: postgresclient.PgTypeText(req.GetUserName()),
		Password: postgresclient.PgTypeText(pwd),
		Name:     postgresclient.PgTypeText(req.GetName()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create user: %v", err.Error())
	}

	return &pb.RegisterResponse{
		UserId: id,
	}, nil
}

// Login implements login business logic.
func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.retrieveUserByUserName(ctx, req.GetUserName())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.InvalidArgument, "username or password is not correctly")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err.Error())
	}

	if err := crypto_util.CheckPassword(req.Password, user.Password.String); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password is not valid")
	}

	tkn, err := s.tknGenerator.Generate(&xcontext.UserInfo{
		UserID: user.ID,
		Role:   string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err.Error())
	}

	session, err := xcontext.ExtractSessionFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unable to extract session: %v", err.Error())
	}

	if err := s.loginHistoryRepo.Create(ctx, s.db, &entity.LoginHistory{
		UserID:      user.ID,
		IP:          session.IP,
		AccessToken: tkn,
		UserAgent:   session.UserAgent,
		LoginAt:     time.Now(),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create login history: %v", err.Error())
	}

	if err := s.userCacheRepo.StoreByUserName(ctx, user.UserName.String, user); err != nil {
		slog.Error("unable to store user cache: %w", err)
	}

	return &pb.LoginResponse{
		UserId:      user.ID,
		AccessToken: tkn,
	}, nil
}

func (s *authService) Logout(ctx context.Context, _ *emptypb.Empty) (*pb.LogoutResponse, error) {
	session, err := xcontext.ExtractSessionFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unable to logout history: %v", err.Error())
	}

	if err := s.loginHistoryRepo.UpdateLogout(ctx, s.db, session.AccessToken); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update logout: %v", err.Error())
	}

	return &pb.LogoutResponse{}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	user, err := s.retrieveUserByEmail(ctx, req.GetEmail())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.InvalidArgument, "username or password is not correctly")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err.Error())
	}

	count, err := s.userCacheRepo.IncrementForgotPassword(ctx, user.Email.String, time.Hour)
	if count > 5 {
		return nil, status.Errorf(codes.Internal, "forgot password count got exceed")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to increment forgot password count: %v", err.Error())
	}

	resetToken := crypto_util.GeneratePassword(15, true, true, true)
	if err := s.userCacheRepo.StoreResetToken(ctx, req.GetEmail(), resetToken, 5*time.Minute); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to store forgot password token: %v", err.Error())
	}

	go func() {
		data, err := proto.Marshal(&msgpb.ForgotPassword{
			UserName:   user.UserName.String,
			Email:      user.Email.String,
			Name:       user.Name.String,
			ResetToken: resetToken,
		})
		if err != nil {
			slog.Error("unable to marshal data", "err", err.Error())
			return
		}
		if err := s.publisher.Publish(ctx, "FORGOT_PASSWORD", []byte(req.GetEmail()), data); err != nil {
			slog.Error("unable to publish forgot password message", "err", err.Error())
		}
	}()

	return &pb.ForgotPasswordResponse{}, nil
}

func (s *authService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if req.GetNewPassword() != req.GetRepeatPassword() {
		return nil, status.Errorf(codes.InvalidArgument, "new password and repeated password is not match")
	}

	if err := s.userCacheRepo.RetrieveResetToken(ctx, req.GetEmail(), req.GetResetToken()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to retrieve forgot password token: %v", err.Error())
	}

	pwd, err := crypto_util.HashPassword(req.GetNewPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to hash password")
	}

	if err := s.userRepo.UpdatePassword(ctx, req.GetEmail(), pwd); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update password: %v", err.Error())
	}

	if err := s.userCacheRepo.RemoveByResetToken(ctx, req.GetEmail(), req.GetResetToken()); err != nil {
		slog.Error("unable to remove reset token", "errpr", err)
	}

	return &pb.ResetPasswordResponse{}, nil

}
