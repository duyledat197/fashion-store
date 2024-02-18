package service

import (
	"context"

	pb "trintech/review/dto/coupon-management/coupon"
)

type couponService struct {
	couponRepo interface {
	}

	pb.UnimplementedCouponServiceServer
}

func (s *couponService) CreateCoupon(context.Context, *pb.CreateCouponRequest) (*pb.CreateCouponResponse, error) {

}
func (s *couponService) DeleteCouponByID(context.Context, *pb.DeleteCouponByIDRequest) (*pb.DeleteCouponByIDResponse, error) {
}
