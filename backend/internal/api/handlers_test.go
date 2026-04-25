package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockExecutor implements system.CommandExecutor for testing
type MockExecutor struct {
	LastScriptName string
	LastArgs       []string
	ReturnOutput   string
	ReturnError    error
}

func (m *MockExecutor) Run(scriptName string, args ...string) (string, error) {
	m.LastScriptName = scriptName
	m.LastArgs = args
	return m.ReturnOutput, m.ReturnError
}

func TestHandleCameraman(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		action         string
		device         string
		useQueryParam  bool
		executorOutput string
		executorError  error
		expectedStatus int
		expectedArgs   []string
	}{
		{
			name:           "Happy Path - JSON",
			action:         "start",
			device:         "V0A1",
			useQueryParam:  false,
			executorOutput: "starting cameraman",
			expectedStatus: http.StatusOK,
			expectedArgs:   []string{"V0A1", "start"},
		},
		{
			name:           "Happy Path - Query Param",
			action:         "stop",
			device:         "V0A2",
			useQueryParam:  true,
			executorOutput: "stopping cameraman",
			expectedStatus: http.StatusOK,
			expectedArgs:   []string{"V0A2", "stop"},
		},
		{
			name:           "Happy Path - No Device",
			action:         "status",
			device:         "",
			executorOutput: "cameraman status",
			expectedStatus: http.StatusOK,
			expectedArgs:   []string{"status"},
		},
		{
			name:           "Happy Path - List Devices",
			action:         "list-dev",
			device:         "",
			executorOutput: "dev1\ndev2",
			expectedStatus: http.StatusOK,
			expectedArgs:   []string{"list-dev"},
		},
		{
			name:           "Error - Invalid Action",
			action:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Invalid Device Pattern",
			action:         "start",
			device:         "; rm -rf /",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Invalid Device Pattern (JSON)",
			action:         "start",
			device:         "V1A1; echo test",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Executor Failure",
			action:         "start",
			device:         "V0A1",
			executorError:  fmt.Errorf("script failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &MockExecutor{
				ReturnOutput: tt.executorOutput,
				ReturnError:  tt.executorError,
			}
			api := NewAPI(mockExec)
			router := gin.New()
			router.POST("/api/cameraman/:action", api.handleCameraman)

			var req *http.Request
			url := fmt.Sprintf("/api/cameraman/%s", tt.action)

			if tt.useQueryParam {
				if tt.device != "" {
					url = fmt.Sprintf("%s?device=%s", url, tt.device)
				}
				req, _ = http.NewRequest(http.MethodPost, url, nil)
			} else {
				body, _ := json.Marshal(CameramanRequest{Device: tt.device})
				req, _ = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "VLX_cameraman", mockExec.LastScriptName)
				assert.Equal(t, tt.expectedArgs, mockExec.LastArgs)

				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.executorOutput, resp["output"])
			}
		})
	}
}
