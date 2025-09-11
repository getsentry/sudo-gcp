# Sudo GCP (Go Version)

A Go implementation of sudo-gcp that provides temporary elevated privileges for Google Cloud commands using short-lived OAuth access tokens with mTLS support.

## Features

- ✅ Service account impersonation via IAM API
- ✅ Graceful fallback to regular authentication
- ✅ Simple TOML configuration
- ✅ Command execution with token environment variables
- 🚧 mTLS support for BeyondCorp environments (framework ready)
- 🚧 Token caching (TODO)
- 🚧 macOS keychain integration (TODO)

## Setup

1. Install Go 1.21 or later
2. Set up Google Cloud authentication (ADC, service account key, etc.)
3. Configure your service account and permissions

## Configuration

Create a `sudo-gcp.toml` file:

```toml
service_account = "my-service-account@my-project.iam.gserviceaccount.com"
lifetime = 3600
scopes = [
    "https://www.googleapis.com/auth/cloud-platform"
]
```

## Usage

```bash
# Run a command with elevated privileges
./sudo-gcp gcloud compute instances list

# Use a different config file
./sudo-gcp --config-file /path/to/config.toml gcloud compute instances list
```

## Development Status

### ✅ Implemented
- Basic CLI structure with flag parsing
- TOML configuration parsing with defaults
- Service account impersonation via IAM API
- Command execution with token environment variables
- Graceful fallback from mTLS to regular authentication
- Clear error messages and logging

### 🚧 In Progress
- mTLS certificate loading framework (needs certificate access implementation)
- Error handling improvements

### 📋 TODO
- **Certificate Access**: The go-keychain library only supports generic passwords, not certificates. We need to:
  - Use a different library (like `github.com/99designs/keyring` or direct Security framework calls)
  - Or implement certificate loading from files with keychain integration
- **Token Caching**: Implement secure token storage with expiration
- **Service Account Delegation**: Add support for delegation chains
- **Tests**: Add comprehensive test coverage
- **Build/Release**: Add proper build and release process

## mTLS Implementation Notes

The current implementation has the framework for mTLS support but needs certificate access:

1. **Certificate Loading**: The `loadCertificateFromKeychain()` function is a placeholder
2. **Fallback**: The tool gracefully falls back to regular authentication when mTLS isn't available
3. **Next Steps**: To implement mTLS, we need to:
   - Find a library that can access macOS keychain certificates
   - Or implement direct Security framework calls using cgo
   - Parse the certificate and private key data
   - Create a proper `tls.Certificate`

## Dependencies

- `github.com/BurntSushi/toml` - TOML configuration parsing
- `golang.org/x/oauth2` - OAuth2 authentication
- `google.golang.org/api` - Google Cloud APIs
- `github.com/keybase/go-keychain` - Keychain access (currently limited to passwords)

## Building

```bash
go build -o sudo-gcp .
```

## Testing

```bash
# Test configuration loading (requires GCP credentials)
./sudo-gcp echo "test"

# Test with actual GCP command (requires authentication)
./sudo-gcp gcloud auth list
```

## Advantages over Rust Version

1. **Better Delegation Support**: Using IAM API directly allows for service account delegation chains
2. **Simpler Dependencies**: Fewer moving parts than the Rust version
3. **GCP SDK**: First-class support from Google's Go SDKs
4. **Easier mTLS**: Go's crypto libraries make mTLS implementation straightforward
5. **Cross-platform**: Go's standard library handles many platform differences
