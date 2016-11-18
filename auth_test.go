package gopd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"io/ioutil"
	"net/url"
)

var auth Auth = Auth{"some-client-id", "some-client-secret"}

func TestAuth_AuthenticateUrl(t *testing.T) {

	authUrl := auth.AuthenticateUrl("http://redirect/", "read+write")

	a := assert.New(t)
	a.Equal(
		fmt.Sprintf(
			"https://app.pandadoc.com/oauth2/authorize?response_type=code&client_id=%s&redirect_url=http://redirect/&scope=read+write",
			auth.ClientId,
		),
		authUrl,
	)
}

func TestAuth_CreateAccessToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := ioutil.ReadFile("./fixtures/auth/access.json")

	httpmock.RegisterResponder("POST", ACCESS_TOKEN_URL,
		func(req *http.Request) (*http.Response, error) {
			a := assert.New(t)
			a.Equal("application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

			body, _ := ioutil.ReadAll(req.Body)
			bodyStr := string(body[:])
			expectedBody := url.Values{
				"client_id": {auth.ClientId},
				"client_secret": {auth.ClientSecret},
				"code": {"auth-code"},
				"grant_type": {"authorization_code"},
				"redirect_uri": {"http://redirect"},
				"scope": {"read+write"},
			}.Encode()
			a.Equal(expectedBody, bodyStr)
			return httpmock.NewStringResponse(200, string(file[:])), nil
		})

	creds, err := auth.CreateAccessToken("auth-code", "read+write", "http://redirect")
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal("2ff2dfe36322448c6953616740a910be57bbd4ca", creds.AccessToken)
	a.Equal("4c82f23d91a75961f4d08134fc5ad0dfe6a4c36a", creds.RefreshToken)
	a.Equal(31535999, creds.ExpiredIn)
	a.Equal("read+write", creds.Scope)
}

func TestAuth_RefreshToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := ioutil.ReadFile("./fixtures/auth/refresh.json")

	httpmock.RegisterResponder("POST", ACCESS_TOKEN_URL,
		func(req *http.Request) (*http.Response, error) {
			a := assert.New(t)
			a.Equal("application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

			body, _ := ioutil.ReadAll(req.Body)
			bodyStr := string(body[:])
			expectedBody := url.Values{
				"client_id": {auth.ClientId},
				"client_secret": {auth.ClientSecret},
				"grant_type": {"refresh_token"},
				"refresh_token": {"refresh-code"},
				"scope": {"read+write"},
			}.Encode()
			a.Equal(expectedBody, bodyStr)
			return httpmock.NewStringResponse(200, string(file[:])), nil
		})

	creds, err := auth.RefreshToken("refresh-code", "read+write")
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal("gf22dfac6322448chj5sv16740a910be57bad45b", creds.AccessToken)
	a.Equal("4c82f23d91a75961f4d08134fc5ad0dfe6a4c36a", creds.RefreshToken)
	a.Equal(31535999, creds.ExpiredIn)
	a.Equal("read+write", creds.Scope)

	a.Equal(creds, GetCredentials())
}
