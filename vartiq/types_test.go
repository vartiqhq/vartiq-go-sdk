package vartiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{Message: "something went wrong", Code: 400}
	assert.Equal(t, "something went wrong", err.Error())
}
