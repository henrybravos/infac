package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	URL      string
	Username string
	Password string
	Timeout  time.Duration
	client   *http.Client
}

type Envelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Xmlns   string   `xml:"xmlns:soap,attr"`
	Header  *Header  `xml:"soap:Header,omitempty"`
	Body    Body     `xml:"soap:Body"`
}

type Header struct {
	Security Security `xml:"wsse:Security"`
}

type Security struct {
	Xmlns            string           `xml:"xmlns:wsse,attr"`
	UsernameToken    UsernameToken    `xml:"wsse:UsernameToken"`
}

type UsernameToken struct {
	Username string `xml:"wsse:Username"`
	Password string `xml:"wsse:Password"`
}

type Body struct {
	Content interface{} `xml:",omitempty"`
	Fault   *Fault      `xml:"soap:Fault,omitempty"`
}

type Fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

type SendBillRequest struct {
	XMLName  xml.Name `xml:"ser:sendBill"`
	Xmlns    string   `xml:"xmlns:ser,attr"`
	FileName string   `xml:"fileName"`
	ContentFile string `xml:"contentFile"`
}

type SendBillResponse struct {
	XMLName            xml.Name `xml:"sendBillResponse"`
	ApplicationResponse string   `xml:"applicationResponse"`
}

type SendSummaryRequest struct {
	XMLName  xml.Name `xml:"ser:sendSummary"`
	Xmlns    string   `xml:"xmlns:ser,attr"`
	FileName string   `xml:"fileName"`
	ContentFile string `xml:"contentFile"`
}

type SendSummaryResponse struct {
	XMLName xml.Name `xml:"sendSummaryResponse"`
	Ticket  string   `xml:"ticket"`
}

type GetStatusRequest struct {
	XMLName xml.Name `xml:"ser:getStatus"`
	Xmlns   string   `xml:"xmlns:ser,attr"`
	Ticket  string   `xml:"ticket"`
}

type GetStatusResponse struct {
	XMLName xml.Name `xml:"getStatusResponse"`
	Status  Status   `xml:"status"`
}

type Status struct {
	StatusCode string `xml:"statusCode"`
	Content    string `xml:"content,omitempty"`
	Error      string `xml:"error,omitempty"`
}

func NewClient(url, username, password string) *Client {
	return &Client{
		URL:      url,
		Username: username,
		Password: password,
		Timeout:  30 * time.Second,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
	}
}

func (c *Client) SendBill(fileName string, content []byte) (*SendBillResponse, error) {
	request := &SendBillRequest{
		Xmlns:       "http://service.sunat.gob.pe",
		FileName:    fileName,
		ContentFile: string(content),
	}
	
	envelope := &Envelope{
		Xmlns: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: &Header{
			Security: Security{
				Xmlns: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
				UsernameToken: UsernameToken{
					Username: c.Username,
					Password: c.Password,
				},
			},
		},
		Body: Body{
			Content: request,
		},
	}
	
	var response SendBillResponse
	err := c.call(envelope, &response)
	return &response, err
}

func (c *Client) SendSummary(fileName string, content []byte) (*SendSummaryResponse, error) {
	request := &SendSummaryRequest{
		Xmlns:       "http://service.sunat.gob.pe",
		FileName:    fileName,
		ContentFile: string(content),
	}
	
	envelope := &Envelope{
		Xmlns: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: &Header{
			Security: Security{
				Xmlns: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
				UsernameToken: UsernameToken{
					Username: c.Username,
					Password: c.Password,
				},
			},
		},
		Body: Body{
			Content: request,
		},
	}
	
	var response SendSummaryResponse
	err := c.call(envelope, &response)
	return &response, err
}

func (c *Client) GetStatus(ticket string) (*GetStatusResponse, error) {
	request := &GetStatusRequest{
		Xmlns:  "http://service.sunat.gob.pe",
		Ticket: ticket,
	}
	
	envelope := &Envelope{
		Xmlns: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: &Header{
			Security: Security{
				Xmlns: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
				UsernameToken: UsernameToken{
					Username: c.Username,
					Password: c.Password,
				},
			},
		},
		Body: Body{
			Content: request,
		},
	}
	
	var response GetStatusResponse
	err := c.call(envelope, &response)
	return &response, err
}

func (c *Client) call(envelope *Envelope, response interface{}) error {
	// Marshal the envelope to XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SOAP envelope: %w", err)
	}
	
	// Add XML declaration
	xmlRequest := append([]byte(xml.Header), xmlData...)
	
	// Create HTTP request
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(xmlRequest))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")
	
	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}
	
	// Log the raw response for debugging
	fmt.Printf("SUNAT Response: %s\n", string(body))
	
	// Try to parse with different envelope formats
	var soapResp Envelope
	err = xml.Unmarshal(body, &soapResp)
	if err != nil {
		// Try with standard Envelope format
		var standardResp struct {
			XMLName xml.Name `xml:"Envelope"`
			Body    struct {
				Content interface{} `xml:",omitempty"`
				Fault   *Fault      `xml:"Fault,omitempty"`
			} `xml:"Body"`
		}
		
		err2 := xml.Unmarshal(body, &standardResp)
		if err2 != nil {
			return fmt.Errorf("failed to unmarshal SOAP response (tried both formats): %w, %w", err, err2)
		}
		
		// Check for fault in standard format
		if standardResp.Body.Fault != nil {
			return fmt.Errorf("SOAP fault: %s - %s", standardResp.Body.Fault.Code, standardResp.Body.Fault.String)
		}
	} else {
		// Check for SOAP fault
		if soapResp.Body.Fault != nil {
			return fmt.Errorf("SOAP fault: %s - %s", soapResp.Body.Fault.Code, soapResp.Body.Fault.String)
		}
	}
	
	// Extract the actual response from the SOAP body
	// Try to parse the response directly
	err = xml.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return nil
}