package vartiq

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// Helper to create a WebhookMessageService with a mock client
func newMockWebhookMessageService() (*WebhookMessageService, *resty.Client) {
	r := resty.New()
	c := &Client{resty: r}
	return &WebhookMessageService{client: c}, r
}

func TestWebhookMessageService_Create(t *testing.T) {
	wms, _ := newMockWebhookMessageService()
	ctx := context.Background()

	payload := map[string]interface{}{
		"hello": "world",
	}

	message, err := wms.Create(ctx, "app-123", payload)
	assert.Error(t, err) // No server, should error
	assert.Nil(t, message)
}
