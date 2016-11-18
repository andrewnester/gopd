package gopd

import (
	"fmt"
	"net/url"
	"encoding/json"
	"bytes"
)

const AUTH_URL = "https://app.pandadoc.com/oauth2/authorize?response_type=code&client_id=%s&redirect_url=%s&scope=%s"
const ACCESS_TOKEN_URL = "https://api.pandadoc.com/oauth2/access_token"

var credentials Credentials = Credentials{}

type Auth struct {
	ClientId     string
	ClientSecret string
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	ExpiredIn    int  `json:"expired_in"`
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
	body, err := SendRequest(
		"POST",
		ACCESS_TOKEN_URL,
		bytes.NewBufferString(url.Values{
			"grant_type": {"authorization_code"},
			"client_id": {a.ClientId},
			"client_secret": {a.ClientSecret},
			"code": {code},
			"scope": {scope},
			"redirect_uri": {redirect_uri},
		}.Encode()),
		"application/x-www-form-urlencoded",
		"200 OK")
	if err != nil {
		return nil, err
	}

	var creds Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return nil, err
	}
	SetCredentials(creds)
	return &creds, nil
}

func (a Auth) RefreshToken(refresh_token string, scope string) (*Credentials, error) {
	body, err := SendRequest(
		"POST",
		ACCESS_TOKEN_URL,
		bytes.NewBufferString(url.Values{
			"grant_type": {"refresh_token"},
			"client_id": {a.ClientId},
			"client_secret": {a.ClientSecret},
			"refresh_token": {refresh_token},
			"scope": {scope},
		}.Encode()),
		"application/x-www-form-urlencoded",
		"200 OK")
	if err != nil {
		return nil, err
	}

	var creds Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return nil, err
	}

	SetCredentials(creds)
	return &creds, nil
}

func GetAccessToken() string {
	return credentials.AccessToken
}

func GetCredentials() *Credentials {
	return &credentials
}

func SetCredentials(creds Credentials) {
	credentials = creds
}

