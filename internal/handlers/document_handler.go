package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"infac/internal/models"
	"infac/internal/services"
)

type DocumentHandler struct {
	documentService *services.DocumentService
}

func NewDocumentHandler(documentService *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var req models.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.documentService.CreateDocument(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (h *DocumentHandler) SendDocument(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.documentService.SendDocument(&doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *DocumentHandler) CheckStatus(c *gin.Context) {
	ticket := c.Param("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket parameter is required"})
		return
	}

	cdr, err := h.documentService.CheckStatus(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cdr)
}

func (h *DocumentHandler) VoidDocument(c *gin.Context) {
	var req models.VoidDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.documentService.VoidDocument(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document voided successfully"})
}

func (h *DocumentHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		documents := api.Group("/documents")
		{
			documents.POST("", h.CreateDocument)
			documents.POST("/send", h.SendDocument)
			documents.POST("/void", h.VoidDocument)
			documents.GET("/status/:ticket", h.CheckStatus)
		}
	}
}