// Package service ...
package service

import (
	"context"
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	couponpb "trintech/review/dto/coupon-management/coupon"
	pb "trintech/review/dto/product-management/product"
	"trintech/review/internal/product-management/entity"
	userEntity "trintech/review/internal/user-management/entity"
	"trintech/review/mocks"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
	"trintech/review/pkg/postgres_client"
)

func Test_productService_PurchaseProduct(t *testing.T) {
	type fields struct {
		productRepo          *mocks.ProductRepository
		purchasedProductRepo *mocks.PurchasedProductRepository

		db                  *postgres_client.PostgresClient
		couponServiceClient *mocks.CouponServiceClient
	}

	db, smock, _ := sqlmock.New()
	type args struct {
		ctx context.Context
		req *pb.PurchaseProductRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.PurchaseProductResponse
		wantErr error
		setup   func(ctx context.Context, fields fields)
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			fields: fields{
				productRepo:          &mocks.ProductRepository{},
				purchasedProductRepo: &mocks.PurchasedProductRepository{},
				db: &postgres_client.PostgresClient{
					DB: db,
				},
				couponServiceClient: &mocks.CouponServiceClient{},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), http_server.ImportUserInfoToMD(&xcontext.UserInfo{
					UserID: 1,
					Role:   userEntity.UserRole_User,
				})),
				req: &pb.PurchaseProductRequest{
					Id:     1,
					Coupon: wrapperspb.String("ABC"),
				},
			},
			want: &pb.PurchaseProductResponse{},
			setup: func(ctx context.Context, fields fields) {
				fields.productRepo.On("RetrieveByID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Product{
					ID:    pg_util.NullInt64(1),
					Price: pg_util.NullFloat64(100),
				}, nil)

				fields.couponServiceClient.On("RetrieveCouponByCode", mock.Anything, mock.Anything, mock.Anything).
					Return(&couponpb.RetrieveCouponByCodeResponse{
						CanUse:       true,
						Value:        50,
						DiscountType: couponpb.DiscountType_DiscountType_VALUE,
					}, nil)

				smock.ExpectBegin()
				fields.purchasedProductRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.couponServiceClient.On("ApplyCoupon", mock.Anything, mock.Anything, mock.Anything).Return(&couponpb.ApplyCouponResponse{}, nil)
				smock.ExpectCommit()
			},
		},
		{
			name: "err invalid user",
			fields: fields{
				productRepo:          &mocks.ProductRepository{},
				purchasedProductRepo: &mocks.PurchasedProductRepository{},
				db: &postgres_client.PostgresClient{
					DB: db,
				},
				couponServiceClient: &mocks.CouponServiceClient{},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &pb.PurchaseProductResponse{},
			wantErr: status.Errorf(codes.PermissionDenied, "user doesn't have permission"),
			setup: func(ctx context.Context, fields fields) {
			},
		},
		{
			name: "err product not found",
			fields: fields{
				productRepo:          &mocks.ProductRepository{},
				purchasedProductRepo: &mocks.PurchasedProductRepository{},
				db: &postgres_client.PostgresClient{
					DB: db,
				},
				couponServiceClient: &mocks.CouponServiceClient{},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), http_server.ImportUserInfoToMD(&xcontext.UserInfo{
					UserID: 1,
					Role:   userEntity.UserRole_User,
				})),
				req: &pb.PurchaseProductRequest{
					Id:     1,
					Coupon: wrapperspb.String("ABC"),
				},
			},
			want:    &pb.PurchaseProductResponse{},
			wantErr: status.Errorf(codes.NotFound, "product not found"),
			setup: func(ctx context.Context, fields fields) {
				fields.productRepo.On("RetrieveByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

			},
		},
		{
			name: "err cannot use coupon",
			fields: fields{
				productRepo:          &mocks.ProductRepository{},
				purchasedProductRepo: &mocks.PurchasedProductRepository{},
				db: &postgres_client.PostgresClient{
					DB: db,
				},
				couponServiceClient: &mocks.CouponServiceClient{},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), http_server.ImportUserInfoToMD(&xcontext.UserInfo{
					UserID: 1,
					Role:   userEntity.UserRole_User,
				})),
				req: &pb.PurchaseProductRequest{
					Id:     1,
					Coupon: wrapperspb.String("ABC"),
				},
			},
			want:    &pb.PurchaseProductResponse{},
			wantErr: status.Errorf(codes.NotFound, "coupon not found"),
			setup: func(ctx context.Context, fields fields) {
				fields.productRepo.On("RetrieveByID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Product{
					ID:    pg_util.NullInt64(1),
					Price: pg_util.NullFloat64(100),
				}, nil)

				fields.couponServiceClient.On("RetrieveCouponByCode", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.NotFound, "coupon not found"))
			},
		},

		{
			name: "err cannot apply coupon",
			fields: fields{
				productRepo:          &mocks.ProductRepository{},
				purchasedProductRepo: &mocks.PurchasedProductRepository{},
				db: &postgres_client.PostgresClient{
					DB: db,
				},
				couponServiceClient: &mocks.CouponServiceClient{},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), http_server.ImportUserInfoToMD(&xcontext.UserInfo{
					UserID: 1,
					Role:   userEntity.UserRole_User,
				})),
				req: &pb.PurchaseProductRequest{
					Id:     1,
					Coupon: wrapperspb.String("ABC"),
				},
			},
			want:    &pb.PurchaseProductResponse{},
			wantErr: status.Errorf(codes.FailedPrecondition, "unable to apply this coupon"),
			setup: func(ctx context.Context, fields fields) {
				fields.productRepo.On("RetrieveByID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Product{
					ID:    pg_util.NullInt64(1),
					Price: pg_util.NullFloat64(100),
				}, nil)

				fields.couponServiceClient.On("RetrieveCouponByCode", mock.Anything, mock.Anything, mock.Anything).
					Return(&couponpb.RetrieveCouponByCodeResponse{
						CanUse:       true,
						Value:        50,
						DiscountType: couponpb.DiscountType_DiscountType_VALUE,
					}, nil)

				smock.ExpectBegin()
				fields.purchasedProductRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.couponServiceClient.On("ApplyCoupon", mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.FailedPrecondition, "unable to apply coupon"))
				smock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.args.ctx, tt.fields)
			s := &productService{
				productRepo:          tt.fields.productRepo,
				purchasedProductRepo: tt.fields.purchasedProductRepo,
				db:                   tt.fields.db,
				couponServiceClient:  tt.fields.couponServiceClient,
			}
			_, err := s.PurchaseProduct(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
