package worker

import (
	"context"
	"encoding/json"

	"github.com/thinhnguyenwilliam/user-management-api/internal/events"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/mail"
)

func (w *Worker) handleMessage(
	ctx context.Context,
	msg events.EmailMessage,
) error {

	err := w.handleEmail(ctx, msg)
	if err == nil {
		return nil
	}

	// 🔥 retry logic
	msg.Retry++

	body, _ := json.Marshal(msg)

	if msg.Retry >= msg.MaxRetry {
		// 👉 send DLQ
		return w.mq.Publish(ctx, "send_email_dlq", body)
	}

	// 👉 retry queue
	return w.mq.Publish(ctx, "send_email_retry", body)
}

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
