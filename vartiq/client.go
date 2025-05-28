// Package vartiq provides a Go SDK for the Vartiq API.
// You must provide an API key to use this client.
package vartiq

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Client represents a Vartiq API client
type Client struct {
	baseURL string
	apiKey  string
	resty   *resty.Client

	Project        *ProjectService
	App            *AppService
	Webhook        *WebhookService
	WebhookMessage *WebhookMessageService
}

// New creates a new Vartiq API client. If baseURL is not provided, it defaults to https://api.us.vartiq.com
func New(apiKey string, baseURL ...string) *Client {
	url := "https://api.us.vartiq.com"
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	}
	r := resty.New().SetBaseURL(url).SetHeader("x-api-key", apiKey)
	c := &Client{
		baseURL: url,
		apiKey:  apiKey,
		resty:   r,
	}
	c.Project = &ProjectService{client: c}
	c.App = &AppService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.WebhookMessage = &WebhookMessageService{client: c}
	return c
}

// Verify checks the signature of a webhook payload.
// It takes the raw payload bytes, the signature string from the header, and the webhook secret.
// It returns the payload bytes if the signature is valid, otherwise returns an error.
func (c *Client) Verify(payload []byte, signature, secret string) ([]byte, error) {
	if signature == "" {
		return nil, errors.New("signature header is missing")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := mac.Sum(nil)

	// Assuming the signature is hex encoded
	receivedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(receivedSignature, expectedSignature) != 1 {
		return nil, errors.New("signature verification failed")
	}

	return payload, nil
}
