package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/datateamsix/email-sentinel/internal/config"
)

// LoadCredentials reads the OAuth credentials from credentials.json
func LoadCredentials(credPath string) (*oauth2.Config, error) {
	data, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w", err)
	}

	config, err := google.ConfigFromJSON(data, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %w", err)
	}

	return config, nil
}

// GetTokenFromWeb starts the OAuth flow and returns a token
func GetTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Generate auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("")
	fmt.Println("🔐 Gmail Authorization Required")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Println("1. Open this URL in your browser:")
	fmt.Println("")
	fmt.Println(authURL)
	fmt.Println("")
	fmt.Println("2. Authorize the application")
	fmt.Println("3. Copy the authorization code and paste it below")
	fmt.Println("")
	fmt.Print("Enter authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	// Exchange auth code for token
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange code for token: %w", err)
	}

	return token, nil
}

// SaveToken saves the OAuth token to the config directory
func SaveToken(token *oauth2.Token) error {
	tokenPath, err := config.TokenPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return err
	}

	file, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("unable to create token file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("unable to encode token: %w", err)
	}

	return nil
}

// LoadToken loads a previously saved token
func LoadToken() (*oauth2.Token, error) {
	tokenPath, err := config.TokenPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// TokenExists checks if a valid token file exists
func TokenExists() bool {
	tokenPath, err := config.TokenPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(tokenPath)
	return err == nil
}
