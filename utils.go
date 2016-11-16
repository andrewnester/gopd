package gopd

import (
	"io"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Response struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	Name         string `json:"name"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type PandadocError struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func (e PandadocError) Error() string {
	err := e.Type
	return fmt.Sprintf("%s; %s", err, e.Detail)
}

func SendRequest(method string, url string, data io.Reader, contentType string, expectedStatus string) ([]byte, error) {
	req, err := http.NewRequest(method, url, data)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GetAccessToken()))
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

	if resp.Status != expectedStatus {
		var respErr PandadocError
		_ = json.Unmarshal(body, &respErr)
		respErr.Detail = string(body[:])
		return nil, respErr
	}

	return body, nil
}

