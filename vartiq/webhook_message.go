package vartiq

import (
	"context"
	"encoding/json"
	"fmt"
)

// Error represents an API error response
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

type WebhookMessageService struct {
	client *Client
}

type WebhookMessage struct {
	ID          string      `json:"id"`
	AppID       string      `json:"app"`
	Payload     interface{} `json:"payload"`
	Signature   string      `json:"signature"`
	IsDelivered bool        `json:"isDelivered"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type webhookMessageResponse struct {
	Data struct {
		WebhookMessages []struct {
			ID      string `json:"id"`
			AppID   string `json:"app"`
			Payload string `json:"payload"` // API returns payload as JSON string
			Headers []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
			IsDelivered bool   `json:"isDelivered"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"webhookMessages"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type WebhookMessageResponse struct {
	Data    WebhookMessage `json:"data"`
	Message string         `json:"message"`
	Success bool           `json:"success"`
}

// Create sends a message to a webhook. The payload can be any JSON-serializable value.
// Example:
//
//	message, err := client.WebhookMessage.Create(ctx, "APP_ID", map[string]interface{}{
//	    "hello": "world",
//	})
func (s *WebhookMessageService) Create(ctx context.Context, appID string, payload interface{}) (*WebhookMessageResponse, error) {
	resp := &webhookMessageResponse{}
	httpResp, err := s.client.resty.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"appId":   appID,
			"payload": payload,
		}).
		SetResult(resp).
		Post("/webhook-messages")
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Debug logging
	fmt.Printf("Response Status: %s\n", httpResp.Status())
	fmt.Printf("Response Body: %s\n", string(httpResp.Body()))

	if !resp.Success {
		return nil, &Error{Message: resp.Message}
	}

	if len(resp.Data.WebhookMessages) == 0 {
		return nil, &Error{Message: "No webhook messages returned"}
	}

	rawMessage := resp.Data.WebhookMessages[0]

	// Parse the payload JSON string back into an interface{}
	var parsedPayload interface{}
	if err := json.Unmarshal([]byte(rawMessage.Payload), &parsedPayload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Extract signature from headers
	var signature string
	for _, header := range rawMessage.Headers {
		if header.Key == "x-Vartiq-signature" {
			signature = header.Value
			break
		}
	}

	message := WebhookMessage{
		ID:          rawMessage.ID,
		AppID:       appID, // Use the provided appID
		Payload:     parsedPayload,
		Signature:   signature,
		IsDelivered: rawMessage.IsDelivered,
		CreatedAt:   rawMessage.CreatedAt,
		UpdatedAt:   rawMessage.UpdatedAt,
	}

	return &WebhookMessageResponse{
		Data:    message,
		Message: resp.Message,
		Success: resp.Success,
	}, nil
}
