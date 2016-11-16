package gopd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const DOCUMENT_API_ENDPOINT = "https://api.pandadoc.com/public/v1/documents"

type FromTemplateDocument struct {
	Name          string `json:"name"`
	TemplateUuid  string `json:"template_uuid"`
	Recipients    []Recipient `json:"recipients"`
	Tokens        []Token `json:"tokens,omitempty"`
	Fields        []Field `json:"fields,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	PricingTables []PricingTable `json:"pricing_tables,omitempty"`
}

type FromPdfDocument struct {
	//TODO: implement creating from PDF
}

func (d FromTemplateDocument) Create() (*Response, error) {
	jsonStr, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	body, err := SendRequest(
		"POST",
		DOCUMENT_API_ENDPOINT,
		bytes.NewBuffer(jsonStr),
		"application/json",
		"201 Created",
	)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (d FromPdfDocument) Create() {

}

func GetDocumentStatus(docId string) (*DocumentStatus, error) {
	body, err := SendRequest(
		"GET",
		fmt.Sprintf("%s/%s", DOCUMENT_API_ENDPOINT, docId),
		nil,
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}

	var status DocumentStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func GetDocumentDetails(docId string) (*Document, error) {
	body, err := SendRequest(
		"GET",
		fmt.Sprintf("%s/%s/details", DOCUMENT_API_ENDPOINT, docId),
		nil,
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}
	var document Document
	err = json.Unmarshal(body, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func SendDocument(docId string, msg string, silent bool) (*Response, error) {
	type SendMsg struct {
		Message string `json:"message"`
		Silent  bool `json:"silent"`
	}
	jsonStr, err := json.Marshal(SendMsg{msg, silent})
	if err != nil {
		return nil, err
	}
	body, err := SendRequest(
		"POST",
		fmt.Sprintf("%s/%s/send", DOCUMENT_API_ENDPOINT, docId),
		bytes.NewBuffer(jsonStr),
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type ShareResponse struct {
	Id        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}

func ShareDocument(docId string, recipient string, lifetime int) (*ShareResponse, error) {
	type ShareMsg struct {
		Recipient string `json:"recipient"`
		Lifetime  int `json:"lifetime"`
	}
	jsonStr, err := json.Marshal(ShareMsg{recipient, lifetime})
	if err != nil {
		return nil, err
	}
	body, err := SendRequest(
		"POST",
		fmt.Sprintf("%s/%s/session", DOCUMENT_API_ENDPOINT, docId),
		bytes.NewBuffer(jsonStr),
		"application/json",
		"201 Created",
	)
	if err != nil {
		return nil, err
	}

	var resp ShareResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func DownloadDocument(docId string) ([]byte, error) {
	body, err := SendRequest(
		"GET",
		fmt.Sprintf("%s/%s/download", DOCUMENT_API_ENDPOINT, docId),
		nil,
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}

	return body, nil
}

