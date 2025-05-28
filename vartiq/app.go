package vartiq

import (
	"context"
)

type AppService struct {
	client *Client
}

type App struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Company     string `json:"company"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreateAppRequest struct {
	Name        string `json:"name"`
	ProjectID   string `json:"projectId"`
	Description string `json:"description,omitempty"`
}

type CreateAppResponse struct {
	Data    App    `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (s *AppService) Create(ctx context.Context, req *CreateAppRequest) (*CreateAppResponse, error) {
	resp := &CreateAppResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(resp).
		Post("/apps")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// List all apps for a project
func (s *AppService) List(ctx context.Context, projectID string) (*struct {
	Data    []App  `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}, error) {
	resp := &struct {
		Data    []App  `json:"data"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/apps?projectId=" + projectID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get a single app by ID
func (s *AppService) Get(ctx context.Context, appID string) (*struct {
	Data    App    `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}, error) {
	resp := &struct {
		Data    App    `json:"data"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/apps/" + appID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateAppRequest is used for updating an app
type UpdateAppRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update an app by ID
func (s *AppService) Update(ctx context.Context, appID string, req *UpdateAppRequest) (*struct {
	Data    App    `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}, error) {
	resp := &struct {
		Data    App    `json:"data"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(resp).
		Put("/apps/" + appID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete an app by ID
func (s *AppService) Delete(ctx context.Context, appID string) error {
	_, err := s.client.resty.R().
		SetContext(ctx).
		Delete("/apps/" + appID)
	return err
}
