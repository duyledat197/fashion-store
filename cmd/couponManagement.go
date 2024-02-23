/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	pb "trintech/review/dto/coupon-management/coupon"
	"trintech/review/internal/coupon-management/service"
	"trintech/review/pkg/grpc_server"
	"trintech/review/pkg/postgres_client"
)

// couponManagementCmd represents the couponManagement command
var couponManagementCmd = &cobra.Command{
	Use:   "couponManagement",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadDefault()
		loadCouponManagement(ctx)
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
	rootCmd.AddCommand(couponManagementCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// couponManagementCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// couponManagementCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// loadCouponManagement initializes and loads the Coupon Management service.
func loadCouponManagement(_ context.Context) {
	// Create a new PostgreSQL client using the specified address.
	pgClient := postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())

	// Create a new CouponService instance with the PostgreSQL client.
	service := service.NewCouponService(pgClient)

	// Create a new gRPC server using the specified configuration.
	srv := grpc_server.NewGrpcServer(cfgs.CouponService)

	// Register the CouponService implementation with the gRPC server.
	pb.RegisterCouponServiceServer(srv.Server, service)

	// Append the all factory client to the list of factories.
	factories = append(factories, pgClient)

	// Append the all server to the list of processors.
	processors = append(processors, srv)
}
