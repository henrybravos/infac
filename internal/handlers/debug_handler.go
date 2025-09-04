package handlers

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"infac/internal/models"
	"infac/pkg/ubl"
)

type DebugHandler struct {
	issuer *models.Company
}

func NewDebugHandler(issuer *models.Company) *DebugHandler {
	return &DebugHandler{
		issuer: issuer,
	}
}

func (h *DebugHandler) GenerateXML(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var xmlContent []byte
	var err error

	switch doc.Type {
	case models.DocumentTypeFactura, models.DocumentTypeBoleta:
		invoice, invoiceErr := ubl.GenerateInvoiceXML(&doc, h.issuer)
		if invoiceErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": invoiceErr.Error()})
			return
		}
		xmlContent, err = xml.MarshalIndent(invoice, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case models.DocumentTypeNotaCredito:
		creditNote, err := ubl.GenerateCreditNoteXML(&doc, h.issuer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		xmlContent, err = xml.MarshalIndent(creditNote, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case models.DocumentTypeNotaDebito:
		debitNote, err := ubl.GenerateDebitNoteXML(&doc, h.issuer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		xmlContent, err = xml.MarshalIndent(debitNote, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported document type"})
		return
	}

	// Add XML declaration
	xmlWithDeclaration := append([]byte(xml.Header), xmlContent...)

	c.Header("Content-Type", "application/xml")
	c.Data(http.StatusOK, "application/xml", xmlWithDeclaration)
}

func (h *DebugHandler) RegisterRoutes(r *gin.Engine) {
	debug := r.Group("/debug")
	{
		debug.POST("/xml", h.GenerateXML)
	}
}