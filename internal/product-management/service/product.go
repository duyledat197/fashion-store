// Package service ...
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	couponpb "trintech/review/dto/coupon-management/coupon"
	pb "trintech/review/dto/product-management/product"
	"trintech/review/internal/product-management/entity"
	"trintech/review/internal/product-management/repository/postgres"
	userEntity "trintech/review/internal/user-management/entity"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
)

// productService is representation of
type productService struct {
	productRepo interface {
		List(ctx context.Context, db database.Executor, offset, limit int64) ([]*entity.Product, error)
		Count(ctx context.Context, db database.Executor) (int64, error)
		RetrieveByID(ctx context.Context, db database.Executor, id int64) (*entity.Product, error)
		Create(ctx context.Context, db database.Executor, data *entity.Product) (int64, error)
		UpdateByID(ctx context.Context, db database.Executor, id int64, data *entity.Product) error
		DeleteByID(ctx context.Context, db database.Executor, id int64) error
		DeleteByIDs(ctx context.Context, db database.Executor, ids []int64) error
	}

	purchasedProductRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.PurchasedProduct) error
	}

	pb.UnimplementedProductServiceServer

	db database.Database

	couponServiceClient couponpb.CouponServiceClient
}

// NewProductService ...
func NewProductService(
	db database.Database,
	couponServiceClient couponpb.CouponServiceClient,
) pb.ProductServiceServer {
	return &productService{
		db:                   db,
		couponServiceClient:  couponServiceClient,
		productRepo:          postgres.NewProductRepository(),
		purchasedProductRepo: postgres.NewPurchasedProductRepository(),
	}
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

// CreateProduct ...
func (s *productService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	userCtx, err := validAdmin(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.productRepo.Create(ctx, s.db, &entity.Product{
		Name:        pg_util.NullString(req.GetName()),
		Type:        pg_util.NullString(req.GetType()),
		Description: pg_util.NullString(req.GetDescription()),
		ImageURLs:   pg_util.StringArray(req.GetImageUrls()),
		CreatedBy:   pg_util.NullInt64(userCtx.UserID),
		Price:       pg_util.NullFloat64(req.GetPrice()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create product: %v", err.Error())
	}

	return &pb.CreateProductResponse{
		Id: id,
	}, nil
}

// DeleteProductByID ...
func (s *productService) DeleteProductByID(ctx context.Context, req *pb.DeleteProductByIDRequest) (*pb.DeleteProductByIDResponse, error) {
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	if err := s.productRepo.DeleteByID(ctx, s.db, req.GetId()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}

		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	return &pb.DeleteProductByIDResponse{}, nil
}

// DeleteProductByIDs ...
func (s *productService) DeleteProductByIDs(ctx context.Context, req *pb.DeleteProductByIDsRequest) (*pb.DeleteProductByIDsResponse, error) {
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	if err := s.productRepo.DeleteByIDs(ctx, s.db, req.GetIds()); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	return &pb.DeleteProductByIDsResponse{}, nil
}

// ListProduct ...
func (s *productService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	list, err := s.productRepo.List(ctx, s.db, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve list product: %v", err.Error())
	}

	respData := make([]*pb.Product, 0, len(list))

	for _, product := range list {
		respData = append(respData, &pb.Product{
			Name:        product.Name.String,
			Type:        product.Type.String,
			ImageUrls:   pg_util.StringArrayValue(product.ImageURLs),
			Description: product.Description.String,
			Price:       product.Price.Float64,
		})
	}

	total, err := s.productRepo.Count(ctx, s.db)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to count product: %v", err.Error())
	}

	return &pb.ListProductResponse{
		Data:  respData,
		Total: total,
	}, nil
}

// RetrieveProductByID ...
func (s *productService) RetrieveProductByID(ctx context.Context, req *pb.RetrieveProductByIDRequest) (*pb.RetrieveProductByIDResponse, error) {
	product, err := s.productRepo.RetrieveByID(ctx, s.db, req.GetId())
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, status.Errorf(codes.NotFound, "product not found")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "unable to retrieve product: %v", err.Error())
	}

	return &pb.RetrieveProductByIDResponse{
		Data: &pb.Product{
			Name:        product.Name.String,
			Type:        product.Type.String,
			ImageUrls:   pg_util.StringArrayValue(product.ImageURLs),
			Description: product.Description.String,
		},
	}, nil
}

// UpdateProductByID ...
func (s *productService) UpdateProductByID(ctx context.Context, req *pb.UpdateProductByIDRequest) (*pb.UpdateProductByIDResponse, error) {
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	if err := s.productRepo.UpdateByID(ctx, s.db, req.GetId(), &entity.Product{
		Name:        pg_util.NullString(req.GetName()),
		Type:        pg_util.NullString(req.GetType()),
		Description: pg_util.NullString(req.GetDescription()),
		ImageURLs:   pg_util.StringArray(req.GetImageUrls()),
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}

		return nil, status.Errorf(codes.Internal, "unable to update product: %v", err.Error())
	}

	return &pb.UpdateProductByIDResponse{}, nil
}

func (s *productService) PurchaseProduct(ctx context.Context, req *pb.PurchaseProductRequest) (*pb.PurchaseProductResponse, error) {
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	product, err := s.productRepo.RetrieveByID(ctx, s.db, req.GetId())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.NotFound, "product not found")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "unable to retrieve product: %v", err.Error())
	}

	var (
		discountTotal float64 = product.Price.Float64
	)
	if req.GetCoupon() != nil {
		coupon, err := s.couponServiceClient.RetrieveCouponByCode(http_server.InjectIncomingCtxToOutgoingCtx(ctx), &couponpb.RetrieveCouponByCodeRequest{
			Code:     req.GetCoupon().GetValue(),
			CheckUse: true,
			UserId:   wrapperspb.Int64(userCtx.UserID),
		})
		if err != nil {
			stt, _ := status.FromError(err)
			switch stt.Code() {
			case codes.NotFound:
				return nil, status.Errorf(codes.NotFound, "coupon not found")
			case codes.FailedPrecondition:
				return nil, status.Errorf(codes.FailedPrecondition, "unable to apply this coupon")
			default:
				return nil, status.Errorf(codes.Internal, "unable to apply this coupon: %v", err.Error())
			}
		}

		if !coupon.CanUse {
			return nil, status.Errorf(codes.FailedPrecondition, "user cannot apply this coupon")
		}

		switch coupon.DiscountType {
		case couponpb.DiscountType_DiscountType_PERCENT:
			discountTotal = product.Price.Float64 * coupon.Value / 100
		case couponpb.DiscountType_DiscountType_VALUE:
			discountTotal = coupon.GetValue()
		}
	}

	if err := database.Transaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.purchasedProductRepo.Create(ctx, tx, &entity.PurchasedProduct{
			ProductID: product.ID,
			UserID:    pg_util.NullInt64(userCtx.UserID),
			Price:     product.Price,
			Discount:  pg_util.NullFloat64(discountTotal),
			Total:     pg_util.NullFloat64(max(0, discountTotal)),
			Coupon:    pg_util.NullString(req.GetCoupon().Value),
		}); err != nil {
			return fmt.Errorf("unable to create purchase product: %v", err)
		}

		if _, err := s.couponServiceClient.ApplyCoupon(http_server.InjectIncomingCtxToOutgoingCtx(ctx), &couponpb.ApplyCouponRequest{
			Code: req.Coupon.Value,
		}); err != nil {
			stt, _ := status.FromError(err)
			switch stt.Code() {
			case codes.NotFound:
				return status.Errorf(codes.NotFound, "coupon not found")
			case codes.FailedPrecondition:
				return status.Errorf(codes.FailedPrecondition, "unable to apply this coupon")
			default:
				return status.Errorf(codes.Internal, "unable to apply this coupon: %v", err.Error())
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.PurchaseProductResponse{}, nil
}
