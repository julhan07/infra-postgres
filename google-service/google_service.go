package googleservice

import (
	"context"
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"
const userInfoScope = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"

type GoogleService struct {
	Key              string
	FirebaseAudience string
}

func NewGoogleService(Key string, FirebaseAudience string) GoogleService {
	return GoogleService{Key, FirebaseAudience}
}

// newTokenProvider function to get token for fcm-send
func (gs GoogleService) LoadAccessToken(gcLocation string) (string, error) {
	jsonKey, err := os.ReadFile(gcLocation)
	if err != nil {
		return "", errors.New("fcm: failed to read credentials file")
	}

	cfg, err := google.JWTConfigFromJSON(jsonKey, firebaseScope, userInfoScope)
	if err != nil {
		return "", errors.New("fcm: failed to get JWT config for the firebase.messaging and userinfo scopes")
	}

	ts := cfg.TokenSource(context.Background())
	token, err := ts.Token()
	return token.AccessToken, err
}

func (gs GoogleService) LoadAccessTokenByKey() (string, error) {
	// Membuat konfigurasi OAuth secara langsung
	cfg := &oauth2.Config{
		Scopes:   []string{userInfoScope},
		Endpoint: google.Endpoint,
		// ClientID dan ClientSecret tidak diperlukan karena tidak ada service-account.json
	}
	// Menambahkan API Key ke konfigurasi OAuth
	cfg.ClientID = gs.Key
	// Membuat token source dari konfigurasi OAuth
	ts := cfg.TokenSource(context.Background(), nil)

	// Mendapatkan token
	token, err := ts.Token()
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
