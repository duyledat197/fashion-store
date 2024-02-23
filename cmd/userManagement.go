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

// loadUserManagement initializes and loads the User Management service.
func loadUserManagement(_ context.Context) {
	// Create a new PostgreSQL client using the specified address.
	pgClient := postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())

	// Create a new AuthService instance with the PostgreSQL client, mock publisher, and token generator.
	service := service.NewAuthService(pgClient, &mocks.Publisher{}, tokenGenerator)

	// Create a new gRPC server using the specified configuration.
	srv := grpc_server.NewGrpcServer(cfgs.UserService)

	// Register the AuthService implementation with the gRPC server.
	pb.RegisterAuthServiceServer(srv.Server, service)

	// Append the PostgreSQL client to the list of factories.
	factories = append(factories, pgClient)

	// Append the gRPC server to the list of processors.
	processors = append(processors, srv)
}
