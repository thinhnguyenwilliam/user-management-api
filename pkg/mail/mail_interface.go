// user-management-api/pkg/mail/mail_interface.go
package mail

import "context"

type EmailProviderService interface {
	SendMail(ctx context.Context, email *Email) error
}
