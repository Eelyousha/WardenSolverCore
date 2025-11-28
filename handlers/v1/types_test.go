package v1

import (
	"testing"
)

// TestSocialNetworksValidation тестирует валидацию социальных сетей
func TestSocialNetworksValidation(t *testing.T) {
	tests := []struct {
		name    string
		network string
		valid   bool
	}{
		{"X is valid", "x", true},
		{"Telegram is valid", "telegram", true},
		{"Twitch is valid", "twitch", true},
		{"Facebook is invalid", "facebook", false},
		{"Empty is invalid", "", false},
		{"YouTube is invalid", "youtube", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if SocialNetworks[tt.network] != tt.valid {
				t.Errorf("SocialNetworks[%q] = %v, want %v", tt.network, SocialNetworks[tt.network], tt.valid)
			}
		})
	}
}

// TestResponsePayloadStructure тестирует структуру ResponsePayload
func TestResponsePayloadStructure(t *testing.T) {
	payload := ResponsePayload{
		Status:  "success",
		Message: "Test message",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	if payload.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", payload.Status)
	}

	if payload.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", payload.Message)
	}

	if payload.Data == nil {
		t.Error("Expected Data to be non-nil")
	}
}

// TestResponsePayloadWithoutData тестирует ResponsePayload без Data
func TestResponsePayloadWithoutData(t *testing.T) {
	payload := ResponsePayload{
		Status:  "error",
		Message: "Error occurred",
	}

	if payload.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", payload.Status)
	}

	if payload.Data != nil {
		t.Errorf("Expected Data to be nil, got %v", payload.Data)
	}
}
