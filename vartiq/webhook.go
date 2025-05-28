package vartiq

import (
	"context"
	"errors"
	"fmt"
)

type WebhookService struct {
	client *Client
}

type AuthMethod string

const (
	AuthMethodAPIKey AuthMethod = "apiKey"
	AuthMethodBasic  AuthMethod = "basic"
	AuthMethodHMAC   AuthMethod = "hmac"
)

type WebhookAuth struct {
	Method       AuthMethod `json:"method"`
	UserName     string     `json:"userName,omitempty"`
	Password     string     `json:"password,omitempty"`
	APIKey       string     `json:"apiKey,omitempty"`
	APIKeyHeader string     `json:"apiKeyHeader,omitempty"`
	HMACHeader   string     `json:"hmacHeader,omitempty"`
	HMACSecret   string     `json:"hmacSecret,omitempty"`
}

type Webhook struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	URL           string       `json:"url"`
	AppID         string       `json:"app"`
	Secret        string       `json:"secret"`
	CustomHeaders []Header     `json:"customHeaders"`
	Headers       []Header     `json:"headers"`
	Auth          *WebhookAuth `json:"auth,omitempty"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateWebhookRequest struct {
	Name          string   `json:"name"`
	URL           string   `json:"url"`
	AppID         string   `json:"appId"`
	CustomHeaders []Header `json:"customHeaders,omitempty"`
	AuthMethod    string   `json:"authMethod,omitempty"`
	// Basic Auth
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	// API Key Auth
	APIKey       string `json:"apiKey,omitempty"`
	APIKeyHeader string `json:"apiKeyHeader,omitempty"`
	// HMAC Auth
	HMACHeader string `json:"hmacHeader,omitempty"`
	HMACSecret string `json:"hmacSecret,omitempty"`
}

type WebhookResponse struct {
	Data    Webhook `json:"data"`
	Message string  `json:"message"`
	Success bool    `json:"success"`
}

type WebhookListResponse struct {
	Data    []Webhook `json:"data"`
	Message string    `json:"message"`
	Success bool      `json:"success"`
}

func validateWebhookAuth(req *CreateWebhookRequest) error {
	if req.AuthMethod == "" {
		return nil
	}

	switch AuthMethod(req.AuthMethod) {
	case AuthMethodBasic:
		if req.UserName == "" || req.Password == "" {
			return errors.New("for basic auth, userName and password are required")
		}
	case AuthMethodHMAC:
		if req.HMACHeader == "" || req.HMACSecret == "" {
			return errors.New("for hmac auth, hmacHeader and hmacSecret are required")
		}
	case AuthMethodAPIKey:
		if req.APIKey == "" || req.APIKeyHeader == "" {
			return errors.New("for apiKey auth, apiKey and apiKeyHeader are required")
		}
	default:
		return fmt.Errorf("invalid auth method: %s", req.AuthMethod)
	}

	return nil
}

func (s *WebhookService) Create(ctx context.Context, req *CreateWebhookRequest) (*WebhookResponse, error) {
	if err := validateWebhookAuth(req); err != nil {
		return nil, err
	}

	// Convert the flattened request to the internal structure
	requestBody := struct {
		Name          string       `json:"name"`
		URL           string       `json:"url"`
		AppID         string       `json:"appId"`
		CustomHeaders []Header     `json:"customHeaders,omitempty"`
		Auth          *WebhookAuth `json:"auth,omitempty"`
	}{
		Name:          req.Name,
		URL:           req.URL,
		AppID:         req.AppID,
		CustomHeaders: req.CustomHeaders,
	}

	if req.AuthMethod != "" {
		requestBody.Auth = &WebhookAuth{
			Method:       AuthMethod(req.AuthMethod),
			UserName:     req.UserName,
			Password:     req.Password,
			APIKey:       req.APIKey,
			APIKeyHeader: req.APIKeyHeader,
			HMACHeader:   req.HMACHeader,
			HMACSecret:   req.HMACSecret,
		}
	}

	resp := &WebhookResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(requestBody).
		SetResult(resp).
		Post("/webhooks")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *WebhookService) GetAll(ctx context.Context, appID string) (*WebhookListResponse, error) {
	resp := &WebhookListResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetQueryParam("appId", appID).
		SetResult(resp).
		Get("/webhooks")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *WebhookService) GetOne(ctx context.Context, webhookID string) (*WebhookResponse, error) {
	resp := &WebhookResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/webhooks/" + webhookID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *WebhookService) Update(ctx context.Context, webhookID string, req map[string]interface{}) (*WebhookResponse, error) {
	resp := &WebhookResponse{}
	_, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(resp).
		Put("/webhooks/" + webhookID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *WebhookService) Delete(ctx context.Context, webhookID string) error {
	_, err := s.client.resty.R().
		SetContext(ctx).
		Delete("/webhooks/" + webhookID)
	return err
}
