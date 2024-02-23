// Package service ...
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "trintech/review/dto/coupon-management/coupon"
	"trintech/review/internal/coupon-management/entity"
	"trintech/review/internal/coupon-management/repository/postgres"
	userEntity "trintech/review/internal/user-management/entity"
	"trintech/review/pkg/crypto_util"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
)

// couponService provides coupon handling operations.
type couponService struct {
	couponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.Coupon) (int64, error)
		DeleteByID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCode(ctx context.Context, db database.Executor, code string) (*entity.Coupon, error)
	}

	userCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.UserCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCouponIDUserID(ctx context.Context, db database.Executor, couponID, userID int64) (*entity.UserCoupon, error)
	}

	productCouponRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.ProductCoupon) error
		DeleteByCouponID(ctx context.Context, db database.Executor, id int64) error
		RetrieveByCouponID(ctx context.Context, db database.Executor, couponID int64) (*entity.ProductCoupon, error)
	}

	usedCouponRepo interface {
		ListUsedCouponByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entity.CouponUsedCoupon, error)
		Create(ctx context.Context, db database.Executor, data *entity.UsedCoupon) error
	}

	db database.Database
	pb.UnimplementedCouponServiceServer
}

// NewCouponService returns coupon service that implements coupon handling operations.
func NewCouponService(db database.Database) pb.CouponServiceServer {
	return &couponService{
		db:                db,
		couponRepo:        postgres.NewCouponRepository(),
		productCouponRepo: postgres.NewProductCouponRepository(),
		userCouponRepo:    postgres.NewUserCouponRepository(),
		usedCouponRepo:    postgres.NewUsedCouponRepository(),
	}
}

// validAdmin checks if the user in the provided context is an admin or super admin.
// If the user is an admin or super admin, it returns the user information.
// If not, it returns a permission denied error.
func validAdmin(ctx context.Context) (*xcontext.UserInfo, error) {
	// Extract user information from the context
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		// If user information cannot be extracted, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// Check if the user has the role of admin or superadmin
	if !slices.Contains([]string{userEntity.UserRole_Admin, userEntity.UserRole_SuperAdmin}, userCtx.Role) {
		// If the user does not have the required role, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// If the user is an admin or superadmin, return the user information
	return userCtx, nil
}

// CreateCoupon is a method of the couponService that creates a new coupon based on the provided request.
// It performs various validations and database transactions to ensure the integrity of the data.
func (s *couponService) CreateCoupon(ctx context.Context, req *pb.CreateCouponRequest) (*pb.CreateCouponResponse, error) {
	// Validate if the user is an admin
	userCtx, err := validAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var id int64

	// Start a database transaction
	if err := database.Transaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// Create a new coupon using the coupon repository
		id, err = s.couponRepo.Create(ctx, tx, &entity.Coupon{
			Code:         pg_util.NullString(crypto_util.GenerateCode("COUPON")),
			From:         pg_util.NullTime(req.From.AsTime()),
			To:           pg_util.NullTime(req.To.AsTime()),
			Description:  pg_util.NullString(req.GetDescription()),
			ImageURL:     pg_util.NullString(req.GetImageUrl()),
			CreatedBy:    pg_util.NullInt64(userCtx.UserID),
			Value:        pg_util.NullFloat64(req.GetValue()),
			Total:        pg_util.NullInt64(req.GetTotal()),
			CreatedAt:    pg_util.NullTime(time.Now()),
			Type:         pg_util.NullString(req.GetType().String()),
			DiscountType: pg_util.NullString(req.GetDiscountType().String()),
		})
		if err != nil {
			return fmt.Errorf("unable to create coupon: %w", err)
		}

		// Depending on the ApplyId type, associate the coupon with a user or a product
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
				return fmt.Errorf("unable to create product coupon: %w", err)
			}
		}

		return nil
	}); err != nil {
		// If there's an error during the transaction, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to create coupon: %v", err)
	}

	// Return the ID of the created coupon in the response
	return &pb.CreateCouponResponse{
		Id: id,
	}, nil
}

// DeleteCouponByID is a method of the couponService that deletes a coupon based on the provided ID.
// It checks if the user has admin privileges, and if so, it proceeds to delete the coupon.
func (s *couponService) DeleteCouponByID(ctx context.Context, req *pb.DeleteCouponByIDRequest) (*pb.DeleteCouponByIDResponse, error) {
	// Check if the user is an admin
	if _, err := validAdmin(ctx); err != nil {
		// If not an admin, return the permission error
		return nil, err
	}

	// Attempt to delete the coupon by ID using the coupon repository
	if err := s.couponRepo.DeleteByID(ctx, s.db, req.GetId()); err != nil {
		// If there is an error during deletion, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to delete coupon by id: %v", err.Error())
	}

	// Return an empty response indicating successful deletion
	return &pb.DeleteCouponByIDResponse{}, nil
}

