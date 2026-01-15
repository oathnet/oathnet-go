package oathnet

import "fmt"

// OathNetError is the base error type.
type OathNetError struct {
	Message    string
	StatusCode int
}

func (e *OathNetError) Error() string {
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
}

func (e *ValidationError) Error() string {
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
