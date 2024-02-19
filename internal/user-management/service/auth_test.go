// service
package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trintech/review/dto/user-management/auth"
	"trintech/review/internal/user-management/entity"
	"trintech/review/mocks"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
	"trintech/review/pkg/pubsub"
	"trintech/review/pkg/token_util"
)

func Test_authService_Register(t *testing.T) {
	type fields struct {
		tknGenerator                   token_util.Authenticator[*xcontext.UserInfo]
		db                             database.Database
		publisher                      pubsub.Publisher
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		userCacheRepo                  *mocks.UserCacheRepository
		userRepo                       *mocks.UserRepository
		loginHistoryRepo               *mocks.LoginHistoryRepository
	}
	type args struct {
		ctx context.Context
		req *pb.RegisterRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RegisterResponse
		wantErr error
		setup   func(ctx context.Context, fields fields)
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "password",
				},
			},
			want: &pb.RegisterResponse{
				UserId: 1,
			},
			wantErr: nil,
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

				fields.userRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
			},
		},
		{
			name: "err password is not match",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "wrong-password",
				},
			},
			wantErr: status.Errorf(codes.InvalidArgument, "password and repeated password is not match"),
			setup: func(ctx context.Context, fields fields) {
			},
		},
		{
			name: "err username existed in cache",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "password",
				},
			},
			want:    nil,
			wantErr: status.Errorf(codes.AlreadyExists, "username already exists"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
			},
		},
		{
			name: "err username existed in db",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "password",
				},
			},
			wantErr: status.Errorf(codes.AlreadyExists, "username already exists"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
			},
		},
		{
			name: "err email existed in cache",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "password",
				},
			},
			wantErr: status.Errorf(codes.AlreadyExists, "email already register"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
			},
		},
		{
			name: "err email existed in db",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RegisterRequest{
					UserName:       "user-name",
					Password:       "password",
					Name:           "test",
					Email:          "user@gmail.com",
					RepeatPassword: "password",
				},
			},
			want:    nil,
			wantErr: status.Errorf(codes.AlreadyExists, "email already register"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.args.ctx, tt.fields)

			s := &authService{
				userRepo:         tt.fields.userRepo,
				loginHistoryRepo: tt.fields.loginHistoryRepo,
				userCacheRepo:    tt.fields.userCacheRepo,

				tknGenerator:                   tt.fields.tknGenerator,
				db:                             tt.fields.db,
				publisher:                      tt.fields.publisher,
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
			}
			got, err := s.Register(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	type fields struct {
		tknGenerator                   *mocks.Authenticator[*xcontext.UserInfo]
		db                             database.Database
		publisher                      pubsub.Publisher
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		userCacheRepo                  *mocks.UserCacheRepository
		userRepo                       *mocks.UserRepository
		loginHistoryRepo               *mocks.LoginHistoryRepository
	}
	type args struct {
		ctx context.Context
		req *pb.LoginRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.LoginResponse
		wantErr error
		setup   func(ctx context.Context, fields fields)
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				tknGenerator:     &mocks.Authenticator[*xcontext.UserInfo]{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					UserName: "user-name",
					Password: "password",
				},
			},

			setup: func(ctx context.Context, fields fields) {
				pwd, _ := crypto_util.HashPassword("password")
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       1,
					UserName: pg_util.PgTypeText("user-name"),
					Password: pg_util.PgTypeText(pwd),
				}, nil)
				fields.tknGenerator.On("Generate", mock.Anything, mock.Anything).Return("token", nil)
				fields.loginHistoryRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.userCacheRepo.On("StoreByUserName", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},

		{
			name: "err wrong password",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				tknGenerator:     &mocks.Authenticator[*xcontext.UserInfo]{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					UserName: "user-name",
					Password: "wrong-password",
				},
			},
			wantErr: status.Errorf(codes.InvalidArgument, "username or password is not correctly"),
			setup: func(ctx context.Context, fields fields) {
				pwd, _ := crypto_util.HashPassword("password")
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       1,
					UserName: pg_util.PgTypeText("user-name"),
					Password: pg_util.PgTypeText(pwd),
				}, nil)
			},
		},

		{
			name: "err user not exist",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				tknGenerator:     &mocks.Authenticator[*xcontext.UserInfo]{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					UserName: "user-name",
					Password: "password",
				},
			},
			wantErr: status.Errorf(codes.InvalidArgument, "username or password is not correctly"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
		},

		{
			name: "err store user to cache fail",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				tknGenerator:     &mocks.Authenticator[*xcontext.UserInfo]{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.LoginRequest{
					UserName: "user-name",
					Password: "password",
				},
			},
			setup: func(ctx context.Context, fields fields) {
				pwd, _ := crypto_util.HashPassword("password")
				fields.userCacheRepo.On("RetrieveByUserName", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       1,
					UserName: pg_util.PgTypeText("user-name"),
					Password: pg_util.PgTypeText(pwd),
				}, nil)
				fields.tknGenerator.On("Generate", mock.Anything, mock.Anything).Return("token", nil)
				fields.loginHistoryRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.userCacheRepo.On("StoreByUserName", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("something wrong"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.args.ctx, tt.fields)
			s := &authService{
				userRepo:         tt.fields.userRepo,
				loginHistoryRepo: tt.fields.loginHistoryRepo,
				userCacheRepo:    tt.fields.userCacheRepo,
				tknGenerator:     tt.fields.tknGenerator,
			}
			_, err := s.Login(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_authService_ForgotPassword(t *testing.T) {
	type fields struct {
		tknGenerator                   *mocks.Authenticator[*xcontext.UserInfo]
		db                             database.Database
		publisher                      *mocks.Publisher
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		userCacheRepo                  *mocks.UserCacheRepository
		userRepo                       *mocks.UserRepository
		loginHistoryRepo               *mocks.LoginHistoryRepository
	}
	type args struct {
		ctx context.Context
		req *pb.ForgotPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ForgotPasswordResponse
		wantErr error
		setup   func(ctx context.Context, fields fields)
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				publisher:        &mocks.Publisher{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{
					Email: "user@gmail.com",
				},
			},
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
				fields.userCacheRepo.On("IncrementForgotPassword", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
				fields.userCacheRepo.On("StoreResetToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "err email is not exist",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{
					Email: "user@gmail.com",
				},
			},
			wantErr: status.Errorf(codes.InvalidArgument, "email is not correctly"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				fields.userRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
			},
		},
		{
			name: "err forgot password count exceed",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{
					Email: "user@gmail.com",
				},
			},
			wantErr: status.Errorf(codes.Internal, "forgot password count got exceed"),
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
				fields.userCacheRepo.On("IncrementForgotPassword", mock.Anything, mock.Anything, mock.Anything).Return(int64(6), nil)
			},
		},
		{
			name: "err publish msg failed",
			fields: fields{
				userCacheRepo:    &mocks.UserCacheRepository{},
				userRepo:         &mocks.UserRepository{},
				loginHistoryRepo: &mocks.LoginHistoryRepository{},
				publisher:        &mocks.Publisher{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ForgotPasswordRequest{
					Email: "user@gmail.com",
				},
			},
			setup: func(ctx context.Context, fields fields) {
				fields.userCacheRepo.On("RetrieveByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
				fields.userCacheRepo.On("IncrementForgotPassword", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
				fields.userCacheRepo.On("StoreResetToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("somegthing went wrong"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.args.ctx, tt.fields)
			s := &authService{
				userRepo:         tt.fields.userRepo,
				loginHistoryRepo: tt.fields.loginHistoryRepo,
				userCacheRepo:    tt.fields.userCacheRepo,
				tknGenerator:     tt.fields.tknGenerator,
				db:               tt.fields.db,
				publisher:        tt.fields.publisher,
			}
			_, err := s.ForgotPassword(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
