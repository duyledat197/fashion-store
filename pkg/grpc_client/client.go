package grpc_client

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"trintech/review/config"
)

type GrpcClient struct {
	*grpc.ClientConn
	cfg *config.Endpoint
}

func NewGrpcClient(cfg *config.Endpoint) *GrpcClient {
	return &GrpcClient{
		cfg: cfg,
	}
}

func (c *GrpcClient) Connect(ctx context.Context) error {
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
	}

	conn, err := grpc.Dial(
		c.cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_retry.StreamClientInterceptor(optsRetry...),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(optsRetry...),
		)),
	)

	if err != nil {
		return err
	}

	c.ClientConn = conn

	return nil
}

func (c *GrpcClient) Close(ctx context.Context) error {
	return c.Close(ctx)
}
