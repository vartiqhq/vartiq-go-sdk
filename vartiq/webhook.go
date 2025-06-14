package vartiq

import (
	"context"
	"encoding/json"
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
	HMACHeader   string     `json:"hmacHeader,omitempty"`
	HMACSecret   string     `json:"hmacSecret,omitempty"`
	APIKey       string     `json:"apiKey,omitempty"`
	APIKeyHeader string     `json:"apiKeyHeader,omitempty"`
	UserName     string     `json:"userName,omitempty"`
	Password     string     `json:"password,omitempty"`
}

type Webhook struct {
	ID            string       `json:"id"`
	URL           string       `json:"url"`
	AppID         string       `json:"appId"`
	Company       string       `json:"company"`
	CustomHeaders []Header     `json:"customHeaders"`
	Headers       []Header     `json:"headers"`
	AuthMethod    *WebhookAuth `json:"authMethod,omitempty"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateWebhookRequest struct {
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

	// Send the request exactly as provided, since it matches the server's validation schema
	resp := &WebhookResponse{}
	httpResp, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(req). // Send the request directly without restructuring
		SetResult(resp).
		Post("/webhooks")
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	// Debug logging
	fmt.Printf("Webhook Create Response Status: %s\n", httpResp.Status())
	fmt.Printf("Webhook Create Response Body: %s\n", string(httpResp.Body()))

	if !resp.Success {
		return nil, fmt.Errorf("webhook creation failed: %s", resp.Message)
	}

	return resp, nil
}

func (s *WebhookService) GetAll(ctx context.Context, appID string) (*WebhookListResponse, error) {
	resp := &WebhookListResponse{}
	httpResp, err := s.client.resty.R().
		SetContext(ctx).
		SetQueryParam("appId", appID).
		SetResult(resp).
		Get("/webhooks")
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	// Debug logging
	fmt.Printf("Webhook List Response Status: %s\n", httpResp.Status())
	fmt.Printf("Webhook List Response Body: %s\n", string(httpResp.Body()))

	if !resp.Success {
		return nil, fmt.Errorf("webhook list retrieval failed: %s", resp.Message)
	}

	// Ensure each webhook has its auth method properly set
	for i := range resp.Data {
		if resp.Data[i].AuthMethod == nil {
			// If authMethod is nil but we have auth data in the response,
			// try to reconstruct it from the raw response
			var rawData map[string]interface{}
			if err := json.Unmarshal(httpResp.Body(), &rawData); err == nil {
				if data, ok := rawData["data"].([]interface{}); ok && i < len(data) {
					if webhook, ok := data[i].(map[string]interface{}); ok {
						if authMethod, ok := webhook["authMethod"].(map[string]interface{}); ok {
							resp.Data[i].AuthMethod = &WebhookAuth{
								Method:     AuthMethod(authMethod["method"].(string)),
								HMACHeader: authMethod["hmacHeader"].(string),
								HMACSecret: authMethod["hmacSecret"].(string),
							}
						}
					}
				}
			}
		}
	}

	return resp, nil
}

func (s *WebhookService) GetOne(ctx context.Context, webhookID string) (*WebhookResponse, error) {
	resp := &WebhookResponse{}
	httpResp, err := s.client.resty.R().
		SetContext(ctx).
		SetResult(resp).
		Get("/webhooks/" + webhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	// Debug logging
	fmt.Printf("Webhook Get Response Status: %s\n", httpResp.Status())
	fmt.Printf("Webhook Get Response Body: %s\n", string(httpResp.Body()))

	if !resp.Success {
		return nil, fmt.Errorf("webhook retrieval failed: %s", resp.Message)
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
