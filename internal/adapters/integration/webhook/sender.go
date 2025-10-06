package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/ports/output"
)

type HTTPWebhookSender struct {
	httpClient *http.Client
	logger     *logger.Logger
}

func NewHTTPWebhookSender(httpClient *http.Client, logger *logger.Logger) output.WebhookSender {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &HTTPWebhookSender{
		httpClient: httpClient,
		logger:     logger,
	}
}
func (s *HTTPWebhookSender) SendWebhook(ctx context.Context, url string, secret *string, event *output.WebhookEvent) error {
	if url == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	if event == nil {
		return fmt.Errorf("webhook event cannot be nil")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "zpwoot-webhook/1.0")
	req.Header.Set("X-Webhook-Event", event.Type)
	req.Header.Set("X-Session-ID", event.SessionID)
	req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(event.Timestamp.Unix(), 10))

	if secret != nil && *secret != "" {
		signature := GenerateSignature(payload, *secret)
		req.Header.Set("X-Webhook-Signature", signature)
	}

	return s.sendWithRetry(ctx, req, url, event.Type)
}
func (s *HTTPWebhookSender) sendWithRetry(ctx context.Context, req *http.Request, url, eventType string) error {
	maxRetries := 3
	retryDelays := []time.Duration{
		0,
		5 * time.Second,
		15 * time.Second,
	}

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelays[attempt]):
			}
		}

		s.logger.Debug().
			Str("url", url).
			Str("event_type", eventType).
			Int("attempt", attempt+1).
			Int("max_retries", maxRetries).
			Msg("Sending webhook")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed (attempt %d/%d): %w", attempt+1, maxRetries, err)
			s.logger.Warn().
				Err(err).
				Str("url", url).
				Str("event_type", eventType).
				Int("attempt", attempt+1).
				Msg("Webhook request failed")

			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			s.logger.Info().
				Str("url", url).
				Str("event_type", eventType).
				Int("status_code", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("Webhook sent successfully")

			return nil
		}

		resp.Body.Close()
		lastErr = fmt.Errorf("webhook returned status %d (attempt %d/%d)", resp.StatusCode, attempt+1, maxRetries)
		s.logger.Warn().
			Str("url", url).
			Str("event_type", eventType).
			Int("status_code", resp.StatusCode).
			Int("attempt", attempt+1).
			Msg("Webhook returned error status")

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			s.logger.Error().
				Str("url", url).
				Str("event_type", eventType).
				Int("status_code", resp.StatusCode).
				Msg("Webhook returned client error, not retrying")

			return lastErr
		}
	}

	s.logger.Error().
		Err(lastErr).
		Str("url", url).
		Str("event_type", eventType).
		Int("max_retries", maxRetries).
		Msg("Webhook failed after all retry attempts")

	return fmt.Errorf("webhook failed after %d attempts: %w", maxRetries, lastErr)
}
