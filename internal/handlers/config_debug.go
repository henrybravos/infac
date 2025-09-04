package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"infac/internal/models"
)

type ConfigDebugHandler struct {
	issuer *models.Company
}

func NewConfigDebugHandler(issuer *models.Company) *ConfigDebugHandler {
	return &ConfigDebugHandler{
		issuer: issuer,
	}
}

func (h *ConfigDebugHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"issuer": h.issuer,
	})
}

func (h *ConfigDebugHandler) RegisterRoutes(r *gin.Engine) {
	debug := r.Group("/config")
	{
		debug.GET("/issuer", h.GetConfig)
	}
}