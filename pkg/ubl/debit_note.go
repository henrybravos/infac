package ubl

import (
	"encoding/xml"
	"fmt"
	
	"infac/internal/models"
)

type DebitNote struct {
	XMLName xml.Name `xml:"DebitNote"`
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
	RequestedMonetaryTotal RequestedMonetaryTotal `xml:"cac:RequestedMonetaryTotal"`
	DebitNoteLine        []DebitNoteLine      `xml:"cac:DebitNoteLine"`
}

type RequestedMonetaryTotal struct {
	LineExtensionAmount MonetaryAmount `xml:"cbc:LineExtensionAmount,omitempty"`
	TaxInclusiveAmount  MonetaryAmount `xml:"cbc:TaxInclusiveAmount,omitempty"`
	TaxExclusiveAmount  MonetaryAmount `xml:"cbc:TaxExclusiveAmount,omitempty"`
	PayableAmount       MonetaryAmount `xml:"cbc:PayableAmount"`
}

type DebitNoteLine struct {
	ID                  string             `xml:"cbc:ID"`
	DebitedQuantity     InvoicedQuantity   `xml:"cbc:DebitedQuantity"`
	LineExtensionAmount MonetaryAmount     `xml:"cbc:LineExtensionAmount"`
	PricingReference    PricingReference   `xml:"cac:PricingReference,omitempty"`
	TaxTotal            []TaxTotal         `xml:"cac:TaxTotal,omitempty"`
	Item                Item               `xml:"cac:Item"`
	Price               Price              `xml:"cac:Price"`
}

func GenerateDebitNoteXML(doc *models.Document, issuer *models.Company) (*DebitNote, error) {
	debitNote := &DebitNote{
		Xmlns:    "urn:oasis:names:specification:ubl:schema:xsd:DebitNote-2",
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
	
	// Discrepancy Response - motivo de la nota de débito
	if len(doc.RelatedDocuments) > 0 {
		debitNote.DiscrepancyResponse = []DiscrepancyResponse{
			{
				ReferenceID:  fmt.Sprintf("%s-%s", doc.RelatedDocuments[0].Serie, doc.RelatedDocuments[0].Number),
				ResponseCode: "02", // Aumento en el valor
				Description:  "Aumento en el valor",
			},
		}
		
		// Billing Reference - documento relacionado
		debitNote.BillingReference = []BillingReference{
			{
				InvoiceDocumentReference: InvoiceDocumentReference{
					ID: fmt.Sprintf("%s-%s", doc.RelatedDocuments[0].Serie, doc.RelatedDocuments[0].Number),
					DocumentTypeCode: string(doc.RelatedDocuments[0].DocumentType),
				},
			},
		}
	}
	
	// Signature (similar a Invoice)
	debitNote.Signature = []Signature{
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
	debitNote.AccountingSupplierParty = AccountingSupplierParty{
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
	
	debitNote.AccountingCustomerParty = AccountingCustomerParty{
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
	
	debitNote.TaxTotal = []TaxTotal{taxTotal}
	
	// Requested Monetary Total
	debitNote.RequestedMonetaryTotal = RequestedMonetaryTotal{
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
	
	// Debit Note Lines
	for i, line := range doc.Lines {
		debitNoteLine := DebitNoteLine{
			ID: fmt.Sprintf("%d", i+1),
			DebitedQuantity: InvoicedQuantity{
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
			debitNoteLine.Item.SellersItemIdentification = SellersItemIdentification{
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
			
			debitNoteLine.TaxTotal = []TaxTotal{lineTaxTotal}
		}
		
		// Pricing reference
		debitNoteLine.PricingReference = PricingReference{
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
		
		debitNote.DebitNoteLine = append(debitNote.DebitNoteLine, debitNoteLine)
	}
	
	return debitNote, nil
}