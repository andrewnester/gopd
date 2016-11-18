package gopd

import (
	"encoding/json"
	"fmt"
)

const TEMPLATE_API_ENDPOINT = "https://api.pandadoc.com/public/v1/templates"

func GetTemplateList(params ...int) (*TemplateList, error) {
	if len(params) > 1 {
		return nil, PandadocError{
			Type:"invalid_function_call",
			Detail:"GetTemplateList expected 0 or 1 argument. You can pass only page number to GetTemplateList",
		}
	}

	page := 1
	if len(params) == 1 {
		page = params[0]
	}

	body, err := SendRequest(
		"GET",
		fmt.Sprintf("%s?page=%d", TEMPLATE_API_ENDPOINT, page),
		nil,
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}
	var tList TemplateList
	err = json.Unmarshal(body, &tList)
	if err != nil {
		return nil, err
	}

	return &tList, nil
}

func GetTemplateDetails(templateId string) (*Template, error) {
	body, err := SendRequest(
		"GET",
		fmt.Sprintf("%s/%s/details", TEMPLATE_API_ENDPOINT, templateId),
		nil,
		"application/json",
		"200 OK",
	)
	if err != nil {
		return nil, err
	}
	var t Template
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

