package vartiq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppIntegration(t *testing.T) {
	client := setupIntegrationTest(t)
	ctx := context.Background()

	// First create a project to hold the app
	projectName := "Test Project for App " + time.Now().Format(time.RFC3339)
	project, err := client.Project.Create(ctx, &CreateProjectRequest{
		Name:        projectName,
		Description: "Integration test project for app tests",
	})
	require.NoError(t, err)
	projectID := project.Data.ID
	defer func() {
		err := client.Project.Delete(ctx, projectID)
		assert.NoError(t, err, "Failed to cleanup test project")
	}()

	// Test App Creation
	appName := "Test App " + time.Now().Format(time.RFC3339)
	app, err := client.App.Create(ctx, &CreateAppRequest{
		Name:        appName,
		ProjectID:   projectID,
		Description: "Integration test app",
	})
	require.NoError(t, err)
	require.NotNil(t, app)
	assert.Equal(t, appName, app.Data.Name)

	// Store app ID for cleanup
	appID := app.Data.ID
	defer func() {
		err := client.App.Delete(ctx, appID)
		assert.NoError(t, err, "Failed to cleanup test app")
	}()

	// Test App Get
	retrieved, err := client.App.Get(ctx, appID)
	require.NoError(t, err)
	assert.Equal(t, appID, retrieved.Data.ID)
	assert.Equal(t, appName, retrieved.Data.Name)
	assert.True(t, retrieved.Success)
	assert.NotEmpty(t, retrieved.Message)

	// Test App List
	apps, err := client.App.List(ctx, projectID)
	require.NoError(t, err)
	assert.NotEmpty(t, apps.Data)
	found := false
	for _, a := range apps.Data {
		if a.ID == appID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created app not found in list")
	assert.True(t, apps.Success)
	assert.NotEmpty(t, apps.Message)

	// Test App Update
	updatedName := "Updated " + appName
	updated, err := client.App.Update(ctx, appID, &UpdateAppRequest{
		Name: updatedName,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.Data.Name)
	assert.True(t, updated.Success)
	assert.NotEmpty(t, updated.Message)
}
