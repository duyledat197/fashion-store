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

type notificationService struct {
	subscriber pubsub.Subscriber

	emailProvider email.Provider
}

// SubscribeForgotPassword ...
func (s *notificationService) SubscribeForgotPassword(ctx context.Context, _, value []byte) {
	var user pb.ForgotPassword
	if err := proto.Unmarshal(value, &user); err != nil {
		slog.Error("unable to unmarshal forgot password data", "error", err)
		return
	}

	if err := s.emailProvider.SendMail(ctx, &email.EmailData{
		From: "trintech@gmail.com",
		To:   user.GetEmail(),
		Content: fmt.Sprintf(`
		Hi %s,
		Please click to link http://trintech.com/%s/%s to reset password
		`,
			user.GetName(),
			user.GetEmail(),
			user.GetResetToken(),
		),
	}); err != nil {
		slog.Error("unable to send email", "error", err)
	}
}
