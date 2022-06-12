package auth

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var app *firebase.App
var client *auth.Client

func init() {
	var err error
	credential := generateFirebaseCredential()
	if err != nil {
		log.Fatalf("Unable to generate firebase credential file: %v\n", err)
	}

	ctx := context.Background()
	opt := option.WithCredentialsJSON(credential)
	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	client, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}
}

func generateFirebaseCredential() []byte {
	credentialDataTemplate := "{\"type\":\"%s\",\"project_id\":\"%s\",\"private_key_id\":\"%s\",\"private_key\":\"%s\",\"client_email\":\"%s\",\"client_id\":\"%s\",\"auth_uri\":\"%s\",\"token_uri\":\"%s\",\"auth_provider_x509_cert_url\":\"%s\",\"client_x509_cert_url\":\"%s\"}"
	credentialData := fmt.Sprintf(credentialDataTemplate, os.Getenv("FB_TYPE"), os.Getenv("FB_PROJECT_ID"), os.Getenv("FB_PRIVATE_KEY_ID"), os.Getenv("FB_PRIVATE_KEY"), os.Getenv("FB_CLIENT_EMAIL"), os.Getenv("FB_CLIENT_ID"), os.Getenv("FB_AUTH_URI"), os.Getenv("FB_TOKEN_URI"), os.Getenv("FB_AUTH_PROVIDER_CERT_URL"), os.Getenv("FB_CLIENT_CERT_URL"))

	return []byte(credentialData)
}
