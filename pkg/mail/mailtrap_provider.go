// user-management-api/pkg/mail/mailtrap_provider.go
package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/logger"
)

type MailtrapProvider struct {
	client *http.Client
	config *config.MailTrapConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(cfg *MailConfig) (EmailProviderService, error) {
	return &MailtrapProvider{
		client: &http.Client{Timeout: cfg.Timeout},
		config: &cfg.Mailtrap,
		logger: cfg.Logger,
	}, nil
}

func (p *MailtrapProvider) SendMail(ctx context.Context, email *Email) error {
	p.logger.Info().
		Str("url", p.config.MailtrapURL).
		Str("api_key_prefix", p.config.MailtrapAPIKey[:5]).
		Msg("mailtrap config")

	traceID := logger.GetTraceID(ctx)
	start := time.Now()

	email.From = Address{
		Email: p.config.MailSender,
		Name:  p.config.NameSender,
	}

	payload, err := json.Marshal(email)
	if err != nil {
		return utils.WrapError("marshal email failed", utils.ErrCodeInternal, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.MailtrapURL, bytes.NewReader(payload))
	if err != nil {
		return utils.WrapError("create request failed", utils.ErrCodeInternal, err)
	}

	req.Header.Add("Authorization", "Bearer "+p.config.MailtrapAPIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error().
			Str("trace_id", traceID).
			Dur("duration", time.Since(start)).
			Err(err).
			Msg("send request failed")

		return utils.WrapError("send request failed", utils.ErrCodeInternal, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)

		p.logger.Error().
			Str("trace_id", traceID).
			Int("status_code", resp.StatusCode).
			Str("response_body", string(body)).
			Msg("mailtrap error")

		return utils.NewError(
			fmt.Sprintf("mailtrap error: %d - %s", resp.StatusCode, string(body)),
			utils.ErrCodeInternal,
		)
	}

	return nil
}
