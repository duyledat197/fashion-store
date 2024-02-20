/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"

	couponpb "trintech/review/dto/coupon-management/coupon"
	productpb "trintech/review/dto/product-management/product"
	userpb "trintech/review/dto/user-management/auth"
	"trintech/review/pkg/grpc_client"
	"trintech/review/pkg/http_server"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadDefault()
		loadGateway(ctx)
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
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadGateway(ctx context.Context) {
	userClientConn := grpc_client.NewGrpcClient(cfgs.UserService)
	productClientConn := grpc_client.NewGrpcClient(cfgs.ProductService)
	couponClientConn := grpc_client.NewGrpcClient(cfgs.CouponService)

	userClient := userpb.NewAuthServiceClient(userClientConn)
	productClient := productpb.NewProductServiceClient(productClientConn)
	couponClient := couponpb.NewCouponServiceClient(couponClientConn)

	httpServer := http_server.NewHttpServer(
		func(mux *runtime.ServeMux) {
			userpb.RegisterAuthServiceHandlerClient(ctx, mux, userClient)
			productpb.RegisterProductServiceHandlerClient(ctx, mux, productClient)
			couponpb.RegisterCouponServiceHandlerClient(ctx, mux, couponClient)
		},
		cfgs.GatewayService,
		tokenGenerator,
	)

	factories = append(factories, userClientConn, productClientConn, couponClientConn)
	processors = append(processors, httpServer)
}
