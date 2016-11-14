package gopd

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"errors"
)

const AUTH_URL = "https://app.pandadoc.com/oauth2/authorize?response_type=code&client_id=%s&redirect_url=%s&scope=%s"
const ACCESS_TOKEN_URL = "https://api.pandadoc.com/oauth2/access_token"

type Auth struct {
	ClientId     string
	ClientSecret string
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Error struct {
	Error string
}

func (a Auth) AuthenticateUrl(redirectUrl string, scope string) string {
	return fmt.Sprintf(AUTH_URL, a.ClientId, redirectUrl, scope)
}

func (a Auth) CreateAccessToken(code string, scope string, redirect_uri string) (*Credentials, error) {
	resp, err := http.PostForm(ACCESS_TOKEN_URL, url.Values{
		"grant_type": {"authorization_code"},
		"client_id": {a.ClientId},
		"client_secret": {a.ClientSecret},
		"code": {code},
		"scope": {scope},
		"redirect_uri": {redirect_uri},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return nil, err
	}

	var respErr Error
	_ = json.Unmarshal(body, &respErr)
	if respErr.Error != "" {
		return nil, errors.New(respErr.Error)
	}
	return &creds, nil
}