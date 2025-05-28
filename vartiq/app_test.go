package vartiq

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// Helper to create an AppService with a mock client
func newMockAppService() (*AppService, *resty.Client) {
	r := resty.New()
	c := &Client{resty: r}
	return &AppService{client: c}, r
}

func TestAppService_Create(t *testing.T) {
	as, _ := newMockAppService()
	ctx := context.Background()
	_, err := as.Create(ctx, &CreateAppRequest{
		Name:      "Test",
		ProjectID: "project-123",
	})
	assert.Error(t, err) // No server, should error
}

func TestAppService_List(t *testing.T) {
	as, _ := newMockAppService()
	ctx := context.Background()
	_, err := as.List(ctx, "project-123")
	assert.Error(t, err)
}

func TestAppService_Get(t *testing.T) {
	as, _ := newMockAppService()
	ctx := context.Background()
	_, err := as.Get(ctx, "id")
	assert.Error(t, err)
}

func TestAppService_Update(t *testing.T) {
	as, _ := newMockAppService()
	ctx := context.Background()
	_, err := as.Update(ctx, "id", &UpdateAppRequest{Name: "New"})
	assert.Error(t, err)
}

func TestAppService_Delete(t *testing.T) {
	as, _ := newMockAppService()
	ctx := context.Background()
	err := as.Delete(ctx, "id")
	assert.Error(t, err)
}
