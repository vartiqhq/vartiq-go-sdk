package vartiq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookMessageIntegration(t *testing.T) {
	client := setupIntegrationTest(t)
	ctx := context.Background()

	// First create a project and app to send webhook messages
	projectName := "Test Project for WebhookMessage " + time.Now().Format(time.RFC3339)
	project, err := client.Project.Create(ctx, &CreateProjectRequest{
		Name:        projectName,
		Description: "Integration test project for webhook message tests",
	})
	require.NoError(t, err)
	projectID := project.Data.ID
	defer func() {
		err := client.Project.Delete(ctx, projectID)
		assert.NoError(t, err, "Failed to cleanup test project")
	}()

	// Create an app
	appName := "Test App for WebhookMessage " + time.Now().Format(time.RFC3339)
	app, err := client.App.Create(ctx, &CreateAppRequest{
		Name:        appName,
		ProjectID:   projectID,
		Description: "Integration test app for webhook message",
	})
	require.NoError(t, err)
	appID := app.Data.ID
	defer func() {
		err := client.App.Delete(ctx, appID)
		assert.NoError(t, err, "Failed to cleanup test app")
	}()

	// Create a webhook
	webhookName := "Test Webhook for Messages " + time.Now().Format(time.RFC3339)
	webhook, err := client.Webhook.Create(ctx, &CreateWebhookRequest{
		Name:  webhookName,
		URL:   "https://example.com/webhook",
		AppID: appID,
		CustomHeaders: []Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, webhook)
	webhookID := webhook.Data.ID
	defer func() {
		err := client.Webhook.Delete(ctx, webhookID)
		assert.NoError(t, err, "Failed to cleanup test webhook")
	}()

	// Test WebhookMessage Creation with different payload types
	testCases := []struct {
		name    string
		payload interface{}
	}{
		{
			name: "string payload",
			payload: map[string]interface{}{
				"message": "Hello, World!",
			},
		},
		{
			name: "complex payload",
			payload: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   float64(123),
					"name": "Test User",
				},
				"action": "signup",
				"metadata": map[string]interface{}{
					"source":    "web",
					"timestamp": float64(time.Now().Unix()),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message, err := client.WebhookMessage.Create(ctx, appID, tc.payload)
			require.NoError(t, err)
			require.NotNil(t, message)
			assert.Equal(t, appID, message.Data.AppID)
			assert.NotEmpty(t, message.Data.ID)
			assert.NotEmpty(t, message.Data.CreatedAt)
			assert.NotEmpty(t, message.Data.Signature)
			assert.Equal(t, tc.payload, message.Data.Payload)
			assert.True(t, message.Success)
			assert.NotEmpty(t, message.Message)
		})
	}
}
