package shared

import (
	"bytescheme/common/auth"
	"bytescheme/common/log"
	"bytescheme/common/util"
	"bytescheme/controller/model"
)

const (
	// APIKeyPrefix is the prefix for api keys in the store
	APIKeyPrefix = "apikey/"

	// UserKeyPrefix is the prefix for users in the store
	UserKeyPrefix = "user/"

	// SecretKeyPrefix is the secret prefix
	SecretKeyPrefix = "secret/"

	// GoogleOAuth2ClientID is the Google OAuth2 client ID
	GoogleOAuth2ClientID = "91456297737-d1p2ha4n2847bpsrdrcp72uhp614ar9q.apps.googleusercontent.com"

	// GoogleOAuth2ClientSecretKey is the Google OAuth2 client secret
	GoogleOAuth2ClientSecretKey = SecretKeyPrefix + "google-oauth2"

	// ControllerRedirectURL is the redirect URL
	ControllerRedirectURL = "https://bytescheme.mynetgear.com/controlboard.html"
)

var (
	// Authenticators contain all the registered authenticators
	Authenticators = []auth.Authenticator{
		auth.AuthenticatorFunc(VerifyToken),
		auth.AuthenticatorFunc(VerifyAPIKey),
	}
)

// VerifyToken verifies the token
func VerifyToken(token string) (*auth.Principal, error) {
	clientSecretKey := SecretKeyPrefix + GoogleOAuth2ClientID
	value, err := Store.Get(clientSecretKey)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	if value == nil {
		return nil, nil
	}
	conf := auth.OAuth2Config{
		ClientID:     GoogleOAuth2ClientID,
		ClientSecret: *value,
		RedirectURL:  ControllerRedirectURL,
	}
	principal, err := conf.GetPrincipal(token)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	if principal == nil {
		return nil, nil
	}
	userKey := UserKeyPrefix + principal.Email
	value, err = Store.Get(userKey)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	if value == nil {
		// Not found..return nil to show 401 error
		log.Errorf("User %s is not authorized", principal.Email)
		return nil, nil
	}
	dbPrincipal := &auth.Principal{}
	err = util.ConvertFromJSON([]byte(*value), dbPrincipal)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	dbPrincipal.Email = principal.Email
	dbPrincipal.Name = principal.Name
	return dbPrincipal, nil
}

// VerifyAPIKey verifies the API key
func VerifyAPIKey(apiKey string) (*auth.Principal, error) {
	storeAPIKey := APIKeyPrefix + apiKey
	value, err := Store.Get(storeAPIKey)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	if value == nil {
		// Not found..return nil to show 401 error
		log.Errorf("Key %s is not authorized", apiKey)
		return nil, nil
	}
	principal := &auth.Principal{}
	err = util.ConvertFromJSON([]byte(*value), principal)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return principal, nil
}
