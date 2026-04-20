package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viruslox/VLX_FrameFlow_GUI/backend/internal/system"
)

type API struct {
	executor *system.Executor
}

func NewAPI(executor *system.Executor) *API {
	return &API{
		executor: executor,
	}
}

func (a *API) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// FrameFlow Client
	api.POST("/frameflow/client/:action", a.handleFrameFlowClient)

	// FrameFlow AP
	api.POST("/frameflow/ap/:action", a.handleFrameFlowAP)

	// FrameFlow Bonding
	api.GET("/frameflow/bonding", a.handleFrameFlowBonding)

	// MediaMTX
	api.POST("/mediamtx/:action", a.handleMediaMTX)

	// GPS Tracker
	api.POST("/gps/:action", a.handleGPS)

	// Cameraman
	api.POST("/cameraman/:action", a.handleCameraman)
}

func (a *API) handleFrameFlowClient(c *gin.Context) {
	action := c.Param("action")
	if action != "start" && action != "stop" && action != "status" && action != "reset" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	out, err := a.executor.Run("VLX_FrameFlow", "client", action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}

func (a *API) handleFrameFlowAP(c *gin.Context) {
	action := c.Param("action")
	if action != "start" && action != "stop" && action != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	out, err := a.executor.Run("VLX_FrameFlow", "AP", action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}

func (a *API) handleFrameFlowBonding(c *gin.Context) {
	out, err := a.executor.Run("VLX_FrameFlow", "bonding")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}

func (a *API) handleMediaMTX(c *gin.Context) {
	action := c.Param("action")
	if action != "start" && action != "stop" && action != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	out, err := a.executor.Run("VLX_mediamtx", action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}

func (a *API) handleGPS(c *gin.Context) {
	action := c.Param("action")
	if action != "start" && action != "stop" && action != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	out, err := a.executor.Run("VLX_gps_tracker", action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}

type CameramanRequest struct {
	Device string `json:"device"` // e.g. V0A1
}

func (a *API) handleCameraman(c *gin.Context) {
	action := c.Param("action")
	if action != "start" && action != "stop" && action != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	var req CameramanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If device is not in JSON, try query param for compatibility
		req.Device = c.Query("device")
	}

	var args []string
	if req.Device != "" {
		args = append(args, req.Device)
	}
	args = append(args, action)

	out, err := a.executor.Run("VLX_cameraman", args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": out})
}
