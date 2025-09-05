package models

type CreateDocumentRequest struct {
	Type         DocumentType `json:"type" binding:"required"`
	Serie        string       `json:"serie" binding:"required"`
	Number       string       `json:"number" binding:"required"`
	IssueDate    string       `json:"issue_date" binding:"required"`
	DueDate      string       `json:"due_date,omitempty"`
	CurrencyCode string       `json:"currency_code" binding:"required"`
	
	Customer Company `json:"customer" binding:"required"`
	
	Lines []CreateDocumentLineRequest `json:"lines" binding:"required,min=1"`
	
	PaymentTerms *PaymentTerms `json:"payment_terms" binding:"required"`
	
	// Para notas de crédito/débito
	RelatedDocuments []RelatedDocument `json:"related_documents,omitempty"`
}

type CreateDocumentLineRequest struct {
	Quantity         float64 `json:"quantity" binding:"required,gt=0"`
	UnitCode         string  `json:"unit_code" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	UnitPrice        float64 `json:"unit_price" binding:"required"`
	
	Taxes []Tax `json:"taxes" binding:"required"`
	
	ProductCode string `json:"product_code,omitempty"`
}

type VoidDocumentRequest struct {
	DocumentType DocumentType `json:"document_type" binding:"required"`
	Serie        string       `json:"serie" binding:"required"`
	Number       string       `json:"number" binding:"required"`
	VoidDate     string       `json:"void_date" binding:"required"`
	Reason       string       `json:"reason" binding:"required"`
}