package ubl

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"

	"infac/internal/models"
)

type Invoice struct {
	XMLName   xml.Name `xml:"Invoice"`
	Xmlns     string   `xml:"xmlns,attr"`
	XmlnsCac  string   `xml:"xmlns:cac,attr"`
	XmlnsCbc  string   `xml:"xmlns:cbc,attr"`
	XmlnsCcts string   `xml:"xmlns:ccts,attr"`
	XmlnsDs   string   `xml:"xmlns:ds,attr"`
	XmlnsExt  string   `xml:"xmlns:ext,attr"`
	XmlnsQdt  string   `xml:"xmlns:qdt,attr"`
	XmlnsUdt  string   `xml:"xmlns:udt,attr"`
	XmlnsXsi  string   `xml:"xmlns:xsi,attr"`

	UBLExtensions        UBLExtensions   `xml:"ext:UBLExtensions"`
	UBLVersionID         string          `xml:"cbc:UBLVersionID"`
	CustomizationID      CustomizationID `xml:"cbc:CustomizationID"`
	ProfileID            ProfileID       `xml:"cbc:ProfileID"`
	ID                   string          `xml:"cbc:ID"`
	IssueDate            string          `xml:"cbc:IssueDate"`
	IssueTime            string          `xml:"cbc:IssueTime"`
	DueDate              string          `xml:"cbc:DueDate,omitempty"`
	InvoiceTypeCode      InvoiceTypeCode `xml:"cbc:InvoiceTypeCode"`
	Note                 []Note          `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode DocumentCurrencyCode `xml:"cbc:DocumentCurrencyCode"`
	LineCountNumeric     int             `xml:"cbc:LineCountNumeric"`

	Signature               []Signature             `xml:"cac:Signature"`
	AccountingSupplierParty AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty `xml:"cac:AccountingCustomerParty"`
	PaymentTerms            []PaymentTerms          `xml:"cac:PaymentTerms,omitempty"`
	TaxTotal                []TaxTotal              `xml:"cac:TaxTotal"`
	LegalMonetaryTotal      LegalMonetaryTotal      `xml:"cac:LegalMonetaryTotal"`
	InvoiceLine             []InvoiceLine           `xml:"cac:InvoiceLine"`
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

type CustomizationID struct {
	SchemeAgencyName string `xml:"schemeAgencyName,attr"`
	Value            string `xml:",chardata"`
}

type ProfileID struct {
	SchemeName       string `xml:"schemeName,attr"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr"`
	SchemeURI        string `xml:"schemeURI,attr"`
	Value            string `xml:",chardata"`
}

type InvoiceTypeCode struct {
	ListAgencyName string `xml:"listAgencyName,attr"`
	ListName       string `xml:"listName,attr"`
	ListURI        string `xml:"listURI,attr"`
	ListID         string `xml:"listID,attr,omitempty"`
	Name           string `xml:"name,attr,omitempty"`
	ListSchemeURI  string `xml:"listSchemeURI,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type DocumentCurrencyCode struct {
	ListID            string `xml:"listID,attr"`
	ListName          string `xml:"listName,attr"`
	ListAgencyName    string `xml:"listAgencyName,attr"`
	Value             string `xml:",chardata"`
}

type Note struct {
	LanguageLocaleID string `xml:"languageLocaleID,attr,omitempty"`
	Value            string `xml:",chardata"`
}

type Signature struct {
	ID                         string                     `xml:"cbc:ID"`
	SignatoryParty             SignatoryParty             `xml:"cac:SignatoryParty"`
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
	Contact             *Contact              `xml:"cac:Contact,omitempty"`
}

type PartyIdentification struct {
	ID IDType `xml:"cbc:ID"`
}

type IDType struct {
	SchemeID         string `xml:"schemeID,attr,omitempty"`
	SchemeName       string `xml:"schemeName,attr,omitempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omitempty"`
	SchemeURI        string `xml:"schemeURI,attr,omitempty"`
	ListID           string `xml:"listID,attr,omitempty"`
	ListName         string `xml:"listName,attr,omitempty"`
	ListAgencyName   string `xml:"listAgencyName,attr,omitempty"`
	Value            string `xml:",chardata"`
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
	RegistrationName    string               `xml:"cbc:RegistrationName"`
	RegistrationAddress *RegistrationAddress `xml:"cac:RegistrationAddress,omitempty"`
}

type RegistrationAddress struct {
	ID               IDType          `xml:"cbc:ID,omitempty"`
	AddressTypeCode  AddressTypeCode `xml:"cbc:AddressTypeCode,omitempty"`
	CityName         string          `xml:"cbc:CityName,omitempty"`
	CountrySubentity string          `xml:"cbc:CountrySubentity,omitempty"`
	District         string          `xml:"cbc:District,omitempty"`
	AddressLine      *AddressLine    `xml:"cac:AddressLine,omitempty"`
	Country          *Country        `xml:"cac:Country,omitempty"`
}

type AddressTypeCode struct {
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	ListName       string `xml:"listName,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type AddressLine struct {
	Line string `xml:"cbc:Line"`
}

