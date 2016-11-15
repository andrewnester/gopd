package gopd

import (
	"net/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

const DOCUMENT_API_ENDPOINT = "https://api.pandadoc.com/public/v1/documents"

type Recipient struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role,omitempty"`
}

type Token struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type FieldName struct {
	Value string `json:"value"`
}

type Field struct {
	FieldName FieldName `json:"field_name"`
}

type SectionRowOptions struct {
	Optional        bool `json:"optional"`
	MayEditQuantity bool `json:"may_edit_quantity"`
}

type SectionRowData struct {
	Qty         int `json:"qty"`
	Name        string `json:"name"`
	Cost        string `json:"cost"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Discount    int `json:"discount"`
}

type SectionRow struct {
	Options       SectionRowOptions `json:"options"`
	Data          SectionRowData `json:"data"`
	CustomFields  map[string]string `json:"custom_fields"`
	CustomColumns map[string]string `json:"custom_columns"`
}

type Section struct {
	Title   string `json:"title"`
	Default bool `json:"default"`
	Rows    []SectionRow `json:"rows"`
}

type PricingTable struct {
	Name     string `json:"name"`
	Sections []Section `json:"sections"`
}

type TemplateDocument struct {
	Name          string `json:"name"`
	TemplateUuid  string `json:"template_uuid"`
	Recipients    []Recipient `json:"recipients"`
	Tokens        []Token `json:"tokens,omitempty"`
	Fields        []Field `json:"fields,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	PricingTables []PricingTable `json:"pricing_tables,omitempty"`
}

type PdfDocument struct {

}

type Response struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	Name         string `json:"name"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type DocumentError struct {
	Type   string `json:"type"`
	Detail map[string][]string `json:"detail"`
}

type DocumentStatus struct {
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	Status       string `json:"status"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

func (e DocumentError) Error() string {
	err := e.Type
	for _, val := range e.Detail {
		err = fmt.Sprintf("%s; %s", err, val)
	}
	return err
}

func (d TemplateDocument) Create(accessToken string) (*Response, error) {
	jsonStr, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", DOCUMENT_API_ENDPOINT, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Status != "201 Created" {
		var respErr DocumentError
		_ = json.Unmarshal(body, &respErr)
		return nil, respErr
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (d PdfDocument) Create() {

}

func GetDocumentStatus(accessToken string, docId string) (*DocumentStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DOCUMENT_API_ENDPOINT, docId), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Status != "200 OK" {
		var respErr DocumentError
		_ = json.Unmarshal(body, &respErr)
		return nil, respErr
	}

	var status DocumentStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}