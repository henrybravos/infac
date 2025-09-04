package models

import "time"

type DocumentType string

const (
	DocumentTypeFactura        DocumentType = "01" // Factura
	DocumentTypeBoleta         DocumentType = "03" // Boleta de Venta
	DocumentTypeNotaCredito    DocumentType = "07" // Nota de Crédito
	DocumentTypeNotaDebito     DocumentType = "08" // Nota de Débito
)

type Document struct {
	ID           string       `json:"id"`
	Serie        string       `json:"serie"`
	Number       string       `json:"number"`
	Type         DocumentType `json:"type"`
	IssueDate    time.Time    `json:"issue_date"`
	DueDate      *time.Time   `json:"due_date,omitempty"`
	CurrencyCode string       `json:"currency_code"`
	
	Issuer   Company `json:"issuer"`
	Customer Company `json:"customer"`
	
	Lines []DocumentLine `json:"lines"`
	
	SubTotal      float64 `json:"sub_total"`
	TotalTaxes    float64 `json:"total_taxes"`
	TotalAmount   float64 `json:"total_amount"`
	
	PaymentTerms *PaymentTerms `json:"payment_terms,omitempty"`
	
	// Para notas de crédito/débito
	RelatedDocuments []RelatedDocument `json:"related_documents,omitempty"`
	
	// Estado del documento
	Status       DocumentStatus `json:"status"`
	SUNATStatus  string        `json:"sunat_status,omitempty"`
	CDR          *CDR          `json:"cdr,omitempty"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentStatus string

const (
	StatusDraft     DocumentStatus = "draft"
	StatusPending   DocumentStatus = "pending"
	StatusSent      DocumentStatus = "sent"
	StatusAccepted  DocumentStatus = "accepted"
	StatusRejected  DocumentStatus = "rejected"
	StatusCancelled DocumentStatus = "cancelled"
)

type Company struct {
	DocumentType   string `json:"document_type" mapstructure:"document_type"`
	DocumentNumber string `json:"document_number" mapstructure:"document_number"`
	Name           string `json:"name" mapstructure:"name"`
	TradeName      string `json:"trade_name,omitempty" mapstructure:"trade_name"`
	Address        string `json:"address" mapstructure:"address"`
	District       string `json:"district" mapstructure:"district"`
	Province       string `json:"province" mapstructure:"province"`
	Department     string `json:"department" mapstructure:"department"`
	Country        string `json:"country" mapstructure:"country"`
	Email          string `json:"email,omitempty" mapstructure:"email"`
	Phone          string `json:"phone,omitempty" mapstructure:"phone"`
}

type DocumentLine struct {
	ID               string  `json:"id"`
	Quantity         float64 `json:"quantity"`
	UnitCode         string  `json:"unit_code"`
	Description      string  `json:"description"`
	UnitPrice        float64 `json:"unit_price"`
	TotalPrice       float64 `json:"total_price"`
	TaxableAmount    float64 `json:"taxable_amount"`
	
	Taxes []Tax `json:"taxes"`
	
	ProductCode string `json:"product_code,omitempty"`
}

type Tax struct {
	Type   TaxType `json:"type"`
	Code   string  `json:"code"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
}

type TaxType string

const (
	TaxTypeIGV  TaxType = "IGV"  // Impuesto General a las Ventas
	TaxTypeISC  TaxType = "ISC"  // Impuesto Selectivo al Consumo
	TaxTypeICBP TaxType = "ICBP" // Impuesto a las Bolsas de Plástico
)

type PaymentTerms struct {
	PaymentMeansCode string    `json:"payment_means_code"`
	DueDate          time.Time `json:"due_date"`
	Amount           float64   `json:"amount"`
}

type RelatedDocument struct {
	DocumentType DocumentType `json:"document_type"`
	Serie        string       `json:"serie"`
	Number       string       `json:"number"`
}

type CDR struct {
	ResponseCode string `json:"response_code"`
	Description  string `json:"description"`
	Notes        string `json:"notes,omitempty"`
}