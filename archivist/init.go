package archivist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	FiveMinutes      = 5 * 60 * 1000
	TokenDescription = "Token for Archivist"
	contentType      = "application/x-www-form-urlencoded"
)

type azureLoginResponse struct {
	token_type   string `json:"token_type"`
	expires_in   int    `json:"expires_in"`
	access_token string `json:"access.token"`
}

func clientSecretLogin(archivistURL, clientTenant, clientID, clientSecret string) (string, error) {
	azureLoginURL := "https://login.microsoftonline.com/" + clientTenant + "/oauth2/token"
	return generateAADClientToken(azureLoginURL, archivistURL, clientID, clientSecret)
}

func tlsLogin(archivistURL, clientTenant, clientID, tlsCert, tlsKey string) (string, error) {
	azureLoginURL := "https://login.microsoftonline.com/" + clientTenant + "/oauth2/token"
	return generateAADTLSToken(azureLoginURL, archivistURL, clientID, tlsCert, tlsKey)
}

func generateAADClientToken(azureLoginURL, archivistURL, clientID, clientSecret string) (string, error) {
	var tokenResponse azureLoginResponse
	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", clientID)
	payload.Set("client_secret", clientSecret)
	payload.Set("resource", archivistURL)

	resp, err := http.Post(azureLoginURL, contentType, strings.NewReader(payload.Encode()))
	if err != nil {
		fmt.Errorf("Error Generating Token from Azure: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))

	tokenErr := json.Unmarshal(body, &tokenResponse)
	if tokenErr != nil {
		fmt.Errorf("%v", tokenErr)
	}
	return tokenResponse.access_token, err
}

func generateAADTLSToken(azureLoginURL, archivistURL, clientID, tlsCert, tlsKey string) (string, error) {
	var tokenResponse azureLoginResponse

	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", clientID)
	payload.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	payload.Set("client_assertion", tlsEncoded)
	payload.Set("resource", archivistURL)

	resp, err := http.Post(azureLoginURL, contentType, strings.NewReader(payload.Encode()))
	if err != nil {
		fmt.Errorf("Error Generating Token from Azure: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))

	tokenErr := json.Unmarshal(body, &tokenResponse)
	if tokenErr != nil {
		fmt.Errorf("%v", tokenErr)
	}
	return tokenResponse.access_token, err
}

func Init() (string, error) {
	archivistURL, err := askForArchivistDetails()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	selectedAuth, err := askForAuthenticationMethod()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	switch selectedAuth {
	case authTypeSecret:
		credentials, err := askForClientSecretCredentials()
		if err != nil {
			fmt.Errorf("%v", err)
		}
		token, err := clientSecretLogin(archivistURL, credentials.clientTenant, credentials.clientID, credentials.clientSecret)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		return token, err
	case authTypeTLS:
		credentials, err := askForTLSFiles()
		if err != nil {
			fmt.Errorf("%v", err)
		}
		token, err := tlsLogin(archivistURL, credentials.clientTenant, credentials.clientID, credentials.tlsCert, credentials.tlsKey)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		return token, err
	}
	return "Shouldn't be here", err
}
