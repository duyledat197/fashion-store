// Package service ...
package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trintech/review/dto/product-management/product"
	"trintech/review/internal/product-management/entity"
	postgresclient "trintech/review/pkg/postgres_client"
)

// ProductService is representation of
type ProductService struct {
	productRepo interface {
		List(ctx context.Context, offset, limit int64) ([]*entity.Product, error)
		RetrieveByID(ctx context.Context, id int64) (*entity.Product, error)
		Create(ctx context.Context, data *entity.Product) (int64, error)
		UpdateByID(ctx context.Context, id int64, data *entity.Product) error
		DeleteByID(ctx context.Context, id int64) error
		DeleteByIDs(ctx context.Context, ids []int64) error
	}

	pb.UnimplementedProductServiceServer
}

// CreateProduct ...
func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	id, err := s.productRepo.Create(ctx, &entity.Product{
		Name:        postgresclient.PgTypeText(req.GetName()),
		Type:        postgresclient.PgTypeText(req.GetType()),
		Description: postgresclient.PgTypeText(req.GetDescription()),
		ImageURL:    postgresclient.PgTypeText(req.GetImageUrl()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create product: %v", err.Error())
	}

	return &pb.CreateProductResponse{
		Id: id,
	}, nil
}

// DeleteProductByID ...
func (s *ProductService) DeleteProductByID(ctx context.Context, req *pb.DeleteProductByIDRequest) (*pb.DeleteProductByIDResponse, error) {
	if err := s.productRepo.DeleteByID(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	return &pb.DeleteProductByIDResponse{}, nil
}

// DeleteProductByIDs ...
func (s *ProductService) DeleteProductByIDs(ctx context.Context, req *pb.DeleteProductByIDsRequest) (*pb.DeleteProductByIDsResponse, error) {
	if err := s.productRepo.DeleteByIDs(ctx, req.GetIds()); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to delete product: %v", err.Error())
	}

	return &pb.DeleteProductByIDsResponse{}, nil
}

// ListProduct ...
func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	list, err := s.productRepo.List(ctx, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve list product: %v", err.Error())
	}

	respData := make([]*pb.Product, 0, len(list))

	for _, product := range list {
		respData = append(respData, &pb.Product{
			Name:        product.Name.String,
			Type:        product.Type.String,
			ImageUrl:    product.ImageURL.String,
			Description: product.Description.String,
		})
	}

	return &pb.ListProductResponse{
		Data: respData,
	}, nil
}

// RetrieveProductByID ...
func (s *ProductService) RetrieveProductByID(ctx context.Context, req *pb.RetrieveProductByIDRequest) (*pb.RetrieveProductByIDResponse, error) {
	product, err := s.productRepo.RetrieveByID(ctx, req.GetId())
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
			ImageUrl:    product.ImageURL.String,
			Description: product.Description.String,
		},
	}, nil
}

// UpdateProductByID ...
func (s *ProductService) UpdateProductByID(ctx context.Context, req *pb.UpdateProductByIDRequest) (*pb.UpdateProductByIDResponse, error) {
	err := s.productRepo.UpdateByID(ctx, req.GetId(), &entity.Product{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update product: %v", err.Error())
	}

	return &pb.UpdateProductByIDResponse{}, nil
}
