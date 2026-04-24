package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/api"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/config"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/system"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/ui"
)

func setupRouter(cfg *config.Config, apiHandler *api.API, wsHub *api.WSHub) *gin.Engine {
	r := gin.Default()

	// Health check endpoint (unauthenticated)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Setup BasicAuth
	auth := gin.BasicAuth(gin.Accounts{
		cfg.AuthUser: cfg.AuthPass,
	})

	// Apply authentication middleware to all subsequent routes
	r.Use(auth)

	// Serve embedded frontend
	ui.ServeFrontend(r)

	// Register API Routes
	apiHandler.RegisterRoutes(r)

	// Register WebSocket endpoint
	r.GET("/ws", wsHub.HandleWebSocket)

	return r
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize System Executor
	executor := system.NewExecutor()

	// Initialize API and WebSocket Hub
	apiHandler := api.NewAPI(executor)
	wsHub := api.NewWSHub()

	// Start WebSocket Hub in background
	go wsHub.Run()

	// Start Telemetry Worker
	StartTelemetryWorker(wsHub)

	// Setup Router
	r := setupRouter(cfg, apiHandler, wsHub)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
