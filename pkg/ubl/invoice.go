package ubl

import (
	"encoding/xml"
	"fmt"
	
	"infac/internal/models"
)

type Invoice struct {
	XMLName xml.Name `xml:"Invoice"`
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
	DueDate       string        `xml:"cbc:DueDate,omitempty"`
	InvoiceTypeCode InvoiceTypeCode `xml:"cbc:InvoiceTypeCode"`
	Note          []Note        `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode string `xml:"cbc:DocumentCurrencyCode"`
	
	Signature            []Signature          `xml:"cac:Signature"`
	AccountingSupplierParty AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty `xml:"cac:AccountingCustomerParty"`
	PaymentTerms         []PaymentTerms       `xml:"cac:PaymentTerms,omitempty"`
	TaxTotal             []TaxTotal           `xml:"cac:TaxTotal"`
	LegalMonetaryTotal   LegalMonetaryTotal   `xml:"cac:LegalMonetaryTotal"`
	InvoiceLine          []InvoiceLine        `xml:"cac:InvoiceLine"`
}

type UBLExtensions struct {
	UBLExtension UBLExtension `xml:"ext:UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"ext:ExtensionContent"`
}

type ExtensionContent struct {
	// Este contenido será añadido por el proveedor de firma digital
	Content string `xml:",innerxml"`
}

type InvoiceTypeCode struct {
	ListAgencyName string `xml:"listAgencyName,attr"`
	ListName       string `xml:"listName,attr"`
	ListURI        string `xml:"listURI,attr"`
	Value          string `xml:",chardata"`
}

type Note struct {
	LanguageLocaleID string `xml:"languageLocaleID,attr,omitempty"`
	Value            string `xml:",chardata"`
}

type Signature struct {
	ID                  string              `xml:"cbc:ID"`
	SignatoryParty      SignatoryParty      `xml:"cac:SignatoryParty"`
	DigitalSignatureAttachment DigitalSignatureAttachment `xml:"cac:DigitalSignatureAttachment"`
}

type SignatoryParty struct {
	PartyIdentification PartyIdentification `xml:"cac:PartyIdentification"`
	PartyName           PartyName           `xml:"cac:PartyName"`
}

type DigitalSignatureAttachment struct {
	ExternalReference ExternalReference `xml:"cac:ExternalReference"`
}

type ExternalReference struct {
	URI string `xml:"cbc:URI"`
}

type AccountingSupplierParty struct {
	Party Party `xml:"cac:Party"`
}

type AccountingCustomerParty struct {
	Party Party `xml:"cac:Party"`
}

type Party struct {
	PartyIdentification []PartyIdentification `xml:"cac:PartyIdentification"`
	PartyName           PartyName             `xml:"cac:PartyName"`
	PartyTaxScheme      []PartyTaxScheme      `xml:"cac:PartyTaxScheme,omitempty"`
	PartyLegalEntity    []PartyLegalEntity    `xml:"cac:PartyLegalEntity,omitempty"`
	Contact             Contact               `xml:"cac:Contact,omitempty"`
}

type PartyIdentification struct {
	ID IDType `xml:"cbc:ID"`
}

type IDType struct {
	SchemeID   string `xml:"schemeID,attr,omitempty"`
	SchemeName string `xml:"schemeName,attr,omitempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omitempty"`
	SchemeURI  string `xml:"schemeURI,attr,omitempty"`
	Value      string `xml:",chardata"`
}

type PartyName struct {
	Name string `xml:"cbc:Name"`
}

type PartyTaxScheme struct {
	RegistrationName string    `xml:"cbc:RegistrationName"`
	CompanyID        IDType    `xml:"cbc:CompanyID"`
	TaxScheme        TaxScheme `xml:"cac:TaxScheme"`
}

type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
}

type Contact struct {
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
	Telephone      string `xml:"cbc:Telephone,omitempty"`
}

type TaxScheme struct {
	ID   string `xml:"cbc:ID"`
	Name string `xml:"cbc:Name"`
}

type PaymentTerms struct {
	ID               string           `xml:"cbc:ID"`
	PaymentMeansID   string           `xml:"cbc:PaymentMeansID"`
	Amount           MonetaryAmount   `xml:"cbc:Amount,omitempty"`
	PaymentDueDate   string           `xml:"cbc:PaymentDueDate,omitempty"`
}

type TaxTotal struct {
	TaxAmount   MonetaryAmount `xml:"cbc:TaxAmount"`
	TaxSubtotal []TaxSubtotal  `xml:"cac:TaxSubtotal"`
}

type TaxSubtotal struct {
	TaxableAmount MonetaryAmount `xml:"cbc:TaxableAmount"`
	TaxAmount     MonetaryAmount `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategory    `xml:"cac:TaxCategory"`
}

type TaxCategory struct {
	ID                string    `xml:"cbc:ID"`
	Percent           float64   `xml:"cbc:Percent,omitempty"`
	TaxExemptionReasonCode string `xml:"cbc:TaxExemptionReasonCode,omitempty"`
	TaxScheme         TaxScheme `xml:"cac:TaxScheme"`
}

type MonetaryAmount struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Value      float64 `xml:",chardata"`
}

type LegalMonetaryTotal struct {
	LineExtensionAmount MonetaryAmount `xml:"cbc:LineExtensionAmount,omitempty"`
	TaxInclusiveAmount  MonetaryAmount `xml:"cbc:TaxInclusiveAmount,omitempty"`
	TaxExclusiveAmount  MonetaryAmount `xml:"cbc:TaxExclusiveAmount,omitempty"`
	PayableAmount       MonetaryAmount `xml:"cbc:PayableAmount"`
}

type InvoiceLine struct {
	ID                  string             `xml:"cbc:ID"`
	InvoicedQuantity    InvoicedQuantity   `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount MonetaryAmount     `xml:"cbc:LineExtensionAmount"`
	PricingReference    PricingReference   `xml:"cac:PricingReference,omitempty"`
	TaxTotal            []TaxTotal         `xml:"cac:TaxTotal,omitempty"`
	Item                Item               `xml:"cac:Item"`
	Price               Price              `xml:"cac:Price"`
}

