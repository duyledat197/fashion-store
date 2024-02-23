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
	memcache "trintech/review/internal/user-management/repository/cache"
	"trintech/review/internal/user-management/repository/postgres"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
	"trintech/review/pkg/pubsub"
	"trintech/review/pkg/token_util"
)

// AuthService is representation of
type AuthService interface {
}

type authService struct {
	userRepo interface {
		RetrieveByEmail(context.Context, database.Executor, string) (*entity.User, error)
		RetrieveByUserName(context.Context, database.Executor, string) (*entity.User, error)
		Create(context.Context, database.Executor, *entity.User) (int64, error)
		UpdatePassword(ctx context.Context, db database.Executor, email, password string) error
	}

	loginHistoryRepo interface {
		Create(context.Context, database.Executor, *entity.LoginHistory) error
		UpdateLogout(ctx context.Context, db database.Executor, userID int64, accessToken string) error
	}

	userCacheRepo interface {
		RetrieveByUserName(context.Context, string) (*entity.User, error)
		StoreByUserName(context.Context, string, *entity.User) error
		RemoveByUserName(context.Context, string) error

		RetrieveByEmail(context.Context, string) (*entity.User, error)
		StoreByEmail(context.Context, string, *entity.User) error
		RemoveByEmail(context.Context, string) error

		IncrementForgotPassword(ctx context.Context, email string) (int64, error)

		StoreResetToken(ctx context.Context, email string, resetToken string) error
		IsExistResetToken(ctx context.Context, email string, resetToken string) error
		RemoveByResetToken(context.Context, string, string) error
	}

	tknGenerator token_util.JWTAuthenticator
	db           database.Database

	publisher pubsub.Publisher

	pb.UnimplementedAuthServiceServer
}

// NewAuthService is representation of
func NewAuthService(
	db database.Database,
	publisher pubsub.Publisher,
	tknGenerator token_util.JWTAuthenticator,
) pb.AuthServiceServer {
	return &authService{
		db:               db,
		publisher:        publisher,
		tknGenerator:     tknGenerator,
		userRepo:         postgres.NewUserRepository(),
		loginHistoryRepo: postgres.NewLoginHistoryRepository(),
		userCacheRepo:    memcache.NewUserCacheRepository(),
	}
}

func (s *authService) retrieveUserByUserName(ctx context.Context, userName string) (*entity.User, error) {
	user, err := s.userCacheRepo.RetrieveByUserName(ctx, userName)
	if user == nil {
		user, err = s.userRepo.RetrieveByUserName(ctx, s.db, userName)
	}

	return user, err
}

func (s *authService) retrieveUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userCacheRepo.RetrieveByEmail(ctx, email)
	if user == nil {
		user, err = s.userRepo.RetrieveByEmail(ctx, s.db, email)
	}

	return user, err
}

// Register is a method of the authService that handles user registration.
// It validates the registration request, checks for existing usernames and emails,
// hashes the password, and creates a new user in the repository.
func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate if the password and repeated password match
	if req.Password != req.RepeatPassword {
		return nil, status.Errorf(codes.InvalidArgument, "password and repeated password do not match")
	}

	// Check if a user with the given username already exists
	user, err := s.retrieveUserByUserName(ctx, req.GetUserName())
	switch {
	case user != nil:
		// If a user with the same username exists, return an already exists error
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	case !errors.Is(err, sql.ErrNoRows):
		// If there is an internal error during username retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve user by username: %v", err.Error())
	}

	// Check if a user with the given email already exists
	user, err = s.retrieveUserByEmail(ctx, req.GetEmail())
	switch {
	case user != nil:
		// If a user with the same email exists, return an already exists error
		return nil, status.Errorf(codes.AlreadyExists, "email already registered")
	case !errors.Is(err, sql.ErrNoRows):
		// If there is an internal error during email retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve user by email: %v", err.Error())
	}

	// Hash the provided password
	pwd, err := crypto_util.HashPassword(req.GetPassword())
	if err != nil {
		// If there is an internal error during password hashing, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to hash password")
	}

	// Create a new user in the repository
	id, err := s.userRepo.Create(ctx, s.db, &entity.User{
		UserName: pg_util.NullString(req.GetUserName()),
		Email:    pg_util.NullString(req.GetEmail()),
		Password: pg_util.NullString(pwd),
		Name:     pg_util.NullString(req.GetName()),
		Role:     entity.UserRole_User,
	})
	if err != nil {
		// If there is an internal error during user creation, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to create user: %v", err.Error())
	}

	// Return the user ID in the response
	return &pb.RegisterResponse{
		UserId: id,
	}, nil
}

