// user-management-api/pkg/mail/mail_factory.go
package mail

import (
	"fmt"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type ProviderFactory interface {
	CreateProvider(config *MailConfig) (EmailProviderService, error)
}

type MailtrapProviderFactory struct{}

func (f *MailtrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error) {
	return NewMailtrapProvider(config)
}

func NewProvider(config *MailConfig) (EmailProviderService, error) {
	switch config.ProviderType {
	case ProviderMailtrap:
		return NewMailtrapProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider")
	}
}
