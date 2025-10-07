package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		return errors.New("webhook URL cannot be empty")
	}

	if event == nil {
		return errors.New("webhook event cannot be nil")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
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
	config := s.getRetryConfig()
	var lastErr error

	for attempt := 0; attempt < config.maxRetries; attempt++ {
		if err := s.waitForRetryDelay(ctx, attempt, config.retryDelays); err != nil {
			return err
		}

		s.logAttempt(url, eventType, attempt+1, config.maxRetries)

		result, err := s.executeSingleAttempt(req, url, eventType, attempt+1, config.maxRetries)
		if result.shouldReturn {
			return result.err
		}

		lastErr = err
	}

	return s.handleAllAttemptsExhausted(lastErr, url, eventType, config.maxRetries)
}

type retryConfig struct {
	maxRetries  int
	retryDelays []time.Duration
}

type attemptResult struct {
	shouldReturn bool
	err          error
}

func (s *HTTPWebhookSender) getRetryConfig() retryConfig {
	return retryConfig{
		maxRetries: 3,
		retryDelays: []time.Duration{
			0,
			5 * time.Second,
			15 * time.Second,
		},
	}
}

func (s *HTTPWebhookSender) waitForRetryDelay(ctx context.Context, attempt int, retryDelays []time.Duration) error {
	if attempt > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelays[attempt]):
		}
	}
	return nil
}

func (s *HTTPWebhookSender) logAttempt(url, eventType string, attempt, maxRetries int) {
	s.logger.Debug().
		Str("url", url).
		Str("event_type", eventType).
		Int("attempt", attempt).
		Int("max_retries", maxRetries).
		Msg("Sending webhook")
}

func (s *HTTPWebhookSender) executeSingleAttempt(req *http.Request, url, eventType string, attempt, maxRetries int) (attemptResult, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logRequestFailure(err, url, eventType, attempt)
		return attemptResult{shouldReturn: false}, fmt.Errorf("HTTP request failed (attempt %d/%d): %w", attempt, maxRetries, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if s.isSuccessStatusCode(resp.StatusCode) {
		s.logSuccess(url, eventType, resp.StatusCode, attempt)
		return attemptResult{shouldReturn: true, err: nil}, nil
	}

	err = fmt.Errorf("webhook returned status %d (attempt %d/%d)", resp.StatusCode, attempt, maxRetries)
	s.logStatusError(url, eventType, resp.StatusCode, attempt)

	if s.isClientError(resp.StatusCode) {
		s.logClientError(url, eventType, resp.StatusCode)
		return attemptResult{shouldReturn: true, err: err}, nil
	}

	return attemptResult{shouldReturn: false}, err
}

func (s *HTTPWebhookSender) isSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func (s *HTTPWebhookSender) isClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func (s *HTTPWebhookSender) logRequestFailure(err error, url, eventType string, attempt int) {
	s.logger.Warn().
		Err(err).
		Str("url", url).
		Str("event_type", eventType).
		Int("attempt", attempt).
		Msg("Webhook request failed")
}

func (s *HTTPWebhookSender) logSuccess(url, eventType string, statusCode, attempt int) {
	s.logger.Info().
		Str("url", url).
		Str("event_type", eventType).
		Int("status_code", statusCode).
		Int("attempt", attempt).
		Msg("Webhook sent successfully")
}

func (s *HTTPWebhookSender) logStatusError(url, eventType string, statusCode, attempt int) {
	s.logger.Warn().
		Str("url", url).
		Str("event_type", eventType).
		Int("status_code", statusCode).
		Int("attempt", attempt).
		Msg("Webhook returned error status")
}

func (s *HTTPWebhookSender) logClientError(url, eventType string, statusCode int) {
	s.logger.Error().
		Str("url", url).
		Str("event_type", eventType).
		Int("status_code", statusCode).
		Msg("Webhook returned client error, not retrying")
}

func (s *HTTPWebhookSender) handleAllAttemptsExhausted(lastErr error, url, eventType string, maxRetries int) error {
	s.logger.Error().
		Err(lastErr).
		Str("url", url).
		Str("event_type", eventType).
		Int("max_retries", maxRetries).
		Msg("Webhook failed after all retry attempts")

	return fmt.Errorf("webhook failed after %d attempts: %w", maxRetries, lastErr)
}
