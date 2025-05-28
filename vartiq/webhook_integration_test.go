package vartiq

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookIntegration(t *testing.T) {
	client := setupIntegrationTest(t)
	ctx := context.Background()

	// First create a project and app to hold the webhook
	projectName := "Test Project for Webhook " + time.Now().Format(time.RFC3339)
	project, err := client.Project.Create(ctx, &CreateProjectRequest{
		Name:        projectName,
		Description: "Integration test project for webhook tests",
	})
	require.NoError(t, err)
	projectID := project.Data.ID
	defer func() {
		err := client.Project.Delete(ctx, projectID)
		assert.NoError(t, err, "Failed to cleanup test project")
	}()

	// Create an app
	appName := "Test App for Webhook " + time.Now().Format(time.RFC3339)
	app, err := client.App.Create(ctx, &CreateAppRequest{
		Name:        appName,
		ProjectID:   projectID,
		Description: "Integration test app for webhook",
	})
	require.NoError(t, err)
	appID := app.Data.ID
	defer func() {
		err := client.App.Delete(ctx, appID)
		assert.NoError(t, err, "Failed to cleanup test app")
	}()

	// Test Webhook Creation
	webhookName := "Test Webhook " + time.Now().Format(time.RFC3339)
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
	assert.Equal(t, webhookName, webhook.Data.Name)

	// Store webhook ID for cleanup
	webhookID := webhook.Data.ID
	defer func() {
		err := client.Webhook.Delete(ctx, webhookID)
		assert.NoError(t, err, "Failed to cleanup test webhook")
	}()

	// Test Webhook Get
	retrieved, err := client.Webhook.GetOne(ctx, webhookID)
	require.NoError(t, err)
	assert.Equal(t, webhookID, retrieved.Data.ID)
	assert.Equal(t, webhookName, retrieved.Data.Name)
	assert.True(t, retrieved.Success)
	assert.NotEmpty(t, retrieved.Message)

	// Test Webhook List
	webhooks, err := client.Webhook.GetAll(ctx, appID)
	require.NoError(t, err)
	assert.NotEmpty(t, webhooks.Data)
	found := false
	for _, w := range webhooks.Data {
		if w.ID == webhookID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created webhook not found in list")
	assert.True(t, webhooks.Success)
	assert.NotEmpty(t, webhooks.Message)

	// Test Webhook Update
	updatedName := "Updated " + webhookName
	updated, err := client.Webhook.Update(ctx, webhookID, map[string]interface{}{
		"name": updatedName,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.Data.Name)
	assert.True(t, updated.Success)
	assert.NotEmpty(t, updated.Message)

	// Test Webhook Verification
	payload := []byte(`{"test":"data"}`)
	secret := webhook.Data.Secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	verifiedPayload, err := client.Verify(payload, signature, secret)
	require.NoError(t, err)
	assert.Equal(t, payload, verifiedPayload)
}
