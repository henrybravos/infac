package ubl

import (
	"encoding/xml"
	"fmt"
	
	"infac/internal/models"
)

type CreditNote struct {
	XMLName xml.Name `xml:"CreditNote"`
	Xmlns   string   `xml:"xmlns,attr"`
	XmlnsCac string  `xml:"xmlns:cac,attr"`
	XmlnsCbc string  `xml:"xmlns:cbc,attr"`
	XmlnsExt string  `xml:"xmlns:ext,attr"`
	
	UBLExtensions UBLExtensions `xml:"ext:UBLExtensions"`
	UBLVersionID  string        `xml:"cbc:UBLVersionID"`
	CustomizationID string      `xml:"cbc:CustomizationID"`
	ID            string        `xml:"cbc:ID"`
	IssueDate     string        `xml:"cbc:IssueDate"`
	IssueTime     string        `xml:"cbc:IssueTime"`
	Note          []Note        `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode string `xml:"cbc:DocumentCurrencyCode"`
	
	DiscrepancyResponse []DiscrepancyResponse `xml:"cac:DiscrepancyResponse"`
	BillingReference    []BillingReference    `xml:"cac:BillingReference"`
	
	Signature            []Signature          `xml:"cac:Signature"`
	AccountingSupplierParty AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty `xml:"cac:AccountingCustomerParty"`
	TaxTotal             []TaxTotal           `xml:"cac:TaxTotal"`
	LegalMonetaryTotal   LegalMonetaryTotal   `xml:"cac:LegalMonetaryTotal"`
	CreditNoteLine       []CreditNoteLine     `xml:"cac:CreditNoteLine"`
}

type DiscrepancyResponse struct {
	ReferenceID   string `xml:"cbc:ReferenceID"`
	ResponseCode  string `xml:"cbc:ResponseCode"`
	Description   string `xml:"cbc:Description"`
}

type BillingReference struct {
	InvoiceDocumentReference InvoiceDocumentReference `xml:"cac:InvoiceDocumentReference"`
}

type InvoiceDocumentReference struct {
	ID           string `xml:"cbc:ID"`
	DocumentTypeCode string `xml:"cbc:DocumentTypeCode"`
}

type CreditNoteLine struct {
	ID                  string             `xml:"cbc:ID"`
	CreditedQuantity    InvoicedQuantity   `xml:"cbc:CreditedQuantity"`
	LineExtensionAmount MonetaryAmount     `xml:"cbc:LineExtensionAmount"`
	PricingReference    PricingReference   `xml:"cac:PricingReference,omitempty"`
	TaxTotal            []TaxTotal         `xml:"cac:TaxTotal,omitempty"`
	Item                Item               `xml:"cac:Item"`
	Price               Price              `xml:"cac:Price"`
}

func GenerateCreditNoteXML(doc *models.Document, issuer *models.Company) (*CreditNote, error) {
	creditNote := &CreditNote{
		Xmlns:    "urn:oasis:names:specification:ubl:schema:xsd:CreditNote-2",
		XmlnsCac: "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsCbc: "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		XmlnsExt: "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2",
		
		UBLVersionID:    "2.1",
		CustomizationID: "2.0",
		ID:              fmt.Sprintf("%s-%s", doc.Serie, doc.Number),
		IssueDate:       doc.IssueDate.Format("2006-01-02"),
		IssueTime:       doc.IssueDate.Format("15:04:05"),
		DocumentCurrencyCode: doc.CurrencyCode,
	}
	
	// Discrepancy Response - motivo de la nota de crédito
	if len(doc.RelatedDocuments) > 0 {
		creditNote.DiscrepancyResponse = []DiscrepancyResponse{
			{
				ReferenceID:  fmt.Sprintf("%s-%s", doc.RelatedDocuments[0].Serie, doc.RelatedDocuments[0].Number),
				ResponseCode: "01", // Anulación de la operación
				Description:  "Anulación de la operación",
			},
		}
		
		// Billing Reference - documento relacionado
		creditNote.BillingReference = []BillingReference{
			{
				InvoiceDocumentReference: InvoiceDocumentReference{
					ID: fmt.Sprintf("%s-%s", doc.RelatedDocuments[0].Serie, doc.RelatedDocuments[0].Number),
					DocumentTypeCode: string(doc.RelatedDocuments[0].DocumentType),
				},
			},
		}
	}
	
	// Signature (similar a Invoice)
	creditNote.Signature = []Signature{
		{
			ID: "IDSignST",
			SignatoryParty: SignatoryParty{
				PartyIdentification: PartyIdentification{
					ID: IDType{Value: issuer.DocumentNumber},
				},
				PartyName: PartyName{
					Name: issuer.Name,
				},
			},
			DigitalSignatureAttachment: DigitalSignatureAttachment{
				ExternalReference: ExternalReference{
					URI: "#SignatureST",
				},
			},
		},
	}
	
	// Supplier y Customer (similar a Invoice)
	creditNote.AccountingSupplierParty = AccountingSupplierParty{
		Party: Party{
			PartyIdentification: []PartyIdentification{
				{
					ID: IDType{
						SchemeID:         getDocumentTypeScheme(issuer.DocumentType),
						SchemeName:       "Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            issuer.DocumentNumber,
					},
				},
			},
			PartyName: PartyName{
				Name: issuer.TradeName,
			},
			PartyTaxScheme: []PartyTaxScheme{
				{
					RegistrationName: issuer.Name,
					CompanyID: IDType{
						SchemeID:         getDocumentTypeScheme(issuer.DocumentType),
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            issuer.DocumentNumber,
					},
					TaxScheme: TaxScheme{
						ID:   "9999",
						Name: "SUNAT",
					},
				},
			},
			PartyLegalEntity: []PartyLegalEntity{
				{
					RegistrationName: issuer.Name,
				},
			},
		},
	}
	
	creditNote.AccountingCustomerParty = AccountingCustomerParty{
		Party: Party{
			PartyIdentification: []PartyIdentification{
				{
					ID: IDType{
						SchemeID:         getDocumentTypeScheme(doc.Customer.DocumentType),
						SchemeName:       "Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            doc.Customer.DocumentNumber,
					},
				},
			},
			PartyLegalEntity: []PartyLegalEntity{
				{
					RegistrationName: doc.Customer.Name,
				},
			},
		},
	}
	
	// Tax Total (similar a Invoice)
	taxTotal := TaxTotal{
		TaxAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.TotalTaxes,
		},
	}
	
	taxGroups := make(map[string]*TaxSubtotal)
	for _, line := range doc.Lines {
		for _, tax := range line.Taxes {
			key := fmt.Sprintf("%s_%s", tax.Type, tax.Code)
			if taxGroups[key] == nil {
				taxGroups[key] = &TaxSubtotal{
					TaxableAmount: MonetaryAmount{CurrencyID: doc.CurrencyCode, Value: 0},
					TaxAmount:     MonetaryAmount{CurrencyID: doc.CurrencyCode, Value: 0},
					TaxCategory: TaxCategory{
						ID:      tax.Code,
						Percent: tax.Rate,
						TaxScheme: TaxScheme{
							ID:   getTaxSchemeID(tax.Type),
							Name: string(tax.Type),
						},
					},
				}
			}
			taxGroups[key].TaxableAmount.Value += line.TaxableAmount
			taxGroups[key].TaxAmount.Value += tax.Amount
		}
	}
	
	for _, taxSubtotal := range taxGroups {
		taxTotal.TaxSubtotal = append(taxTotal.TaxSubtotal, *taxSubtotal)
	}
	
	creditNote.TaxTotal = []TaxTotal{taxTotal}
	
	// Legal Monetary Total
	creditNote.LegalMonetaryTotal = LegalMonetaryTotal{
		LineExtensionAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.SubTotal,
		},
		TaxInclusiveAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.TotalAmount,
		},
		TaxExclusiveAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.SubTotal,
		},
		PayableAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.TotalAmount,
		},
	}
	
	// Credit Note Lines
	for i, line := range doc.Lines {
		creditNoteLine := CreditNoteLine{
			ID: fmt.Sprintf("%d", i+1),
			CreditedQuantity: InvoicedQuantity{
				UnitCode: line.UnitCode,
				Value:    line.Quantity,
			},
			LineExtensionAmount: MonetaryAmount{
				CurrencyID: doc.CurrencyCode,
				Value:      line.TotalPrice,
			},
			Item: Item{
				Description: []string{line.Description},
			},
			Price: Price{
				PriceAmount: MonetaryAmount{
					CurrencyID: doc.CurrencyCode,
					Value:      line.UnitPrice,
				},
			},
		}
		
		if line.ProductCode != "" {
			creditNoteLine.Item.SellersItemIdentification = SellersItemIdentification{
				ID: line.ProductCode,
			}
		}
		
		// Line taxes
		if len(line.Taxes) > 0 {
			lineTaxTotal := TaxTotal{
				TaxAmount: MonetaryAmount{CurrencyID: doc.CurrencyCode, Value: 0},
			}
			
			for _, tax := range line.Taxes {
				lineTaxTotal.TaxAmount.Value += tax.Amount
				lineTaxTotal.TaxSubtotal = append(lineTaxTotal.TaxSubtotal, TaxSubtotal{
					TaxableAmount: MonetaryAmount{
						CurrencyID: doc.CurrencyCode,
						Value:      line.TaxableAmount,
					},
					TaxAmount: MonetaryAmount{
						CurrencyID: doc.CurrencyCode,
						Value:      tax.Amount,
					},
					TaxCategory: TaxCategory{
						ID:      tax.Code,
						Percent: tax.Rate,
						TaxScheme: TaxScheme{
							ID:   getTaxSchemeID(tax.Type),
							Name: string(tax.Type),
						},
					},
				})
			}
			
			creditNoteLine.TaxTotal = []TaxTotal{lineTaxTotal}
		}
		
		// Pricing reference
		creditNoteLine.PricingReference = PricingReference{
			AlternativeConditionPrice: []AlternativeConditionPrice{
				{
					PriceAmount: MonetaryAmount{
						CurrencyID: doc.CurrencyCode,
						Value:      line.UnitPrice * (1 + getTotalTaxRate(line.Taxes)/100),
					},
					PriceTypeCode: "01",
				},
			},
		}
		
		creditNote.CreditNoteLine = append(creditNote.CreditNoteLine, creditNoteLine)
	}
	
	return creditNote, nil
}