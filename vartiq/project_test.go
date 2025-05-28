package vartiq

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// Helper to create a ProjectService with a mock client
func newMockProjectService() (*ProjectService, *resty.Client) {
	r := resty.New()
	c := &Client{resty: r}
	return &ProjectService{client: c}, r
}

func TestProjectService_Create(t *testing.T) {
	ps, _ := newMockProjectService()
	ctx := context.Background()
	_, err := ps.Create(ctx, &CreateProjectRequest{Name: "Test", Description: "Desc"})
	assert.Error(t, err) // No server, should error
}

func TestProjectService_List(t *testing.T) {
	ps, _ := newMockProjectService()
	ctx := context.Background()
	_, err := ps.List(ctx)
	assert.Error(t, err)
}

func TestProjectService_Get(t *testing.T) {
	ps, _ := newMockProjectService()
	ctx := context.Background()
	_, err := ps.Get(ctx, "id")
	assert.Error(t, err)
}

func TestProjectService_Update(t *testing.T) {
	ps, _ := newMockProjectService()
	ctx := context.Background()
	_, err := ps.Update(ctx, "id", &UpdateProjectRequest{Name: "New"})
	assert.Error(t, err)
}

func TestProjectService_Delete(t *testing.T) {
	ps, _ := newMockProjectService()
	ctx := context.Background()
	err := ps.Delete(ctx, "id")
	assert.Error(t, err)
}
