package services

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"infac/internal/models"
	"infac/pkg/ubl"

	"github.com/henrybravos/sunatlib"
)

type DocumentService struct {
	issuer      *models.Company
	sunatClient *sunatlib.SUNATClient
}

func NewDocumentService(issuer *models.Company) *DocumentService {
	// Create SUNAT client
	sunatClient := sunatlib.NewSUNATClient(
		issuer.DocumentNumber, // RUC
		"MODDATOS",            // SOL username
		"moddatos",            // SOL password
		"https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService", // Beta endpoint
	)

	// Set certificate from PFX file
	err := sunatClient.SetCertificateFromPFX(
		"pkg/signature/certificate_fixed.pfx",
		"20612790168NEOFORCE",
		"pkg/signature/temp",
	)
	if err != nil {
		fmt.Printf("Warning: Failed to set certificate: %v\n", err)
		fmt.Printf("Continuing without certificate - signatures will not be valid\n")
	} else {
		fmt.Printf("Successfully initialized SUNAT client with certificate\n")
	}

	return &DocumentService{
		issuer:      issuer,
		sunatClient: sunatClient,
	}
}

func (s *DocumentService) CreateDocument(req *models.CreateDocumentRequest) (*models.Document, error) {
	doc := &models.Document{
		ID:               fmt.Sprintf("%s-%s", req.Serie, req.Number),
		Serie:            req.Serie,
		Number:           req.Number,
		Type:             req.Type,
		CurrencyCode:     req.CurrencyCode,
		Issuer:           *s.issuer,
		Customer:         req.Customer,
		PaymentTerms:     req.PaymentTerms,
		RelatedDocuments: req.RelatedDocuments,
		Status:           models.StatusDraft,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Parse dates
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issue date format: %w", err)
	}
	doc.IssueDate = issueDate

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format: %w", err)
		}
		doc.DueDate = &dueDate
	}

	// Process lines
	var subTotal float64
	var totalTaxes float64

	for i, lineReq := range req.Lines {
		line := models.DocumentLine{
			ID:          fmt.Sprintf("%d", i+1),
			Quantity:    lineReq.Quantity,
			UnitCode:    lineReq.UnitCode,
			Description: lineReq.Description,
			UnitPrice:   lineReq.UnitPrice,
			TotalPrice:  lineReq.Quantity * lineReq.UnitPrice,
			ProductCode: lineReq.ProductCode,
			Taxes:       lineReq.Taxes,
		}

		// Calculate taxable amount (base for taxes)
		line.TaxableAmount = line.TotalPrice

		// Calculate taxes
		for j, tax := range line.Taxes {
			line.Taxes[j].Amount = line.TaxableAmount * tax.Rate / 100
			totalTaxes += line.Taxes[j].Amount
		}

		subTotal += line.TotalPrice
		doc.Lines = append(doc.Lines, line)
	}

	doc.SubTotal = subTotal
	doc.TotalTaxes = totalTaxes
	doc.TotalAmount = subTotal + totalTaxes

	return doc, nil
}
func (s *DocumentService) SendDocument(doc *models.Document) error {
	// 1. Generar XML según tipo
	var xmlContent []byte
	var err error

	switch doc.Type {
	case models.DocumentTypeFactura, models.DocumentTypeBoleta:
		invoice, err := ubl.GenerateInvoiceXML(doc, &doc.Issuer)
		if err != nil {
			return fmt.Errorf("failed to generate invoice XML: %w", err)
		}
		xmlContent, err = xml.MarshalIndent(invoice, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal invoice XML: %w", err)
		}
	case models.DocumentTypeNotaCredito:
		creditNote, err := ubl.GenerateCreditNoteXML(doc, &doc.Issuer)
		if err != nil {
			return fmt.Errorf("failed to generate credit note XML: %w", err)
		}
		xmlContent, err = xml.MarshalIndent(creditNote, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal credit note XML: %w", err)
		}
	case models.DocumentTypeNotaDebito:
		debitNote, err := ubl.GenerateDebitNoteXML(doc, &doc.Issuer)
		if err != nil {
			return fmt.Errorf("failed to generate debit note XML: %w", err)
		}
		xmlContent, err = xml.MarshalIndent(debitNote, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal debit note XML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported document type: %s", doc.Type)
	}

	// 2. Agregar declaración XML
	xmlWithDeclaration := append([]byte(xml.Header), xmlContent...)

	// 3. Firmar XML (antes de enviar)
	signedXML, err := s.sunatClient.SignXML(xmlWithDeclaration)
	if err != nil {
		return fmt.Errorf("failed to sign XML: %w", err)
	}

	// 4. Guardar el XML firmado para depuración
	xmlFileName := fmt.Sprintf("%s-%s-%s-%s-signed.xml", s.issuer.DocumentNumber, string(doc.Type), doc.Serie, doc.Number)
	xmlPath := filepath.Join("storage", "xml", xmlFileName)
	if err := os.WriteFile(xmlPath, signedXML, 0644); err != nil {
		fmt.Printf("Warning: Failed to save signed XML: %v\n", err)
	}

	// 5. Generar ID de documento
	documentID := fmt.Sprintf("%s-%s", doc.Serie, doc.Number)

	// 6. Enviar a SUNAT
	response, err := s.sunatClient.SendToSUNAT(signedXML, string(doc.Type), documentID)
	if err != nil {
		doc.Status = models.StatusRejected
		return fmt.Errorf("failed to send document to SUNAT: %w", err)
	}

	// 7. Manejar respuesta
	if response.Success {
		doc.Status = models.StatusAccepted
		doc.CDR = &models.CDR{
			ResponseCode: "0",
			Description:  "Accepted",
		}

		if len(response.ApplicationResponse) > 0 {
			cdrPath := filepath.Join("storage", "cdr", fmt.Sprintf("R-%s-%s-%s-%s.zip",
				s.issuer.DocumentNumber, string(doc.Type), doc.Serie, doc.Number))
			err = response.SaveApplicationResponse(cdrPath)
			if err != nil {
				fmt.Printf("Warning: Failed to save CDR: %v\n", err)
			}
		}
	} else {
		doc.Status = models.StatusRejected
		fmt.Printf("SUNAT response: %+v\n", response)
		
		// Try to get the most specific error message available
		var errorMsg string
		if response.Message != "" {
			errorMsg = response.Message
		} else if response.Error != nil {
			errorMsg = response.Error.Error()
		} else {
			errorMsg = "unknown error"
		}
		
		return fmt.Errorf("document rejected by SUNAT: %s", errorMsg)
	}

	doc.UpdatedAt = time.Now()
	return nil
}

func (s *DocumentService) CheckStatus(ticket string) (*models.CDR, error) {
	// For now, return a placeholder CDR
	// TODO: Implement status checking with sunatlib if available
	cdr := &models.CDR{
		ResponseCode: "0",
		Description:  "Status check not implemented yet",
		Notes:        fmt.Sprintf("Ticket: %s", ticket),
	}

	return cdr, nil
}

func (s *DocumentService) VoidDocument(req *models.VoidDocumentRequest) error {
	// Generate void communication XML
	// This would involve creating a summary of voided documents
	// Implementation depends on SUNAT's specific requirements
	return fmt.Errorf("void document not implemented yet")
}

func (s *DocumentService) createZipFile(fileName string, content []byte) ([]byte, error) {
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
