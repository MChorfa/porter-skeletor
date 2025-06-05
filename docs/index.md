# Porter Mixin Generator (Skeletor)

[![Build Status](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml/badge.svg)](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml)
[![GitHub Pages](https://github.com/getporter/skeletor/actions/workflows/pages.yml/badge.svg)](https://getporter.github.io/skeletor/)
[![Release](https://img.shields.io/github/v/release/getporter/skeletor)](https://github.com/getporter/skeletor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/getporter/skeletor)](https://goreportcard.com/report/github.com/getporter/skeletor)

**Skeletor** is an enterprise-grade command-line tool designed to streamline the creation of new Porter mixins with built-in security, compliance, authentication, and observability features. It scaffolds production-ready mixin projects based on configurable templates, providing a solid foundation for enterprise environments.

## Quick Start

### Basic Mixin Generation

```bash
# Install Skeletor
go install github.com/getporter/skeletor/cmd/skeletor@latest

# Create a basic mixin
skeletor create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin"
```

### Enterprise Mixin with All Features

```bash
# Create an enterprise-ready mixin with all security, compliance, auth, and observability features
skeletor create \
  --name enterprise-mixin \
  --author "Enterprise Team" \
  --module "github.com/your-org/enterprise-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,policy_enforcement" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa" \
  --enable-auth \
  --auth-features "rbac,ldap,sso,mfa,vault" \
  --enable-observability \
  --observability-features "apm,opentelemetry,audit_logging,distributed_tracing"
```

## Core Capabilities

* **Rapid Scaffolding:** Generate production-ready Porter mixin projects in seconds
* **Enterprise-Grade Templates:** Built-in support for security, compliance, and observability
* **Flexible Configuration:** Interactive and non-interactive modes with extensive customization
* **Template Engine:** Advanced Go template support with custom functions and conditional rendering
* **Post-Generation Validation:** Automatic code formatting, dependency resolution, and build verification

## Enterprise Features

### Security Features (`--enable-security`)

* **Input Validation:** Comprehensive input sanitization and validation
* **Rate Limiting:** Configurable request throttling and abuse prevention  
* **Secure Headers:** HTTP security headers and CORS configuration
* **Vulnerability Scanning:** Integrated security scanning with Gosec
* **Policy Enforcement:** Role-based access control and security policies

### Compliance Frameworks (`--enable-compliance`)

* **SOC 2:** System and Organization Controls compliance templates
* **GDPR:** General Data Protection Regulation compliance features
* **HIPAA:** Health Insurance Portability and Accountability Act support
* **PCI DSS:** Payment Card Industry Data Security Standard templates

### Authentication & Authorization (`--enable-auth`)

* **RBAC:** Role-Based Access Control implementation
* **LDAP Integration:** Enterprise directory service integration
* **SSO Support:** Single Sign-On with SAML/OAuth2/OIDC
* **MFA:** Multi-Factor Authentication implementation
* **HashiCorp Vault:** Secrets management integration
* **Session Management:** Secure session handling and lifecycle management

### Enhanced Observability (`--enable-observability`)

* **APM Integration:** Application Performance Monitoring setup
* **Infrastructure Monitoring:** System metrics and health checks
* **Custom Metrics:** Business-specific metric collection
* **Health Checks:** Comprehensive health endpoint implementation
* **OpenTelemetry:** Distributed tracing and observability
* **Audit Logging:** Comprehensive audit trail and compliance logging
* **Distributed Tracing:** End-to-end request tracing across services

## Installation

### Package Managers

**Homebrew (macOS/Linux):**
```bash
brew install getporter/tap/skeletor
```

**Go Install:**
```bash
go install github.com/getporter/skeletor/cmd/skeletor@latest
```

### Binary Downloads

Download pre-built binaries from the [releases page](https://github.com/getporter/skeletor/releases):

```bash
# Linux (amd64)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_linux_amd64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# macOS (amd64)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_darwin_amd64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# Windows (amd64)
# Download skeletor_windows_amd64.zip from releases page
```

### Docker

```bash
# Pull the latest image
docker pull ghcr.io/getporter/skeletor:latest

# Create a mixin using Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin"
```

### Build from Source

```bash
git clone https://github.com/getporter/skeletor.git
cd skeletor
go run mage.go build install
# Binary will be in ./bin/skeletor
```

## Next Steps

* [Command Reference](command-reference.md) - Complete CLI documentation
* [Enterprise Features Guide](enterprise-features.md) - Detailed enterprise features documentation
* [Template Customization](template-customization.md) - How to customize templates
* [Contributing](contributing.md) - How to contribute to Skeletor
* [Examples](examples.md) - Real-world usage examples

## License

Apache 2.0 License - see [LICENSE](https://github.com/getporter/skeletor/blob/main/LICENSE) for details.
