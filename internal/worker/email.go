package worker

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/events"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/mail"
)

func (w *Worker) handleEmail(
	ctx context.Context,
	msg events.EmailMessage,
) error {

	switch msg.Type {

	case "reset_password":
		return w.sendResetPassword(ctx, msg)

	default:
		return nil
	}
}

func (w *Worker) sendResetPassword(
	ctx context.Context,
	msg events.EmailMessage,
) error {

	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: msg.To},
		},
		Subject: msg.Subject,
		Text:    msg.Body,
	}

	return w.mailService.SendMail(ctx, mailContent)
}
