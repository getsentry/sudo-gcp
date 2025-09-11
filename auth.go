package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

// getFreshToken implements service account impersonation with mTLS support
func getFreshToken(ctx context.Context, config *Config) (string, error) {
	// Create HTTP client with mTLS if certificate is available
	client, err := createMTLSClient(ctx)
	if err != nil {
		log.Printf("mTLS setup failed, falling back to regular client: %v", err)
		// Fall back to regular client if mTLS setup fails
		client, err = google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			return "", fmt.Errorf("failed to create HTTP client: %w", err)
		}
	} else {
		log.Printf("Using mTLS client with certificate from keychain")
	}

	// Create IAM credentials service
	iamService, err := iamcredentials.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("failed to create IAM credentials service: %w", err)
	}

	// Prepare the request
	req := &iamcredentials.GenerateAccessTokenRequest{
		Scope:    config.Scopes,
		Lifetime: fmt.Sprintf("%ds", config.Lifetime),
	}

	// Generate access token
	resp, err := iamService.Projects.ServiceAccounts.GenerateAccessToken(
		fmt.Sprintf("projects/-/serviceAccounts/%s", config.ServiceAccount),
		req,
	).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return resp.AccessToken, nil
}

// createMTLSClient creates an HTTP client with mTLS support using certificates from keychain
func createMTLSClient(ctx context.Context) (*http.Client, error) {
	// Try to load certificate from keychain
	cert, err := loadCertificateFromKeychain()
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate from keychain: %w", err)
	}

	// Create TLS config with client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		// Add other TLS settings as needed for your environment
		MinVersion: tls.VersionTLS12,
	}

	// Create HTTP client with custom transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Use a reasonable timeout
	}

	return client, nil
}
