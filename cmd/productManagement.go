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

func loadProductManagement(_ context.Context) {
	pgClient := postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())
	couponClientConn := grpc_client.NewGrpcClient(cfgs.CouponService)
	couponClient := couponpb.NewCouponServiceClient(couponClientConn)

	log.Println(cfgs.CouponService.Address())
	service := service.NewProductService(pgClient, couponClient)

	srv := grpc_server.NewGrpcServer(cfgs.ProductService)

	pb.RegisterProductServiceServer(srv.Server, service)

	factories = append(factories, pgClient, couponClientConn)
	processors = append(processors, srv)
}
