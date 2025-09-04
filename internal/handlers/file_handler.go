package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	
	"infac/internal/models"
	"infac/pkg/ubl"
)

type FileHandler struct {
	issuer *models.Company
}

func NewFileHandler(issuer *models.Company) *FileHandler {
	return &FileHandler{
		issuer: issuer,
	}
}

func (h *FileHandler) GenerateFiles(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set issuer data properly
	doc.Issuer = *h.issuer

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
		creditNote, creditErr := ubl.GenerateCreditNoteXML(&doc, h.issuer)
		if creditErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": creditErr.Error()})
			return
		}
		xmlContent, err = xml.MarshalIndent(creditNote, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case models.DocumentTypeNotaDebito:
		debitNote, debitErr := ubl.GenerateDebitNoteXML(&doc, h.issuer)
		if debitErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": debitErr.Error()})
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

	// Generate file names
	xmlFileName := fmt.Sprintf("%s-%s-%s.xml", h.issuer.DocumentNumber, doc.Serie, doc.Number)
	zipFileName := fmt.Sprintf("%s-%s-%s.zip", h.issuer.DocumentNumber, doc.Serie, doc.Number)

	// Save XML file
	xmlPath := filepath.Join("storage", "xml", xmlFileName)
	err = os.WriteFile(xmlPath, xmlWithDeclaration, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save XML: %v", err)})
		return
	}

	// Create ZIP file
	zipContent, err := createZipFile(xmlFileName, xmlWithDeclaration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create ZIP: %v", err)})
		return
	}

	// Save ZIP file
	zipPath := filepath.Join("storage", "zip", zipFileName)
	err = os.WriteFile(zipPath, zipContent, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save ZIP: %v", err)})
		return
	}

	// Encode to base64
	base64Content := base64.StdEncoding.EncodeToString(zipContent)

	result := gin.H{
		"document":      doc,
		"xml":          string(xmlWithDeclaration),
		"zip_base64":   base64Content,
		"files": gin.H{
			"xml_filename": xmlFileName,
			"xml_path":     xmlPath,
			"zip_filename": zipFileName,
			"zip_path":     zipPath,
		},
		"message": "Files saved successfully",
	}

	c.JSON(http.StatusOK, result)
}

func createZipFile(fileName string, content []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	f, err := w.Create(fileName)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(content)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *FileHandler) RegisterRoutes(r *gin.Engine) {
	files := r.Group("/files")
	{
		files.POST("/generate", h.GenerateFiles)
	}
}