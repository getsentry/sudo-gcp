package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `service_account = "test@test.iam.gserviceaccount.com"
lifetime = 1800
scopes = [
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/compute"
]`

	tmpfile, err := os.CreateTemp("", "test-config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test loading the config
	config, err := loadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the config was loaded correctly
	if config.ServiceAccount != "test@test.iam.gserviceaccount.com" {
		t.Errorf("Expected service account 'test@test.iam.gserviceaccount.com', got '%s'", config.ServiceAccount)
	}

	if config.Lifetime != 1800 {
		t.Errorf("Expected lifetime 1800, got %d", config.Lifetime)
	}

	if len(config.Scopes) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(config.Scopes))
	}

	expectedScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
	}

	for i, scope := range expectedScopes {
		if config.Scopes[i] != scope {
			t.Errorf("Expected scope '%s', got '%s'", scope, config.Scopes[i])
		}
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Create a minimal config file
	configContent := `service_account = "test@test.iam.gserviceaccount.com"`

	tmpfile, err := os.CreateTemp("", "test-config-minimal-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test loading the config
	config, err := loadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify defaults are applied
	if config.Lifetime != 3600 {
		t.Errorf("Expected default lifetime 3600, got %d", config.Lifetime)
	}

	if len(config.Scopes) != 1 {
		t.Errorf("Expected 1 default scope, got %d", len(config.Scopes))
	}

	if config.Scopes[0] != "https://www.googleapis.com/auth/cloud-platform" {
		t.Errorf("Expected default scope 'https://www.googleapis.com/auth/cloud-platform', got '%s'", config.Scopes[0])
	}
}
