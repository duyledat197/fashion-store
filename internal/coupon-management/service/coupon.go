// Package service ...
package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trintech/review/dto/coupon-management/coupon"
	"trintech/review/internal/coupon-management/entity"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/pg_util"
)

type couponService struct {
	couponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error)
		DeleteByID(ctx context.Context, db database.Executor, id int64) error
	}

	userCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
	}

	productCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
	}

	db database.Database
	pb.UnimplementedCouponServiceServer
}

func (s *couponService) CreateCoupon(ctx context.Context, req *pb.CreateCouponRequest) (*pb.CreateCouponResponse, error) {
	var (
		id  int64
		err error
	)
	if err := s.db.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		id, err = s.couponRepo.Create(ctx, tx, &entity.Coupon{
			Code:        pg_util.PgTypeText(crypto_util.GenerateCode("COUPON")),
			From:        pg_util.PgTypeTimestamptz(req.From.AsTime()),
			To:          pg_util.PgTypeTimestamptz(req.To.AsTime()),
			Rules:       req.GetRules(),
			Description: pg_util.PgTypeText(req.GetDescription()),
			ImageURL:    pg_util.PgTypeText(req.GetImageUrl()),
			CreatedBy:   pgtype.Int8{},
			CreatedAt:   pgtype.Timestamptz{},
			UpdatedAt:   pgtype.Timestamptz{},
		})
		if err != nil {
			return fmt.Errorf("unable to create coupon: %w", err)
		}

		switch req.GetApplyId().(type) {
		case *pb.CreateCouponRequest_UserId:
			if err := s.userCouponRepo.Create(ctx, tx, &entity.UserCoupon{
				CouponID: pg_util.PgTypeInt8(id),
				UserID:   pg_util.PgTypeInt8(req.GetUserId().GetValue()),
			}); err != nil {
				return fmt.Errorf("unable to create user coupon: %w", err)
			}
		case *pb.CreateCouponRequest_ProductId:
			if err := s.productCouponRepo.Create(ctx, tx, &entity.ProductCoupon{
				CouponID:  pg_util.PgTypeInt8(id),
				ProductID: pg_util.PgTypeInt8(req.GetProductId().GetValue()),
			}); err != nil {
				return fmt.Errorf("unable to create user coupon: %w", err)
			}
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create coupon: %v", err)
	}

	return &pb.CreateCouponResponse{
		Id: id,
	}, nil
}

func (s *couponService) DeleteCouponByID(ctx context.Context, req *pb.DeleteCouponByIDRequest) (*pb.DeleteCouponByIDResponse, error) {
	if err := s.couponRepo.DeleteByID(ctx, s.db, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete coupon by id: %v", err.Error())
	}

	return &pb.DeleteCouponByIDResponse{}, nil
}
