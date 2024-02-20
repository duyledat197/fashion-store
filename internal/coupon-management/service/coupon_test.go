// Package service ...
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pb "trintech/review/dto/coupon-management/coupon"
	"trintech/review/internal/coupon-management/entity"
	"trintech/review/mocks"
	"trintech/review/pkg/database"
	"trintech/review/pkg/pg_util"
)

func Test_couponService_RetrieveCouponByCode_CanUse(t *testing.T) {
	type fields struct {
		couponRepo        *mocks.CouponRepository
		userCouponRepo    *mocks.UserCouponRepository
		productCouponRepo *mocks.ProductCouponRepository
		usedCouponRepo    *mocks.UsedCouponRepository
		db                database.Database
	}
	type args struct {
		ctx context.Context
		req *pb.RetrieveCouponByCodeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RetrieveCouponByCodeResponse
		wantErr error
		setup   func(ctx context.Context, fields fields)
	}{
		// TODO: Add test cases.
		{
			name: "happy case not check",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code: "ABC",
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: true,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type:  pg_util.NullString(pb.CouponType_CouponType_LIMITED.String()),
					Used:  pg_util.NullInt64(1),
					Total: pg_util.NullInt64(10),
				}, nil)

			},
		},
		{
			name: "happy case check use type limited",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: true,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type:  pg_util.NullString(pb.CouponType_CouponType_LIMITED.String()),
					Used:  pg_util.NullInt64(1),
					Total: pg_util.NullInt64(10),
				}, nil)

			},
		},

		{
			name: "happy case check use type user",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: true,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type: pg_util.NullString(pb.CouponType_CouponType_USER.String()),
				}, nil)
				fields.userCouponRepo.On("RetrieveByCouponID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.UserCoupon{
					Used:  pg_util.NullInt64(1),
					Total: pg_util.NullInt64(10),
				}, nil)
			},
		},
		{
			name: "happy case check use type product",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: true,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type: pg_util.NullString(pb.CouponType_CouponType_PRODUCT.String()),
				}, nil)
				fields.productCouponRepo.On("RetrieveByCouponID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.ProductCoupon{
					Used:  pg_util.NullInt64(1),
					Total: pg_util.NullInt64(10),
				}, nil)
			},
		},

		{
			name: "err exceed type limited exceed",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type:  pg_util.NullString(pb.CouponType_CouponType_LIMITED.String()),
					Used:  pg_util.NullInt64(10),
					Total: pg_util.NullInt64(10),
				}, nil)
			},
		},

		{
			name: "err exceed type user exceed",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Type: pg_util.NullString(pb.CouponType_CouponType_LIMITED.String()),
				}, nil)

				fields.productCouponRepo.On("RetrieveByCouponID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.UserCoupon{
					Used:  pg_util.NullInt64(10),
					Total: pg_util.NullInt64(10),
				}, nil)
			},
		},
		{
			name: "err exceed type product exceed",
			fields: fields{
				couponRepo:        &mocks.CouponRepository{},
				userCouponRepo:    &mocks.UserCouponRepository{},
				productCouponRepo: &mocks.ProductCouponRepository{},
				usedCouponRepo:    &mocks.UsedCouponRepository{},
			},
			args: args{
				ctx: context.Background(),
				req: &pb.RetrieveCouponByCodeRequest{
					Code:     "ABC",
					CheckUse: true,
				},
			},
			want: &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			},
			setup: func(ctx context.Context, fields fields) {
				fields.couponRepo.On("RetrieveByCode", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Coupon{
					Used:  pg_util.NullInt64(1),
					Total: pg_util.NullInt64(10),
					Type:  pg_util.NullString(pb.CouponType_CouponType_PRODUCT.String()),
				}, nil)
				fields.productCouponRepo.On("RetrieveByCouponID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.ProductCoupon{
					Used:  pg_util.NullInt64(10),
					Total: pg_util.NullInt64(10),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		tt.setup(tt.args.ctx, tt.fields)
		t.Run(tt.name, func(t *testing.T) {
			s := &couponService{
				couponRepo:        tt.fields.couponRepo,
				userCouponRepo:    tt.fields.userCouponRepo,
				productCouponRepo: tt.fields.productCouponRepo,
				usedCouponRepo:    tt.fields.usedCouponRepo,
				db:                tt.fields.db,
			}
			got, err := s.RetrieveCouponByCode(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.CanUse, got.CanUse)
			}
		})
	}
}