// RetrieveCouponByCode is a method of the couponService that retrieves a coupon by its code.
// It also checks if the coupon can be used based on the provided conditions.
func (s *couponService) RetrieveCouponByCode(ctx context.Context, req *pb.RetrieveCouponByCodeRequest) (*pb.RetrieveCouponByCodeResponse, error) {
	// Retrieve the coupon by its code using the coupon repository
	coupon, err := s.couponRepo.RetrieveByCode(ctx, s.db, req.GetCode())
	if err != nil {
		// If there is an error during retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve coupon by code: %v", err.Error())
	}

	// If CheckUse is requested, perform additional checks on the usability of the coupon
	if req.GetCheckUse() {
		now := time.Now()

		// Check if the current time is within the valid range of the coupon's From and To times
		if coupon.From.Time.After(now) || coupon.To.Time.Before(now) {
			return &pb.RetrieveCouponByCodeResponse{
				CanUse: false,
			}, nil
		}

		// Perform additional checks based on the coupon type
		switch coupon.Type.String {
		case pb.CouponType_CouponType_LIMITED.String():
			// If the coupon type is LIMITED, check if it has reached its usage limit
			if coupon.Used.Int64 >= coupon.Total.Int64 {
				return &pb.RetrieveCouponByCodeResponse{
					CanUse: false,
				}, nil
			}
		case pb.CouponType_CouponType_USER.String():
			// If the coupon type is USER, retrieve user-specific coupon information
			userCoupon, err := s.userCouponRepo.RetrieveByCouponIDUserID(ctx, s.db, coupon.ID.Int64, req.UserId.Value)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "unable to retrieve user coupon: %v", err.Error())
			}
			// Check if the user-specific coupon has reached its usage limit
			if userCoupon.Used.Int64 >= userCoupon.Total.Int64 {
				return &pb.RetrieveCouponByCodeResponse{
					CanUse: false,
				}, nil
			}
		case pb.CouponType_CouponType_PRODUCT.String():
			// If the coupon type is PRODUCT, retrieve product-specific coupon information
			productCoupon, err := s.productCouponRepo.RetrieveByCouponID(ctx, s.db, coupon.ID.Int64)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "unable to retrieve product coupon: %v", err.Error())
			}
			// Check if the product-specific coupon has reached its usage limit
			if productCoupon.Used.Int64 >= productCoupon.Total.Int64 {
				return &pb.RetrieveCouponByCodeResponse{
					CanUse: false,
				}, nil
			}
		}
	}

	// Return the retrieved coupon information in the response
	return &pb.RetrieveCouponByCodeResponse{
		Type:         new(pb.CouponType).FromString(coupon.Type.String),
		Total:        coupon.Total.Int64,
		From:         timestamppb.New(coupon.From.Time),
		To:           timestamppb.New(coupon.To.Time),
		ImageUrl:     coupon.ImageURL.String,
		Description:  coupon.Description.String,
		DiscountType: new(pb.DiscountType).FromString(coupon.DiscountType.String),
		Value:        coupon.Value.Float64,
		Used:         coupon.Used.Int64,
		CanUse:       true,
	}, nil
}

// ListUsedCoupon is a method of the couponService that retrieves a list of coupons used by the current user.
func (s *couponService) ListUsedCoupon(ctx context.Context, _ *pb.ListUsedCouponRequest) (*pb.ListUsedCouponResponse, error) {
	// Extract user information from the context
	userCtx, err := xcontext.ExtractUserInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve user from context: %v", err.Error())
	}

	// Retrieve the list of used coupons by the current user from the usedCoupon repository
	coupons, err := s.usedCouponRepo.ListUsedCouponByUserID(ctx, s.db, userCtx.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve used coupon: %v", err.Error())
	}

	// Prepare the response data with details of used coupons
	respData := make([]*pb.ListUsedCouponResponse_Coupon, 0, len(coupons))
	for _, coupon := range coupons {
		respData = append(respData, &pb.ListUsedCouponResponse_Coupon{
			Code:        coupon.Coupon.Code.String,
			Description: coupon.Coupon.Description.String,
			ImageUrl:    coupon.Coupon.ImageURL.String,
			ApplyAt:     timestamppb.New(coupon.UsedCoupon.CreatedAt.Time),
		})
	}

	// Return the response with the list of used coupons
	return &pb.ListUsedCouponResponse{
		Data: respData,
	}, nil
}

// ApplyCoupon is a method of the couponService that applies a coupon for the current user.
// It checks if the user has the necessary permissions and retrieves the coupon by its code.
// If the coupon is found, it creates a record in the usedCoupon repository for the applied coupon.
func (s *couponService) ApplyCoupon(ctx context.Context, req *pb.ApplyCouponRequest) (*pb.ApplyCouponResponse, error) {
	// Extract user information from the context
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		// If user information cannot be extracted, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// Retrieve the coupon by its code from the coupon repository
	coupon, err := s.couponRepo.RetrieveByCode(ctx, s.db, req.GetCode())
	if err != nil {
		// If the coupon is not found, return a not found error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "coupon not found")
		}

		// If there is an internal error during coupon retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve coupon by code: %v", err.Error())
	}

	// Create a record in the usedCoupon repository for the applied coupon
	if err := s.usedCouponRepo.Create(ctx, s.db, &entity.UsedCoupon{
		CouponID:  coupon.ID,
		UserID:    pg_util.NullInt64(userCtx.UserID),
		CreatedAt: pg_util.NullTime(time.Now()),
	}); err != nil {
		// If there is an internal error during usedCoupon creation, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to create used coupon: %v", err.Error())
	}

	// Return an empty response indicating successful coupon application
	return &pb.ApplyCouponResponse{}, nil
}