type InvoicedQuantity struct {
	UnitCode string  `xml:"unitCode,attr"`
	Value    float64 `xml:",chardata"`
}

type PricingReference struct {
	AlternativeConditionPrice []AlternativeConditionPrice `xml:"cac:AlternativeConditionPrice"`
}

type AlternativeConditionPrice struct {
	PriceAmount   MonetaryAmount `xml:"cbc:PriceAmount"`
	PriceTypeCode string         `xml:"cbc:PriceTypeCode"`
}

type Item struct {
	Description                     []string                        `xml:"cbc:Description"`
	SellersItemIdentification       SellersItemIdentification       `xml:"cac:SellersItemIdentification,omitempty"`
	CommodityClassification         []CommodityClassification       `xml:"cac:CommodityClassification,omitempty"`
}

type SellersItemIdentification struct {
	ID string `xml:"cbc:ID"`
}

type CommodityClassification struct {
	ItemClassificationCode ItemClassificationCode `xml:"cbc:ItemClassificationCode"`
}

type ItemClassificationCode struct {
	ListID   string `xml:"listID,attr,omitempty"`
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	Value    string `xml:",chardata"`
}

type Price struct {
	PriceAmount MonetaryAmount `xml:"cbc:PriceAmount"`
}

func GenerateInvoiceXML(doc *models.Document, issuer *models.Company) (*Invoice, error) {
	invoice := &Invoice{
		Xmlns:    "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		XmlnsCac: "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsCbc: "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		XmlnsExt: "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2",
		
		UBLVersionID:    "2.1",
		CustomizationID: "2.0",
		ID:              fmt.Sprintf("%s-%s", doc.Serie, doc.Number),
		IssueDate:       doc.IssueDate.Format("2006-01-02"),
		IssueTime:       doc.IssueDate.Format("15:04:05"),
		DocumentCurrencyCode: doc.CurrencyCode,
		
		InvoiceTypeCode: InvoiceTypeCode{
			ListAgencyName: "PE:SUNAT",
			ListName:       "Tipo de Documento",
			ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo01",
			Value:          string(doc.Type),
		},
	}
	
	if doc.DueDate != nil {
		invoice.DueDate = doc.DueDate.Format("2006-01-02")
	}
	
	// Signature
	invoice.Signature = []Signature{
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
	
	// Supplier Party
	invoice.AccountingSupplierParty = AccountingSupplierParty{
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
	
	if issuer.Email != "" || issuer.Phone != "" {
		invoice.AccountingSupplierParty.Party.Contact = Contact{
			ElectronicMail: issuer.Email,
			Telephone:      issuer.Phone,
		}
	}
	
	// Customer Party
	invoice.AccountingCustomerParty = AccountingCustomerParty{
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
	
	// Payment Terms
	if doc.PaymentTerms != nil {
		invoice.PaymentTerms = []PaymentTerms{
			{
				ID:             "FormaPago",
				PaymentMeansID: doc.PaymentTerms.PaymentMeansCode,
				Amount: MonetaryAmount{
					CurrencyID: doc.CurrencyCode,
					Value:      doc.PaymentTerms.Amount,
				},
				PaymentDueDate: doc.PaymentTerms.DueDate.Format("2006-01-02"),
			},
		}
	}
	
	// Tax Total
	taxTotal := TaxTotal{
		TaxAmount: MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.TotalTaxes,
		},
	}
	
	// Agrupar impuestos por tipo
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
	
	invoice.TaxTotal = []TaxTotal{taxTotal}
	
	// Legal Monetary Total
	invoice.LegalMonetaryTotal = LegalMonetaryTotal{
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
	
	// Invoice Lines
	for i, line := range doc.Lines {
		invoiceLine := InvoiceLine{
			ID: fmt.Sprintf("%d", i+1),
			InvoicedQuantity: InvoicedQuantity{
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
			invoiceLine.Item.SellersItemIdentification = SellersItemIdentification{
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
			
			invoiceLine.TaxTotal = []TaxTotal{lineTaxTotal}
		}
		
		// Pricing reference para mostrar precio con impuestos
		invoiceLine.PricingReference = PricingReference{
			AlternativeConditionPrice: []AlternativeConditionPrice{
				{
					PriceAmount: MonetaryAmount{
						CurrencyID: doc.CurrencyCode,
						Value:      line.UnitPrice * (1 + getTotalTaxRate(line.Taxes)/100),
					},
					PriceTypeCode: "01", // Precio unitario (incluye el IGV)
				},
			},
		}
		
		invoice.InvoiceLine = append(invoice.InvoiceLine, invoiceLine)
	}
	
	return invoice, nil
}

func getDocumentTypeScheme(docType string) string {
	switch docType {
	case "6": // RUC
		return "6"
	case "1": // DNI
		return "1"
	default:
		return docType
	}
}

func getTaxSchemeID(taxType models.TaxType) string {
	switch taxType {
	case models.TaxTypeIGV:
		return "1000"
	case models.TaxTypeISC:
		return "2000"
	case models.TaxTypeICBP:
		return "7152"
	default:
		return "9999"
	}
}

func getTotalTaxRate(taxes []models.Tax) float64 {
	total := 0.0
	for _, tax := range taxes {
		total += tax.Rate
	}
	return total
}