type Country struct {
	IdentificationCode IDType `xml:"cbc:IdentificationCode"`
}

type Contact struct {
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
	Telephone      string `xml:"cbc:Telephone,omitempty"`
}

type TaxScheme struct {
	ID          IDType `xml:"cbc:ID"`
	Name        string `xml:"cbc:Name"`
	TaxTypeCode string `xml:"cbc:TaxTypeCode,omitempty"`
}

type PaymentTerms struct {
	ID             string         `xml:"cbc:ID"`
	PaymentMeansID string         `xml:"cbc:PaymentMeansID"`
	Amount         MonetaryAmount `xml:"cbc:Amount,omitempty"`
	PaymentDueDate string         `xml:"cbc:PaymentDueDate,omitempty"`
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
	ID                     IDType                 `xml:"cbc:ID"`
	Percent                float64                `xml:"cbc:Percent,omitempty"`
	TaxExemptionReasonCode TaxExemptionReasonCode `xml:"cbc:TaxExemptionReasonCode,omitempty"`
	TaxScheme              TaxScheme              `xml:"cac:TaxScheme"`
}

type TaxExemptionReasonCode struct {
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	ListName       string `xml:"listName,attr,omitempty"`
	ListURI        string `xml:"listURI,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type MonetaryAmount struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Value      float64 `xml:",chardata"`
}

type LegalMonetaryTotal struct {
	LineExtensionAmount MonetaryAmount `xml:"cbc:LineExtensionAmount,omitempty"`
	TaxExclusiveAmount  MonetaryAmount `xml:"cbc:TaxExclusiveAmount,omitempty"`
	TaxInclusiveAmount  MonetaryAmount `xml:"cbc:TaxInclusiveAmount,omitempty"`
	PayableAmount       MonetaryAmount `xml:"cbc:PayableAmount"`
}

type InvoiceLine struct {
	ID                  string           `xml:"cbc:ID"`
	InvoicedQuantity    InvoicedQuantity `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount MonetaryAmount   `xml:"cbc:LineExtensionAmount"`
	PricingReference    PricingReference `xml:"cac:PricingReference,omitempty"`
	TaxTotal            []TaxTotal       `xml:"cac:TaxTotal,omitempty"`
	Item                Item             `xml:"cac:Item"`
	Price               Price            `xml:"cac:Price"`
}

type InvoicedQuantity struct {
	UnitCode               string  `xml:"unitCode,attr"`
	UnitCodeListID         string  `xml:"unitCodeListID,attr,omitempty"`
	UnitCodeListAgencyName string  `xml:"unitCodeListAgencyName,attr,omitempty"`
	Value                  float64 `xml:",chardata"`
}

type PricingReference struct {
	AlternativeConditionPrice []AlternativeConditionPrice `xml:"cac:AlternativeConditionPrice"`
}

type AlternativeConditionPrice struct {
	PriceAmount   MonetaryAmount `xml:"cbc:PriceAmount"`
	PriceTypeCode PriceTypeCode  `xml:"cbc:PriceTypeCode"`
}

