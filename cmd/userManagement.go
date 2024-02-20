/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	pb "trintech/review/dto/user-management/auth"
	"trintech/review/internal/user-management/service"
	"trintech/review/mocks"
	"trintech/review/pkg/grpc_server"
	"trintech/review/pkg/postgres_client"
)

// userManagementCmd represents the userManagement command
var userManagementCmd = &cobra.Command{
	Use:   "userManagement",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadDefault()
		loadUserManagement(ctx)
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
	rootCmd.AddCommand(userManagementCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userManagementCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userManagementCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadUserManagement(ctx context.Context) {
	pgClient := postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())

	service := service.NewAuthService(pgClient, &mocks.Publisher{}, tokenGenerator)

	srv := grpc_server.NewGrpcServer(cfgs.UserService)

	pb.RegisterAuthServiceServer(srv.Server, service)

	factories = append(factories, pgClient)
	processors = append(processors, srv)
}
