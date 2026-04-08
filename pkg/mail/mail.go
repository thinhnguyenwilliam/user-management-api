// user-management-api/pkg/mail/mail.go
package mail

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/logger"
)

type Email struct {
	From     Address   `json:"from"`
	To       []Address `json:"to"`
	Subject  string    `json:"subject"`
	Text     string    `json:"text"`
	Category string    `json:"category"`
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailConfig struct {
	ProviderType ProviderType
	Mailtrap     config.MailTrapConfig
	MaxRetries   int
	Timeout      time.Duration
	Logger       *zerolog.Logger
}

type MailService struct {
	config   *MailConfig
	provider EmailProviderService
	logger   *zerolog.Logger
}

func NewMailService(cfg *config.Config, logger *zerolog.Logger) (EmailProviderService, error) {
	mailCfg := &MailConfig{
		ProviderType: ProviderType(cfg.MailTrap.MailProviderType),
		Mailtrap:     cfg.MailTrap,
		MaxRetries:   3,
		Timeout:      10 * time.Second,
		Logger:       logger,
	}

	provider, err := NewProvider(mailCfg)
	if err != nil {
		return nil, err
	}

	return &MailService{
		config:   mailCfg,
		provider: provider,
		logger:   logger,
	}, nil
}

func (ms *MailService) SendMail(ctx context.Context, email *Email) error {
	traceID := logger.GetTraceID(ctx)

	ctx, cancel := context.WithTimeout(ctx, ms.config.Timeout)
	defer cancel()

	start := time.Now()
	var lastErr error

	for attempt := 1; attempt <= ms.config.MaxRetries; attempt++ {
		ctxAttempt, cancel := context.WithTimeout(ctx, ms.config.Timeout)

		startAttempt := time.Now()
		err := ms.provider.SendMail(ctxAttempt, email)
		cancel()

		if err == nil {
			ms.logger.Info().
				Str("trace_id", traceID).
				Dur("duration", time.Since(startAttempt)).
				Int("attempt", attempt).
				Msg("Email sent successfully")
			return nil
		}

		lastErr = err

		ms.logger.Warn().
			Str("trace_id", traceID).
			Int("attempt", attempt).
			Err(err).
			Msg("Retry sending email")

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	ms.logger.Error().
		Str("trace_id", traceID).
		Dur("total_duration", time.Since(start)).
		Err(lastErr).
		Msg("Failed after retries")

	return utils.WrapError("send mail failed", utils.ErrCodeInternal, lastErr)
}
