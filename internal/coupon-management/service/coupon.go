// Package service ...
package service

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "trintech/review/dto/coupon-management/coupon"
	"trintech/review/internal/coupon-management/entity"
	userEntity "trintech/review/internal/user-management/entity"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
)

type couponService struct {
	couponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error)
		DeleteByID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCode(ctx context.Context, db database.Executor, code string) (*entity.Coupon, error)
	}

	userCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.UserCoupon, error)
	}

	productCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.ProductCoupon, error)
	}

	usedCouponRepo interface {
		ListUsedCouponByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entity.CouponUsedCoupon, error)
	}

	db database.Database
	pb.UnimplementedCouponServiceServer
}

func validAdmin(ctx context.Context) (*xcontext.UserInfo, error) {
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	if !slices.Contains([]string{userEntity.UserRole_Admin, userEntity.UserRole_SuperAdmin}, userCtx.Role) {
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	return userCtx, nil
}

func (s *couponService) CreateCoupon(ctx context.Context, req *pb.CreateCouponRequest) (*pb.CreateCouponResponse, error) {
	userCtx, err := validAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var id int64

	if err := database.Transaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		id, err = s.couponRepo.Create(ctx, tx, &entity.Coupon{
			Code:        pg_util.NullString(crypto_util.GenerateCode("COUPON")),
			From:        pg_util.NullTime(req.From.AsTime()),
			To:          pg_util.NullTime(req.To.AsTime()),
			Rules:       req.GetRules(),
			Description: pg_util.NullString(req.GetDescription()),
			ImageURL:    pg_util.NullString(req.GetImageUrl()),
			CreatedBy:   pg_util.NullInt64(userCtx.UserID),
		})
		if err != nil {
			return fmt.Errorf("unable to create coupon: %w", err)
		}

		switch req.GetApplyId().(type) {
		case *pb.CreateCouponRequest_UserId:
			if err := s.userCouponRepo.Create(ctx, tx, &entity.UserCoupon{
				CouponID: pg_util.NullInt64(id),
				UserID:   pg_util.NullInt64(req.GetUserId().GetValue()),
			}); err != nil {
				return fmt.Errorf("unable to create user coupon: %w", err)
			}
		case *pb.CreateCouponRequest_ProductId:
			if err := s.productCouponRepo.Create(ctx, tx, &entity.ProductCoupon{
				CouponID:  pg_util.NullInt64(id),
				ProductID: pg_util.NullInt64(req.GetProductId().GetValue()),
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
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	if err := s.couponRepo.DeleteByID(ctx, s.db, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete coupon by id: %v", err.Error())
	}

	return &pb.DeleteCouponByIDResponse{}, nil
}

func (s *couponService) RetrieveCouponByCode(ctx context.Context, req *pb.RetrieveCouponByCodeRequest) (*pb.RetrieveCouponByCodeResponse, error) {
	coupon, err := s.couponRepo.RetrieveByCode(ctx, s.db, req.GetCode())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve coupon by code: %v", err.Error())
	}

	switch coupon.Type.String {
	case pb.CouponType_CouponType_LIMITED.String():
		if coupon.Used.Int64 >= coupon.Total.Int64 {
			return &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			}, nil
		}
	case pb.CouponType_CouponType_USER.String():
		userCoupon, err := s.userCouponRepo.RetrieveByCouponID(ctx, s.db, coupon.ID.Int64)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to retrieve user coupon : %v", err.Error())
		}
		if userCoupon.Used.Int64 >= userCoupon.Total.Int64 {
			return &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			}, nil
		}
	case pb.CouponType_CouponType_PRODUCT.String():
		productCoupon, err := s.productCouponRepo.RetrieveByCouponID(ctx, s.db, coupon.ID.Int64)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to retrieve user coupon : %v", err.Error())
		}
		if productCoupon.Used.Int64 >= productCoupon.Total.Int64 {
			return &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			}, nil
		}
	}

	return &pb.RetrieveCouponByCodeResponse{
		Type:         new(pb.CouponType).FromString(coupon.Type.String),
		Total:        coupon.Total.Int64,
		From:         timestamppb.New(coupon.From.Time),
		To:           timestamppb.New(coupon.To.Time),
		ImageUrl:     coupon.ImageURL.String,
		Description:  coupon.Description.String,
		Rules:        coupon.Rules,
		DiscountType: new(pb.DiscountType).FromString(coupon.DiscountType.String),
		Value:        coupon.Value.Float64,
		Used:         coupon.Used.Int64,
		CanUse:       true,
	}, nil
}

func (s *couponService) ListUsedCoupon(ctx context.Context, _ *pb.ListUsedCouponRequest) (*pb.ListUsedCouponResponse, error) {
	userCtx, err := xcontext.ExtractUserInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user from context: %v", err.Error())
	}
	coupons, err := s.usedCouponRepo.ListUsedCouponByUserID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve used coupon: %v", err.Error())
	}

	respData := make([]*pb.ListUsedCouponResponse_Coupon, 0, len(coupons))
	for _, coupon := range coupons {
		respData = append(respData, &pb.ListUsedCouponResponse_Coupon{
			Code:        coupon.Coupon.Code.String,
			Description: coupon.Coupon.Description.String,
			ImageUrl:    coupon.Coupon.ImageURL.String,
			ApplyAt:     timestamppb.New(coupon.UsedCoupon.CreatedAt.Time),
		})
	}

	return &pb.ListUsedCouponResponse{
		Data: respData,
	}, nil
}
