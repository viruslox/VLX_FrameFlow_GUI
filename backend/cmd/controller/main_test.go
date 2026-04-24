package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/api"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/config"
)

// MockExecutor implements system.CommandExecutor for testing
type MockExecutor struct{}

func (m *MockExecutor) Run(scriptName string, args ...string) (string, error) {
	return "mock output", nil
}

func TestRouterAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		AuthUser: "testuser",
		AuthPass: "testpass",
	}

	mockExec := &MockExecutor{}
	apiHandler := api.NewAPI(mockExec)
	wsHub := api.NewWSHub()

	router := setupRouter(cfg, apiHandler, wsHub)

	tests := []struct {
		name           string
		method         string
		path           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Health Check - Unauthenticated (Allowed)",
			method:         http.MethodGet,
			path:           "/health",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "API Route - Unauthenticated (Denied)",
			method:         http.MethodGet,
			path:           "/api/frameflow/bonding",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "API Route - Incorrect Credentials (Denied)",
			method:         http.MethodGet,
			path:           "/api/frameflow/bonding",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("wrong:creds")),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "API Route - Correct Credentials (Allowed)",
			method:         http.MethodGet,
			path:           "/api/frameflow/bonding",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass")),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