// Login implements login business logic.
// Login is a method of the authService that handles user login.
// It retrieves the user by username, checks the password, generates an access token,
// records login history, and stores the user information in cache.
func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Retrieve the user by username
	user, err := s.retrieveUserByUserName(ctx, req.GetUserName())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// If the user with the given username is not found, return an invalid argument error
		return nil, status.Errorf(codes.InvalidArgument, "username or password is not correct")
	case err != nil:
		// If there is an internal error during user retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err.Error())
	}

	// Check if the provided password matches the hashed password in the database
	if err := crypto_util.CheckPassword(req.Password, user.Password.String); err != nil {
		// If the password is incorrect, return an invalid argument error
		return nil, status.Errorf(codes.InvalidArgument, "username or password is not correct")
	}

	// Generate an access token for the user
	tkn, err := s.tknGenerator.Generate(&xcontext.UserInfo{
		UserID: user.ID.Int64,
		Role:   string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		// If there is an internal error during token generation, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to generate token: %v", err.Error())
	}

	// Extract session information from the context
	session := http_server.ExtractSessionFromCtx(ctx)

	// Record login history for the user
	if err := s.loginHistoryRepo.Create(ctx, s.db, &entity.LoginHistory{
		UserID:      pg_util.NullInt64(user.ID.Int64),
		IP:          pg_util.NullString(session.IP),
		AccessToken: pg_util.NullString(tkn),
		UserAgent:   pg_util.NullString(session.UserAgent),
	}); err != nil {
		// If there is an internal error during login history creation, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to create login history: %v", err.Error())
	}

	// Store user information in cache for future reference
	if err := s.userCacheRepo.StoreByUserName(ctx, user.UserName.String, user); err != nil {
		// Log an error if there is an issue storing user information in cache
		slog.Error("unable to store user cache", "err", err)
	}

	// Return the user ID and access token in the response
	return &pb.LoginResponse{
		UserId:      user.ID.Int64,
		AccessToken: tkn,
	}, nil
}

// Logout is a method of the authService that handles user logout.
// It updates the logout timestamp in the login history repository.
func (s *authService) Logout(ctx context.Context, _ *emptypb.Empty) (*pb.LogoutResponse, error) {
	// Extract session information from the context
	session := http_server.ExtractSessionFromCtx(ctx)

	// Extract user information from the context
	userCtx, _ := http_server.ExtractUserInfoFromCtx(ctx)

	// Update the logout timestamp in the login history repository
	if err := s.loginHistoryRepo.UpdateLogout(ctx, s.db, userCtx.UserID, session.AccessToken); err != nil {
		// If there is an internal error during logout update, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to update logout: %v", err.Error())
	}

	// Return an empty response indicating successful logout
	return &pb.LogoutResponse{}, nil
}

// ForgotPassword is a method of the authService that handles the process of requesting
// a password reset. It retrieves the user by email, increments the forgot password count,
// generates a reset token, stores the reset token in the cache, and publishes a message
// for further processing.
func (s *authService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	// Retrieve the user by email
	user, err := s.retrieveUserByEmail(ctx, req.GetEmail())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// If the user with the given email is not found, return an invalid argument error
		return nil, status.Errorf(codes.InvalidArgument, "email is not correct")
	case err != nil:
		// If there is an internal error during user retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve user: %v", err.Error())
	}

	// Increment the forgot password count and check if it exceeds the limit
	count, err := s.userCacheRepo.IncrementForgotPassword(ctx, user.Email.String)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to increment forgot password count: %v", err.Error())
	}

	if count > 5 {
		// If the forgot password count exceeds the limit, return an internal server error
		return nil, status.Errorf(codes.Internal, "forgot password count exceeded")
	}

	// Generate a reset token and store it in the cache
	resetToken := crypto_util.GeneratePassword(15, true, true, true)
	if err := s.userCacheRepo.StoreResetToken(ctx, req.GetEmail(), resetToken); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to store forgot password token: %v", err.Error())
	}

	// Asynchronously publish a message for further processing (e.g., sending an email)
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

	// Return an empty response indicating successful initiation of the forgot password process
	return &pb.ForgotPasswordResponse{}, nil
}

// ResetPassword is a method of the authService that handles the process of resetting a user's password.
// It checks if the new password and repeated password match, validates the reset token,
// hashes the new password, updates the user's password in the repository, and removes the reset token from the cache.
func (s *authService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// Check if the new password and repeated password match
	if req.GetNewPassword() != req.GetRepeatPassword() {
		return nil, status.Errorf(codes.InvalidArgument, "new password and repeated password is not match")
	}

	// Check if the reset token is valid
	if err := s.userCacheRepo.IsExistResetToken(ctx, req.GetEmail(), req.GetResetToken()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to retrieve forgot password token: %v", err.Error())
	}

	// Hash the new password
	pwd, err := crypto_util.HashPassword(req.GetNewPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to hash password")
	}

	// Update the user's password in the repository
	if err := s.userRepo.UpdatePassword(ctx, s.db, req.GetEmail(), pwd); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update password: %v", err.Error())
	}

	// Remove the reset token from the cache
	if err := s.userCacheRepo.RemoveByResetToken(ctx, req.GetEmail(), req.GetResetToken()); err != nil {
		slog.Error("unable to remove reset token", "error", err)
	}

	// Return an empty response indicating successful password reset
	return &pb.ResetPasswordResponse{}, nil
}
