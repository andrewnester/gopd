package gopd

import (
	"encoding/json"
	"fmt"
)

const TEMPLATE_API_ENDPOINT = "https://api.pandadoc.com/public/v1/templates"

func GetTemplateList() (*TemplateList, error) {
	body, err := SendRequest(
		"GET",
		TEMPLATE_API_ENDPOINT,
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

