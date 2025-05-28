package vartiq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectIntegration(t *testing.T) {
	client := setupIntegrationTest(t)
	ctx := context.Background()

	// Test Project Creation
	projectName := "Test Project " + time.Now().Format(time.RFC3339)
	project, err := client.Project.Create(ctx, &CreateProjectRequest{
		Name:        projectName,
		Description: "Integration test project",
	})
	require.NoError(t, err)
	require.NotNil(t, project)
	assert.Equal(t, projectName, project.Data.Name)

	// Store project ID for cleanup
	projectID := project.Data.ID
	defer func() {
		err := client.Project.Delete(ctx, projectID)
		assert.NoError(t, err, "Failed to cleanup test project")
	}()

	// Test Project Get
	retrieved, err := client.Project.Get(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, projectID, retrieved.Data.ID)
	assert.Equal(t, projectName, retrieved.Data.Name)
	assert.True(t, retrieved.Success)
	assert.NotEmpty(t, retrieved.Message)

	// Test Project List
	projects, err := client.Project.List(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, projects.Data)
	found := false
	for _, p := range projects.Data {
		if p.ID == projectID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created project not found in list")
	assert.True(t, projects.Success)
	assert.NotEmpty(t, projects.Message)

	// Test Project Update
	updatedName := "Updated " + projectName
	updated, err := client.Project.Update(ctx, projectID, &UpdateProjectRequest{
		Name: updatedName,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.Data.Name)
	assert.True(t, updated.Success)
	assert.NotEmpty(t, updated.Message)
}
