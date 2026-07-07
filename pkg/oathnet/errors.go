package oathnet

import (
	"encoding/json"
	"fmt"
)

// OathNetError is the base error type.
type OathNetError struct {
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

func (e *OathNetError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("OathNet error (%d): %s (%s)", e.StatusCode, e.Message, formatErrorDetails(e.Details))
	}
	return fmt.Sprintf("OathNet error (%d): %s", e.StatusCode, e.Message)
}

// AuthenticationError is returned for authentication failures.
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

// ValidationError is returned for validation failures.
type ValidationError struct {
	Message string
	Details map[string]interface{}
}

func (e *ValidationError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("Validation error: %s (%s)", e.Message, formatErrorDetails(e.Details))
	}
	return fmt.Sprintf("Validation error: %s", e.Message)
}

// NotFoundError is returned when a resource is not found.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found: %s", e.Message)
}

// RateLimitError is returned when rate limited.
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit: %s", e.Message)
}

// QuotaExceededError is returned when quota is exceeded.
type QuotaExceededError struct {
	Message string
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("Quota exceeded: %s", e.Message)
}

func formatErrorDetails(details map[string]interface{}) string {
	data, err := json.Marshal(details)
	if err != nil {
		return fmt.Sprintf("%v", details)
	}
	return string(data)
}
