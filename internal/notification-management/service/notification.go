package service

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	pb "trintech/review/dto/msg/common"
	"trintech/review/pkg/email"
	"trintech/review/pkg/pubsub"
)

// notificationService is responsible for handling various notification-related tasks.
type notificationService struct {
	subscriber    pubsub.Subscriber
	emailProvider email.Provider
}

// SubscribeForgotPassword listens for messages related to forgot password events and processes them.
func (s *notificationService) SubscribeForgotPassword(ctx context.Context, _, value []byte) {
	// Unmarshal the received message into a ForgotPassword protobuf message.
	var user pb.ForgotPassword
	if err := proto.Unmarshal(value, &user); err != nil {
		slog.Error("unable to unmarshal forgot password data", "error", err)
		return
	}

	// Send a password reset email to the user.
	if err := s.emailProvider.SendMail(ctx, &email.EmailData{
		From: "trintech@gmail.com",
		To:   user.GetEmail(),
		Content: fmt.Sprintf(`
		Hi %s,
		Please click the link http://trintech.com/%s/%s to reset your password
		`,
			user.GetName(),
			user.GetEmail(),
			user.GetResetToken(),
		),
	}); err != nil {
		slog.Error("unable to send email", "error", err)
	}
}
