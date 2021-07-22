package archivist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	FiveMinutes      = 5 * 60 * 1000
	TokenDescription = "Token for Archivist"
	contentType      = "application/x-www-form-urlencoded"
)

type azureLoginResponse struct {
	Token_type   string `json:"token_type"`
	Expires_in   string `json:"expires_in"`
	Access_token string `json:"access_token"`
}

type Profiles struct {
	Profile []Profile `json:"profiles"`
}

type Profile struct {
	Profile_name  string `json:"name"`
	Rkvst_url     string `json:"rkvst-url"`
	Tenant_id     string `json:"tenant-id"`
	Client_id     string `json:"client-id"`
	Client_secret string `json:"client-secret"`
}

func ClientSecretLogin(archivistURL, ClientTenant, ClientID, ClientSecret string) (string, error) {
	azureLoginURL := "https://login.microsoftonline.com/" + ClientTenant + "/oauth2/token"
	return generateAADClientToken(azureLoginURL, archivistURL, ClientID, ClientSecret)
}

func tlsLogin(archivistURL, ClientTenant, ClientID, tlsCert, tlsKey string) (string, error) {
	azureLoginURL := "https://auth.microsoftonline.com/" + ClientTenant + "/oauth2/token"
	return generateAADTLSToken(azureLoginURL, archivistURL, ClientID, tlsCert, tlsKey)
}

func generateAADClientToken(azureLoginURL, archivistURL, ClientID, ClientSecret string) (string, error) {
	var tokenResponse azureLoginResponse
	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", ClientID)
	payload.Set("client_secret", ClientSecret)
	payload.Set("resource", archivistURL)

	resp, err := http.Post(azureLoginURL, contentType, strings.NewReader(payload.Encode()))
	if err != nil {
		return "\nError Generating Token\n\n", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	tokenErr := json.Unmarshal(body, &tokenResponse)
	if tokenErr != nil {
		return "\nError Unmarshalling Token Response\n\n", tokenErr
	}

	return tokenResponse.Access_token, err
}

func generateAADTLSToken(azureLoginURL, archivistURL, ClientID, tlsCert, tlsKey string) (string, error) {
	var tokenResponse azureLoginResponse

	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", ClientID)
	payload.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	//payload.Set("client_assertion", tlsEncoded)
	payload.Set("resource", archivistURL)

	resp, err := http.Post(azureLoginURL, contentType, strings.NewReader(payload.Encode()))
	if err != nil {
		return "\nError Generating Token from Azure\n\n", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	tokenErr := json.Unmarshal(body, &tokenResponse)
	if tokenErr != nil {
		return "\nError Unmarshalling Token\n\n", err
	}
	return tokenResponse.Access_token, err
}

func InitRead(filePath string, selectedProfile string) (string, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "Failed to read file", err
	}
	var profiles Profiles
	json.Unmarshal(byteValue, &profiles)
	for _, i := range profiles.Profile {
		if i.Profile_name == selectedProfile {
			archivistURL := i.Rkvst_url
			clientTenant := i.Tenant_id
			clientID := i.Client_id
			clientSecret := i.Client_secret
			token, err := ClientSecretLogin(archivistURL, clientTenant, clientID, clientSecret)
			if err != nil {
				return "\nError Generating Token\n\n", err
			}
			return token, err
		}
	}
	return "Failed to Read Profile", err
}

func InitAsk(selectedAuth string) (string, error) {
	archivistURL, err := askForArchivistDetails()
	if err != nil {
		return "\nError Validating URL\n\n", err
	}
	if selectedAuth == "" {
		selectedAuth, err = askForAuthenticationMethod()
		if err != nil {
			return "\nError Choosing Auth\n\n", err
		}
	}
	selectedAuth = fmt.Sprintf("%v", selectedAuth)
	switch selectedAuth {
	case "authTypeSecret":
		credentials, err := askForClientSecretCredentials()
		if err != nil {
			return "\nError Entering Credentials\n\n", err
		}
		token, err := ClientSecretLogin(archivistURL, credentials.ClientTenant, credentials.ClientID, credentials.ClientSecret)
		if err != nil {
			return "\nError Generating Token\n\n", err
		}
		return token, err
	case "authTypeTLS":
		credentials, err := askForTLSFiles()
		if err != nil {
			return "\nError Entering TLS Paths\n\n", err
		}
		token, err := tlsLogin(archivistURL, credentials.ClientTenant, credentials.ClientID, credentials.tlsCert, credentials.tlsKey)
		if err != nil {
			return "\nError Generating Token\n\n", err
		}
		ValidateToken(token)
		return token, err
	}
	return "\nShouldn't be here\n\n", err
}
