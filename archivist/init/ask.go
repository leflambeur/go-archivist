package archivist

import (
	"fmt"
	"net/url"

	"github.com/AlecAivazis/survey/v2"
)

const (
	authTypeTLS    = "TLS"
	authTypeSecret = "Client Secret"
)

type clientCredentials struct {
	clientTenant string
	clientID     string
	clientSecret string
}

type tlsPaths struct {
	tlsCert string
	tlsKey  string
}

func validateURL(val interface{}) error {
	inputURL, ok := val.(string)
	if !ok {
		return fmt.Errorf("non-string input")
	}
	if _, err := url.ParseRequestURI(inputURL); err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	return nil
}

func askForArchivistDetails() (string, error) {
	url := ""
	q := &survey.Input{
		Message: "Enter Jitsuin Archivist URL: ",
		Help:    "Enter the URL of your Jitsuin Archivist Instance.",
	}
	err := survey.AskOne(q, &url, survey.WithValidator(survey.Required), survey.WithValidator(validateURL))
	if err != nil {
		return "", fmt.Errorf("error asking for archivist details: %v", err)
	}
	return url, nil
}

func askForAuthenticationMethod(options []string) (string, error) {
	selectedAuth := ""
	q := &survey.Select{
		Message: "Select authentication method",
		Options: options,
		Help:    "Select either Client Secret or TLS authentication",
	}
	err := survey.AskOne(q, &selectedAuth, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("error asking for authentication method: %v", err)
	}
	log.Debugf("selected auth: %v", selectedAuth)
	return selectedAuth, nil
}

func askForClientSecretCredentials() (*clientCredentials, error) {
	answers := clientCredentials{}
	qs := []*survey.Question{
		{
			Name: "clientTenant",
			Prompt: &survey.Input{
				Message: "Enter your Azure Tenant",
			},
			Validate: survey.Required,
		},
		{
			Name: "clientID",
			Prompt: &survey.Input{
				Message: "Enter your Client ID (API APP ID)",
			},
			Validate: survey.Required,
		},
		{
			Name: "clientSecret",
			Prompt: &survey.Password{
				Message: "Enter your Client Secret (API APP Secret)",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(qs, &answers); err != nil {
		return nil, fmt.Errorf("error fetching credentials: %v", err)
	}
	return &answers, nil
}

func askForTLSFiles() (*tlsPaths, error) {
	answers := tlsPaths{}
	qs := []*survey.Question{
		{
			Name: "tlsCert",
			Prompt: &survey.Input{
				Message: "Enter the full path for your Archivist TLS Cert",
			},
			Validate: survey.Required,
		},
		{
			Name: "tlsKey",
			Prompt: &survey.Input{
				Message: "Enter the full path to your Archivist TLS Key",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(qs, &answers); err != nil {
		return nil, fmt.Errorf("error fetching credentials: %v", err)
	}
	return &answers, nil
}
