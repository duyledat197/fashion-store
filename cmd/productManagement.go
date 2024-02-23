/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	couponpb "trintech/review/dto/coupon-management/coupon"
	pb "trintech/review/dto/product-management/product"
	"trintech/review/internal/product-management/service"
	"trintech/review/pkg/grpc_client"
	"trintech/review/pkg/grpc_server"
	"trintech/review/pkg/postgres_client"
)

// productManagementCmd represents the productManagement command
var productManagementCmd = &cobra.Command{
	Use:   "productManagement",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadDefault()
		loadProductManagement(ctx)
		errChan := make(chan error)
		start(ctx, errChan)
		err := <-errChan
		if err != nil {
			slog.Error(err.Error())
			stop(ctx)
		}
	},
}

func init() {
	rootCmd.AddCommand(productManagementCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// productManagementCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// productManagementCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// loadProductManagement initializes and loads the Product Management service.
func loadProductManagement(_ context.Context) {
	// Create a new PostgreSQL client using the specified address.
	pgClient := postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())

	// Create a gRPC client connection to the Coupon service.
	couponClientConn := grpc_client.NewGrpcClient(cfgs.CouponService)

	// Create a gRPC client instance for the Coupon service.
	couponClient := couponpb.NewCouponServiceClient(couponClientConn)

	// Log the address of the Coupon service.
	log.Println(cfgs.CouponService.Address())

	// Create a new ProductService instance with the PostgreSQL client and Coupon client.
	service := service.NewProductService(pgClient, couponClient)

	// Create a new gRPC server using the specified configuration.
	srv := grpc_server.NewGrpcServer(cfgs.ProductService)

	// Register the ProductService implementation with the gRPC server.
	pb.RegisterProductServiceServer(srv.Server, service)

	// Append the PostgreSQL client and Coupon gRPC client connection to the list of factories.
	factories = append(factories, pgClient, couponClientConn)

	// Append the gRPC server to the list of processors.
	processors = append(processors, srv)
}
