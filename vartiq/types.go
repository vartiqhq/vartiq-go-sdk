package vartiq

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}
