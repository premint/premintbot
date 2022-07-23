package premint

import (
	"testing"

	"go.uber.org/zap"
)

func TestProvidePremint(t *testing.T) {
	t.Run("test ProvidePremint", func(t *testing.T) {
		premintClient := ProvidePremint()

		if premintClient == nil {
			t.Error("Expected premint client to be non-nil")
		}

		if premintClient.httpClient == nil {
			t.Error("Expected premint client HTTP client to be non-nil")
		}

		if premintClient.httpClient.Transport == nil {
			t.Error("Expected premint client HTTP client transport to be non-nil")
		}
	})
}

func TestCheckPremintStatusForUser(t *testing.T) {
	t.Run("test CheckPremintStatusForUser", func(t *testing.T) {
		logger := zap.NewNop().Sugar()
		apiKey := "some-secret-api-key"
		userID := "360541062839926785"
		r, err := CheckPremintStatusForUser(logger, apiKey, userID)
		if err != nil {
			t.Error("Expected no error")
		}

	})
}
