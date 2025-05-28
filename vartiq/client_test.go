package vartiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_DefaultBaseURL(t *testing.T) {
	apiKey := "test-key"
	client := New(apiKey)
	assert.Equal(t, "https://api.us.vartiq.com", client.baseURL)
	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.resty)
	assert.NotNil(t, client.Project)
	assert.NotNil(t, client.App)
	assert.NotNil(t, client.Webhook)
	assert.NotNil(t, client.WebhookMessage)
}

func TestNewClient_CustomBaseURL(t *testing.T) {
	apiKey := "test-key"
	baseURL := "https://custom.example.com"
	client := New(apiKey, baseURL)
	assert.Equal(t, baseURL, client.baseURL)
}
