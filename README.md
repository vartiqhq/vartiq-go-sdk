# vartiq-go-sdk

A Go SDK for interacting with the Vartiq API. Supports Project, App, Webhook, and Webhook Message resources.

## Installation

```sh
go get github.com/vartiqhq/vartiq-go-sdk
```

## Usage

### Import and Initialize

```go
import (
	"github.com/vartiqhq/vartiq-go-sdk/vartiq"
)

client := vartiq.New("YOUR_API_KEY")
```

### Go Types

You can import types for strong typing:

```go
import (
	"github.com/vartiqhq/vartiq-go-sdk/vartiq"
)
// vartiq.Project, vartiq.App, vartiq.Webhook, vartiq.WebhookMessage
```

## API

### Project

```go
// Create a project
projectResp, err := client.Project.Create(ctx, &vartiq.CreateProjectRequest{
	Name:        "Test",
	Description: "desc",
})

// Get all projects
projects, err := client.Project.List(ctx)

// Get a single project
project, err := client.Project.Get(ctx, "PROJECT_ID")

// Update a project
updated, err := client.Project.Update(ctx, "PROJECT_ID", &vartiq.UpdateProjectRequest{
	Name: "New Name",
})

// Delete a project
err := client.Project.Delete(ctx, "PROJECT_ID")
```

### App

```go
// Create an app
appResp, err := client.App.Create(ctx, &vartiq.CreateAppRequest{
	Name:      "App Name",
	ProjectID: "PROJECT_ID",
})

// Get all apps for a project
apps, err := client.App.List(ctx, "PROJECT_ID")

// Get a single app
app, err := client.App.Get(ctx, "APP_ID")

// Update an app
updated, err := client.App.Update(ctx, "APP_ID", &vartiq.UpdateAppRequest{
	Name: "New App Name",
})

// Delete an app
err := client.App.Delete(ctx, "APP_ID")
```

### Webhook

```go
// Create a webhook
webhookResp, err := client.Webhook.Create(ctx, &vartiq.CreateWebhookRequest{
	Name:   "Webhook",
	URL:    "https://your-webhook-url.com",
	AppID:  "APP_ID",
	CustomHeaders: []vartiq.Header{{Key: "x-app", Value: "x-value"}}, // optional
})

// Get all webhooks for an app
webhooks, err := client.Webhook.GetAll(ctx, "APP_ID")

// Get a single webhook
webhook, err := client.Webhook.GetOne(ctx, "WEBHOOK_ID")

// Update a webhook
updated, err := client.Webhook.Update(ctx, "WEBHOOK_ID", map[string]interface{}{
	"name": "New Webhook Name",
})

// Delete a webhook
err := client.Webhook.Delete(ctx, "WEBHOOK_ID")
```

### Webhook Message

The WebhookMessage service allows you to programmatically send messages to your webhooks.

```go
// Create a webhook message
message, err := client.WebhookMessage.Create(ctx, "APP_ID", map[string]interface{}{
	"hello": "world",
})
```

### Webhook Verification

To verify a webhook signature, you can use the `Verify` method. This is useful for ensuring that incoming webhooks are genuinely from Vartiq and have not been tampered with.

```go
import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vartiqhq/vartiq-go-sdk/vartiq"
)

// Assuming you have an http.HandlerFunc to receive webhooks
func handleWebhook(w http.ResponseWriter, req *http.Request) {
	client := vartiq.New("YOUR_API_KEY") // Or use your existing client instance

	webhookSecret := "YOUR_WEBHOOK_SECRET" // Retrieve your webhook secret securely

	signature := req.Header.Get("X-Vartiq-Signature") // Get the signature from the header
	if signature == "" {
		http.Error(w, "Missing signature header", http.StatusBadRequest)
		return
	}

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	verifiedPayload, err := client.Verify(payload, signature, webhookSecret)
	if err != nil {
		// Signature is invalid
		fmt.Printf("Webhook verification failed: %v\n", err)
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}

	// Signature is valid, verifiedPayload contains the raw payload bytes
	fmt.Printf("Webhook verified successfully. Payload: %s\n", string(verifiedPayload))

	// Process the verified payload...

	w.WriteHeader(http.StatusOK)
}
```

The `Verify` function returns the original payload bytes if the signature is valid. If the signature is invalid or missing, it returns an error.