type PriceTypeCode struct {
	ListName       string `xml:"listName,attr,omitempty"`
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	ListURI        string `xml:"listURI,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type Item struct {
	Description               []string                   `xml:"cbc:Description"`
	SellersItemIdentification *SellersItemIdentification `xml:"cac:SellersItemIdentification,omitempty"`
	CommodityClassification   []CommodityClassification  `xml:"cac:CommodityClassification,omitempty"`
}

type SellersItemIdentification struct {
	ID string `xml:"cbc:ID"`
}

type CommodityClassification struct {
	ItemClassificationCode ItemClassificationCode `xml:"cbc:ItemClassificationCode"`
}

type ItemClassificationCode struct {
	ListID         string `xml:"listID,attr,omitempty"`
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type Price struct {
	PriceAmount MonetaryAmount `xml:"cbc:PriceAmount"`
}


func GenerateInvoiceXML(doc *models.Document, issuer *models.Company) (*Invoice, error) {
	invoice := &Invoice{
		Xmlns:     "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		XmlnsCac:  "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsCbc:  "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		XmlnsCcts: "urn:un:unece:uncefact:documentation:2",
		XmlnsDs:   "http://www.w3.org/2000/09/xmldsig#",
		XmlnsExt:  "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2",
		XmlnsQdt:  "urn:oasis:names:specification:ubl:schema:xsd:QualifiedDatatypes-2",
		XmlnsUdt:  "urn:un:unece:uncefact:data:specification:UnqualifiedDataTypesSchemaModule:2",
		XmlnsXsi:  "http://www.w3.org/2001/XMLSchema-instance",

		UBLVersionID: "2.1",
		CustomizationID: CustomizationID{
			SchemeAgencyName: "PE:SUNAT",
			Value:            "2.0",
		},
		ProfileID: ProfileID{
			SchemeName:       "Tipo de Operacion",
			SchemeAgencyName: "PE:SUNAT",
			SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo51",
			Value:            "0101", // Venta interna
		},
		ID:                   fmt.Sprintf("%s-%s", doc.Serie, doc.Number),
		IssueDate:            doc.IssueDate.Format("2006-01-02"),
		IssueTime:            doc.IssueDate.Format("15:04:05"),
		DocumentCurrencyCode: DocumentCurrencyCode{
			ListID:         "ISO 4217 Alpha",
			ListName:       "Currency",
			ListAgencyName: "United Nations Economic Commission for Europe",
			Value:          doc.CurrencyCode,
		},

		InvoiceTypeCode: InvoiceTypeCode{
			ListAgencyName: "PE:SUNAT",
			ListName:       "Tipo de Documento", 
			ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo01",
			ListID:         "0101",
			Name:           "Tipo de Operacion",
			ListSchemeURI:  "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo51",
			Value:          string(doc.Type),
		},
		LineCountNumeric: len(doc.Lines),
	}

	if doc.DueDate != nil {
		invoice.DueDate = doc.DueDate.Format("2006-01-02")
	} else {
		// Set DueDate to IssueDate if not provided
		invoice.DueDate = doc.IssueDate.Format("2006-01-02")
	}

	// Add Note with amount in letters (required by SUNAT)
	amountInWords := convertAmountToWords(doc.TotalAmount, doc.CurrencyCode)
	invoice.Note = []Note{
		{
			LanguageLocaleID: "1000",
			Value:            amountInWords,
		},
	}

	// Signature (using document-specific ID instead of hardcoded)
	documentID := fmt.Sprintf("%s-%s", doc.Serie, doc.Number)
	invoice.Signature = []Signature{
		{
			ID: documentID,
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
					URI: fmt.Sprintf("#%s", documentID),
				},
			},
		},
	}

	// Supplier Party
	supplierPartyIdentifications := []PartyIdentification{
		{
			ID: IDType{
				SchemeID:         getDocumentTypeScheme(issuer.DocumentType),
				SchemeName:       "Documento de Identidad",
				SchemeAgencyName: "PE:SUNAT",
				SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
				Value:            issuer.DocumentNumber,
			},
		},
	}

	// NOTE: Establishment code goes in AddressTypeCode, NOT as a second PartyIdentification
	// Based on accepted SUNAT invoice analysis
	establishmentCode := issuer.EstablishmentCode
	if establishmentCode == "" {
		establishmentCode = "0000" // Valor por defecto si no se especifica
	}

	invoice.AccountingSupplierParty = AccountingSupplierParty{
		Party: Party{
			PartyIdentification: supplierPartyIdentifications,
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
						ID: IDType{
							SchemeID:         "6",
							SchemeName:       "SUNAT:Identificador de Documento de Identidad",
							SchemeAgencyName: "PE:SUNAT",
							SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
							Value:            issuer.DocumentNumber,
						},
						Name: "SUNAT",
					},
				},
			},
			PartyLegalEntity: []PartyLegalEntity{
				{
					RegistrationName: issuer.Name,
					RegistrationAddress: &RegistrationAddress{
						ID: IDType{
							SchemeName:       "Ubigeos",
							SchemeAgencyName: "PE:INEI",
							Value:            "150101", // Ubigeo por defecto para Lima
						},
						AddressTypeCode: AddressTypeCode{
							ListAgencyName: "PE:SUNAT",
							ListName:       "Establecimientos anexos",
							Value:          establishmentCode,
						},
						CityName:         issuer.District,
						CountrySubentity: issuer.Province,
						District:         issuer.District,
						AddressLine: &AddressLine{
							Line: issuer.Address,
						},
						Country: &Country{
							IdentificationCode: IDType{
								ListID:         "ISO 3166-1",
								ListAgencyName: "United Nations Economic Commission for Europe",
								ListName:       "Country",
								Value:          issuer.Country,
							},
						},
					},
				},
			},
		},
	}

	// Temporary: skip Contact to test if this fixes the UBL validation
	// if issuer.Email != "" || issuer.Phone != "" {
	// 	invoice.AccountingSupplierParty.Party.Contact = &Contact{
	// 		ElectronicMail: issuer.Email,
	// 		Telephone:      issuer.Phone,
	// 	}
	// }
	customerID := "00000000" // Valor por defecto si no hay documento
	customerName := "CLIENTE VARIOS"

	if doc.Customer.DocumentNumber != "" {
		customerID = doc.Customer.DocumentNumber
	}
	if doc.Customer.Name != "" {
		customerName = doc.Customer.Name
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
						Value:            customerID,
					},
				},
			},
			PartyName: PartyName{
				Name: customerName,
			},
			PartyLegalEntity: []PartyLegalEntity{
				{
					RegistrationName: customerName,
				},
			},
		},
	}

	// Payment Terms (Required by SUNAT Resolution 000193-2020)
	if doc.PaymentTerms == nil {
		return nil, fmt.Errorf("payment terms are required")
	}
	
	if doc.PaymentTerms.PaymentMeansCode == "" {
		return nil, fmt.Errorf("payment means code is required (e.g., 'Contado' or 'Credito')")
	}
	
	paymentTerms := PaymentTerms{
		ID:             "FormaPago",
		PaymentMeansID: doc.PaymentTerms.PaymentMeansCode,
	}
	
	// Add amount and due date for credit payments
	if doc.PaymentTerms.PaymentMeansCode == "Credito" {
		paymentTerms.Amount = MonetaryAmount{
			CurrencyID: doc.CurrencyCode,
			Value:      doc.PaymentTerms.Amount,
		}
		if !doc.PaymentTerms.DueDate.IsZero() {
			paymentTerms.PaymentDueDate = doc.PaymentTerms.DueDate.Format("2006-01-02")
		}
	}
	
	invoice.PaymentTerms = []PaymentTerms{paymentTerms}

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
						ID: IDType{
							SchemeID:         "UN/ECE 5305",
							SchemeName:       "Tax Category Identifier",
							SchemeAgencyName: "United Nations Economic Commission for Europe",
							Value:            getTaxCategoryID(tax.Type),
						},
						Percent: tax.Rate,
						TaxExemptionReasonCode: TaxExemptionReasonCode{
							ListAgencyName: "PE:SUNAT",
							ListName:       "Afectacion del IGV",
							ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
							Value:          getTaxExemptionCode(tax.Type),
						},
						TaxScheme: TaxScheme{
							ID: IDType{
								SchemeID:         "UN/ECE 5153",
								SchemeName:       "Codigo de tributos",
								SchemeAgencyName: "PE:SUNAT",
								Value:            getTaxSchemeID(tax.Type),
							},
							Name:        string(tax.Type),
							TaxTypeCode: "VAT",
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
				UnitCode:               line.UnitCode,
				UnitCodeListID:         "UN/ECE rec 20",
				UnitCodeListAgencyName: "United Nations Economic Commission for Europe",
				Value:                  line.Quantity,
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
			invoiceLine.Item.SellersItemIdentification = &SellersItemIdentification{
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
						ID: IDType{
							SchemeID:         "UN/ECE 5305",
							SchemeName:       "Tax Category Identifier",
							SchemeAgencyName: "United Nations Economic Commission for Europe",
							Value:            getTaxCategoryID(tax.Type),
						},
						Percent: tax.Rate,
						TaxExemptionReasonCode: TaxExemptionReasonCode{
							ListAgencyName: "PE:SUNAT",
							ListName:       "Afectacion del IGV",
							ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
							Value:          getTaxExemptionCode(tax.Type),
						},
						TaxScheme: TaxScheme{
							ID: IDType{
								SchemeID:         "UN/ECE 5153",
								SchemeName:       "Codigo de tributos",
								SchemeAgencyName: "PE:SUNAT",
								Value:            getTaxSchemeID(tax.Type),
							},
							Name:        string(tax.Type),
							TaxTypeCode: "VAT",
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
					PriceTypeCode: PriceTypeCode{
						ListName:       "Tipo de Precio",
						ListAgencyName: "PE:SUNAT",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo16",
						Value:          "01", // Precio unitario (incluye el IGV)
					},
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
		return "0"
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

func getTaxCategoryID(taxType models.TaxType) string {
	switch taxType {
	case models.TaxTypeIGV:
		return "S" // Standard rate (IGV gravado)
	case models.TaxTypeISC:
		return "S" // Standard rate
	case models.TaxTypeICBP:
		return "S" // Standard rate
	default:
		return "S"
	}
}

func getTaxExemptionCode(taxType models.TaxType) string {
	switch taxType {
	case models.TaxTypeIGV:
		return "10" // Gravado - Operación Onerosa
	case models.TaxTypeISC:
		return "20" // Exonerado - Operación Onerosa
	case models.TaxTypeICBP:
		return "20" // Exonerado - Operación Onerosa
	default:
		return "10"
	}
}

// convertAmountToWords converts a numeric amount to words in Spanish
func convertAmountToWords(amount float64, currency string) string {
	integerPart := int(amount)
	decimalPart := int(math.Round((amount - float64(integerPart)) * 100))

	var currencyWord string
	switch currency {
	case "PEN":
		if integerPart == 1 {
			currencyWord = "SOL"
		} else {
			currencyWord = "SOLES"
		}
	case "USD":
		if integerPart == 1 {
			currencyWord = "DÓLAR AMERICANO"
		} else {
			currencyWord = "DÓLARES AMERICANOS"
		}
	default:
		currencyWord = currency
	}

	integerWords := numberToWords(integerPart)
	return fmt.Sprintf("%s CON %02d/100 %s", strings.ToUpper(integerWords), decimalPart, currencyWord)
}

// numberToWords converts a number to words in Spanish
func numberToWords(n int) string {
	if n == 0 {
		return "CERO"
	}

	ones := []string{"", "UNO", "DOS", "TRES", "CUATRO", "CINCO", "SEIS", "SIETE", "OCHO", "NUEVE", "DIEZ",
		"ONCE", "DOCE", "TRECE", "CATORCE", "QUINCE", "DIECISEIS", "DIECISIETE", "DIECIOCHO", "DIECINUEVE"}

	tens := []string{"", "", "VEINTE", "TREINTA", "CUARENTA", "CINCUENTA", "SESENTA", "SETENTA", "OCHENTA", "NOVENTA"}

	hundreds := []string{"", "CIEN", "DOSCIENTOS", "TRESCIENTOS", "CUATROCIENTOS", "QUINIENTOS",
		"SEISCIENTOS", "SETECIENTOS", "OCHOCIENTOS", "NOVECIENTOS"}

	if n < 20 {
		return ones[n]
	}

	if n < 100 {
		if n%10 == 0 {
			return tens[n/10]
		}
		return tens[n/10] + " Y " + ones[n%10]
	}

	if n < 1000 {
		if n == 100 {
			return "CIEN"
		}
		if n%100 == 0 {
			return hundreds[n/100]
		}
		return hundreds[n/100] + " " + numberToWords(n%100)
	}

	if n < 1000000 {
		thousands := n / 1000
		remainder := n % 1000

		var result string
		if thousands == 1 {
			result = "MIL"
		} else {
			result = numberToWords(thousands) + " MIL"
		}

		if remainder > 0 {
			result += " " + numberToWords(remainder)
		}
		return result
	}

	// For millions and above, simplified version
	return fmt.Sprintf("%d", n)
}
