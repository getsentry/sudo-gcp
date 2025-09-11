package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceAccount string   `toml:"service_account"`
	Scopes         []string `toml:"scopes,omitempty"`
	Lifetime       int      `toml:"lifetime,omitempty"`
}

type TokenInfo struct {
	AccessToken string    `json:"access_token"`
	ExpireTime  time.Time `json:"expire_time"`
	Scopes      []string  `json:"scopes"`
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config-file", "./sudo-gcp.toml", "Path to config file")
	flag.Parse()

	// Parse remaining args as the command to execute
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No command specified")
	}

	// Load configuration
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get access token
	token, err := getAccessToken(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Execute command with token
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(),
		"GOOGLE_OAUTH_ACCESS_TOKEN="+token,
		"CLOUDSDK_AUTH_ACCESS_TOKEN="+token,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

func loadConfig(path string) (*Config, error) {
	var config Config

	// Try to read config file
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Set defaults
	if config.Lifetime == 0 {
		config.Lifetime = 3600 // 1 hour default
	}
	if len(config.Scopes) == 0 {
		config.Scopes = []string{"https://www.googleapis.com/auth/cloud-platform"}
	}

	return &config, nil
}

func getAccessToken(ctx context.Context, config *Config) (string, error) {
	// TODO: Implement token caching and retrieval
	// For now, just get a fresh token
	return getFreshToken(ctx, config)
}
