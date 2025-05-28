package vartiq

import (
	"context"
)

type ProjectService struct {
	client *Client
}

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Company     string `json:"company"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectResponse struct {
	Data    Project `json:"data"`
	Message string  `json:"message"`
	Success bool    `json:"success"`
}

func (s *ProjectService) Create(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	resp := &CreateProjectResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(resp).
		Post("/projects")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// List all projects
func (s *ProjectService) List(ctx context.Context) (*struct {
	Data    []Project `json:"data"`
	Message string    `json:"message"`
	Success bool      `json:"success"`
}, error) {
	resp := &struct {
		Data    []Project `json:"data"`
		Message string    `json:"message"`
		Success bool      `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/projects")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get a single project by ID
func (s *ProjectService) Get(ctx context.Context, projectID string) (*struct {
	Data    Project `json:"data"`
	Message string  `json:"message"`
	Success bool    `json:"success"`
}, error) {
	resp := &struct {
		Data    Project `json:"data"`
		Message string  `json:"message"`
		Success bool    `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/projects/" + projectID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateProjectRequest is used for updating a project
type UpdateProjectRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update a project by ID
func (s *ProjectService) Update(ctx context.Context, projectID string, req *UpdateProjectRequest) (*struct {
	Data    Project `json:"data"`
	Message string  `json:"message"`
	Success bool    `json:"success"`
}, error) {
	resp := &struct {
		Data    Project `json:"data"`
		Message string  `json:"message"`
		Success bool    `json:"success"`
	}{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(resp).
		Put("/projects/" + projectID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete a project by ID
func (s *ProjectService) Delete(ctx context.Context, projectID string) error {
	_, err := s.client.resty.R().
		SetContext(ctx).
		Delete("/projects/" + projectID)
	return err
}
