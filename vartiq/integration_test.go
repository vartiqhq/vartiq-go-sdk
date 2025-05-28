package vartiq

import (
	"os"
	"testing"
)

var (
	testClient *Client
	testAPIKey string
)

func skipIfNoAPIKey(t *testing.T) {
	if os.Getenv("VARTIQ_API_KEY") == "" {
		t.Skip("Skipping integration test: VARTIQ_API_KEY not set")
	}
}

func setupIntegrationTest(t *testing.T) *Client {
	skipIfNoAPIKey(t)

	if testClient == nil {
		testAPIKey = os.Getenv("VARTIQ_API_KEY")
		baseURL := os.Getenv("VARTIQ_API_URL")
		if baseURL == "" {
			testClient = New(testAPIKey)
		} else {
			testClient = New(testAPIKey, baseURL)
		}
	}
	return testClient
}
