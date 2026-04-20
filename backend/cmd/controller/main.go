package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/api"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/config"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/system"
)

func main() {
	// Load configuration
	_ = config.LoadConfig() // Still load config in case we add other fields later

	// Initialize System Executor
	executor := system.NewExecutor()

	// Initialize API and WebSocket Hub
	apiHandler := api.NewAPI(executor)
	wsHub := api.NewWSHub()

	// Start WebSocket Hub in background
	go wsHub.Run()

	// Start Telemetry Worker
	StartTelemetryWorker(wsHub)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Register API Routes
	apiHandler.RegisterRoutes(r)

	// Register WebSocket endpoint
	r.GET("/ws", wsHub.HandleWebSocket)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
