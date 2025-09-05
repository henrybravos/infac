package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"infac/internal/config"
	"infac/internal/handlers"
	"infac/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services (now using sunatlib internally)
	documentService := services.NewDocumentService(&cfg.Issuer)

	// Initialize handlers
	documentHandler := handlers.NewDocumentHandler(documentService)

	// Setup Gin router
	r := gin.Default()
	
	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "infac-api",
		})
	})

	// Register routes
	documentHandler.RegisterRoutes(r)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	log.Fatal(r.Run(addr))
}