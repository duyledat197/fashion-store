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

// validAdmin checks admin validation from context.
// It extracts user information from the context and validates if the user has admin or super admin role.
func validAdmin(ctx context.Context) (*xcontext.UserInfo, error) {
	// Extract user information from the context
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		// If user information is not found, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// Check if the user has admin or super admin role
	if !slices.Contains([]string{userEntity.UserRole_Admin, userEntity.UserRole_SuperAdmin}, userCtx.Role) {
		// If not, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// Return user information if validation is successful
	return userCtx, nil
}

// CreateProduct is a method of the productService that handles the creation of a new product.
// It validates the admin user, creates a product in the repository, and returns the created product's ID.
func (s *productService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// Validate admin user
	userCtx, err := validAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new product in the repository
	id, err := s.productRepo.Create(ctx, s.db, &entity.Product{
		Name:        pg_util.NullString(req.GetName()),
		Type:        pg_util.NullString(req.GetType()),
		Description: pg_util.NullString(req.GetDescription()),
		ImageURLs:   pg_util.StringArray(req.GetImageUrls()),
		CreatedBy:   pg_util.NullInt64(userCtx.UserID),
		Price:       pg_util.NullFloat64(req.GetPrice()),
	})

	if err != nil {
		// If there is an error during product creation, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to create product: %v", err.Error())
	}

	// Return the created product's ID
	return &pb.CreateProductResponse{
		Id: id,
	}, nil
}

// DeleteProductByID is a method of the productService that handles the deletion of a product by ID.
// It validates the admin user, deletes the product in the repository, and returns an empty response.
func (s *productService) DeleteProductByID(ctx context.Context, req *pb.DeleteProductByIDRequest) (*pb.DeleteProductByIDResponse, error) {
	// Validate admin user
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	// Delete the product in the repository by ID
	if err := s.productRepo.DeleteByID(ctx, s.db, req.GetId()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the product is not found, return a not found error
			return nil, status.Errorf(codes.NotFound, "product not found")
		}

		// If there is an error during product deletion, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	// Return an empty response indicating successful deletion
	return &pb.DeleteProductByIDResponse{}, nil
}

// DeleteProductByIDs is a method of the productService that handles the deletion of products by IDs.
// It validates the admin user, deletes the products in the repository by IDs, and returns an empty response.
func (s *productService) DeleteProductByIDs(ctx context.Context, req *pb.DeleteProductByIDsRequest) (*pb.DeleteProductByIDsResponse, error) {
	// Validate admin user
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	// Delete the products in the repository by IDs
	if err := s.productRepo.DeleteByIDs(ctx, s.db, req.GetIds()); err != nil {
		// If there is an error during product deletion, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	// Return an empty response indicating successful deletion
	return &pb.DeleteProductByIDsResponse{}, nil
}

// ListProduct is a method of the productService that retrieves a list of products.
// It validates the admin user, retrieves the list of products from the repository, and returns the response.
func (s *productService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	// Retrieve the list of products from the repository
	list, err := s.productRepo.List(ctx, s.db, req.GetOffset(), req.GetLimit())
	if err != nil {
		// If there is an error during product retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve list product: %v", err.Error())
	}

	// Transform the list of products to the response format
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

	// Get the total count of products
	total, err := s.productRepo.Count(ctx, s.db)
	if err != nil {
		// If there is an error during count retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to count product: %v", err.Error())
	}

	// Return the list of products and total count in the response
	return &pb.ListProductResponse{
		Data:  respData,
		Total: total,
	}, nil
}

// RetrieveProductByID is a method of the productService that retrieves a product by ID.
// It validates the admin user, retrieves the product from the repository by ID, and returns the response.
func (s *productService) RetrieveProductByID(ctx context.Context, req *pb.RetrieveProductByIDRequest) (*pb.RetrieveProductByIDResponse, error) {
	// Retrieve the product from the repository by ID
	product, err := s.productRepo.RetrieveByID(ctx, s.db, req.GetId())
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// If the product is not found, return a not found error
		return nil, status.Errorf(codes.NotFound, "product not found")
	case err != nil:
		// If there is an error during product retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve product: %v", err.Error())
	}

	// Transform the retrieved product to the response format
	return &pb.RetrieveProductByIDResponse{
		Data: &pb.Product{
			Name:        product.Name.String,
			Type:        product.Type.String,
			ImageUrls:   pg_util.StringArrayValue(product.ImageURLs),
			Description: product.Description.String,
			Price:       product.Price.Float64,
		},
	}, nil
}

// UpdateProductByID is a method of the productService that updates a product by ID.
// It validates the admin user, updates the product in the repository by ID, and returns an empty response.
func (s *productService) UpdateProductByID(ctx context.Context, req *pb.UpdateProductByIDRequest) (*pb.UpdateProductByIDResponse, error) {
	// Validate admin user
	if _, err := validAdmin(ctx); err != nil {
		return nil, err
	}

	// Update the product in the repository by ID
	if err := s.productRepo.UpdateByID(ctx, s.db, req.GetId(), &entity.Product{
		Name:        pg_util.NullString(req.GetName()),
		Type:        pg_util.NullString(req.GetType()),
		Description: pg_util.NullString(req.GetDescription()),
		ImageURLs:   pg_util.StringArray(req.GetImageUrls()),
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the product is not found, return a not found error
			return nil, status.Errorf(codes.NotFound, "product not found")
		}

		// If there is an error during product update, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to update product: %v", err.Error())
	}

	// Return an empty response indicating successful update
	return &pb.UpdateProductByIDResponse{}, nil
}

// PurchaseProduct is a method of the productService that handles the purchase of a product.
// It extracts user information, retrieves the product, applies a coupon if provided,
// creates a purchase record in the repository, and returns an empty response.
func (s *productService) PurchaseProduct(ctx context.Context, req *pb.PurchaseProductRequest) (*pb.PurchaseProductResponse, error) {
	// Extract user information from the context
	userCtx, ok := http_server.ExtractUserInfoFromCtx(ctx)
	if !ok {
		// If user information is not found, return a permission denied error
		return nil, status.Errorf(codes.PermissionDenied, "user doesn't have permission")
	}

	// Retrieve the product by ID
	product, err := s.productRepo.RetrieveByID(ctx, s.db, req.GetId())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// If the product is not found, return a not found error
		return nil, status.Errorf(codes.NotFound, "product not found")
	case err != nil:
		// If there is an error during product retrieval, return an internal server error
		return nil, status.Errorf(codes.Internal, "unable to retrieve product: %v", err.Error())
	}

	// Initialize discountTotal with the original product price
	var discountTotal float64 = product.Price.Float64

	// Apply coupon if provided
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
				// If the coupon is not found, return a not found error
				return nil, status.Errorf(codes.NotFound, "coupon not found")
			case codes.FailedPrecondition:
				// If there is a failed precondition, return a failed precondition error
				return nil, status.Errorf(codes.FailedPrecondition, "unable to apply this coupon")
			default:
				// If there is an internal error, return an internal server error
				return nil, status.Errorf(codes.Internal, "unable to apply this coupon: %v", err.Error())
			}
		}

		// Check if the user can use this coupon
		if !coupon.CanUse {
			return nil, status.Errorf(codes.FailedPrecondition, "user cannot apply this coupon")
		}

		// Apply discount based on the coupon type
		switch coupon.DiscountType {
		case couponpb.DiscountType_DiscountType_PERCENT:
			discountTotal = product.Price.Float64 * coupon.Value / 100
		case couponpb.DiscountType_DiscountType_VALUE:
			discountTotal = coupon.GetValue()
		}
	}

	// Perform the purchase operation in a database transaction
	if err := database.Transaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// Create a purchased product record
		purchaseProduct := &entity.PurchasedProduct{
			ProductID: product.ID,
			UserID:    pg_util.NullInt64(userCtx.UserID),
			Price:     product.Price,
			Discount:  pg_util.NullFloat64(discountTotal),
			Total:     pg_util.NullFloat64(max(0, discountTotal)),
		}
		if req.GetCoupon() != nil {
			purchaseProduct.Coupon = pg_util.NullString(req.GetCoupon().Value)
		}

		// Create the purchased product record in the repository
		if err := s.purchasedProductRepo.Create(ctx, tx, purchaseProduct); err != nil {
			return fmt.Errorf("unable to create purchase product: %v", err)
		}

		// Apply the coupon if provided
		if req.GetCoupon() != nil {
			if _, err := s.couponServiceClient.ApplyCoupon(http_server.InjectIncomingCtxToOutgoingCtx(ctx), &couponpb.ApplyCouponRequest{
				Code: req.Coupon.Value,
			}); err != nil {
				stt, _ := status.FromError(err)
				switch stt.Code() {
				case codes.NotFound:
					// If the coupon is not found, return a not found error
					return status.Errorf(codes.NotFound, "coupon not found")
				case codes.FailedPrecondition:
					// If there is a failed precondition, return a failed precondition error
					return status.Errorf(codes.FailedPrecondition, "unable to apply this coupon")
				default:
					// If there is an internal error, return an internal server error
					return status.Errorf(codes.Internal, "unable to apply this coupon: %v", err.Error())
				}
			}
		}

		return nil
	}); err != nil {
		// If there is an error during the transaction, return the error
		return nil, err
	}

	// Return an empty response indicating successful purchase
	return &pb.PurchaseProductResponse{}, nil
}
