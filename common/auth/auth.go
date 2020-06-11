package auth

import (
	"fmt"

	token "github.com/futurenda/google-auth-id-token-verifier"
)

// Principal is the security principal
type Principal struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Roles      []string          `json:"roles"`
	Properties map[string]string `json:"properties"`
}

// PathPermission is the placeholder for REST path permission
type PathPermission struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

// Authenticator is the authenticator interface
type Authenticator interface {
	Authenticate(string) (*Principal, error)
}

// AuthenticatorFunc is the authentication handler
type AuthenticatorFunc func(string) (*Principal, error)

// OAuth2Config is the OAuth2 config
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Authenticate is the implemenation of Authenticator
func (handler AuthenticatorFunc) Authenticate(credential string) (*Principal, error) {
	return handler(credential)
}

// GetPrincipal verifies the ID token and returns the principal
func (oauth2Config *OAuth2Config) GetPrincipal(idToken string) (*Principal, error) {
	verifier := token.Verifier{}
	audience := []string{
		oauth2Config.ClientID,
	}
	err := verifier.VerifyIDToken(idToken, audience)
	if err != nil {
		fmt.Printf("Error in verifying token. Error: %s\n", err.Error())
		return nil, err
	}
	claimSet, err := token.Decode(idToken)
	if err != nil {
		fmt.Printf("Error in decoding token. Error: %s\n", err.Error())
		return nil, err
	}
	principal := &Principal{
		Name:  claimSet.Name,
		Email: claimSet.Email,
	}
	return principal, nil
}
