# Porter Mixin Generator (Skeletor)

[![Build Status](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml/badge.svg)](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml)
[![GitHub Pages](https://github.com/getporter/skeletor/actions/workflows/pages.yml/badge.svg)](https://getporter.github.io/skeletor/)
[![Release](https://img.shields.io/github/v/release/getporter/skeletor)](https://github.com/getporter/skeletor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/getporter/skeletor)](https://goreportcard.com/report/github.com/getporter/skeletor)

**Skeletor** is an enterprise-grade command-line tool designed to streamline the creation of new Porter mixins with built-in security, compliance, authentication, and observability features. It scaffolds production-ready mixin projects based on configurable templates, providing a solid foundation for enterprise environments.

**[📚 View Documentation Site](https://getporter.github.io/skeletor/)** | **[🚀 Quick Start](#quick-start)** | **[🏢 Enterprise Features](#enterprise-features)**

## Features

### 🚀 **Core Capabilities**
* **Rapid Scaffolding:** Generate production-ready Porter mixin projects in seconds
* **Enterprise-Grade Templates:** Built-in support for security, compliance, and observability
* **Flexible Configuration:** Interactive and non-interactive modes with extensive customization
* **Template Engine:** Advanced Go template support with custom functions and conditional rendering
* **Post-Generation Validation:** Automatic code formatting, dependency resolution, and build verification

### 🏢 **Enterprise Features**

#### 🔒 **Security Features** (`--enable-security`)
* **Input Validation:** Comprehensive input sanitization and validation
* **Rate Limiting:** Configurable request throttling and abuse prevention
* **Secure Headers:** HTTP security headers and CORS configuration
* **Vulnerability Scanning:** Integrated security scanning with Gosec
* **Policy Enforcement:** Role-based access control and security policies

#### 📋 **Compliance Frameworks** (`--enable-compliance`)
* **SOC 2:** System and Organization Controls compliance templates
* **GDPR:** General Data Protection Regulation compliance features
* **HIPAA:** Health Insurance Portability and Accountability Act support
* **PCI DSS:** Payment Card Industry Data Security Standard templates

#### 🔐 **Authentication & Authorization** (`--enable-auth`)
* **RBAC:** Role-Based Access Control implementation
* **LDAP Integration:** Enterprise directory service integration
* **SSO Support:** Single Sign-On with SAML/OAuth2/OIDC
* **MFA:** Multi-Factor Authentication implementation
* **HashiCorp Vault:** Secrets management integration
* **Session Management:** Secure session handling and lifecycle management

#### 📊 **Enhanced Observability** (`--enable-observability`)
* **APM Integration:** Application Performance Monitoring setup
* **Infrastructure Monitoring:** System metrics and health checks
* **Custom Metrics:** Business-specific metric collection
* **Health Checks:** Comprehensive health endpoint implementation
* **OpenTelemetry:** Distributed tracing and observability
* **Audit Logging:** Comprehensive audit trail and compliance logging
* **Distributed Tracing:** End-to-end request tracing across services

## Quick Start

### 🚀 **Basic Mixin Generation**

```bash
# Install Skeletor
go install github.com/getporter/skeletor/cmd/skeletor@latest

# Create a basic mixin
skeletor create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin"
```

### 🏢 **Enterprise Mixin with All Features**

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

### 🔒 **Security-Focused Mixin**

```bash
# Create a mixin with enhanced security features
skeletor create \
  --name secure-mixin \
  --author "Security Team" \
  --enable-security \
  --security-features "input_validation,rate_limiting,vulnerability_scanning" \
  --compliance-level "slsa-l3"
```

## Installation

### 📦 **Package Managers**

**Homebrew (macOS/Linux):**
```bash
brew install getporter/tap/skeletor
```

**Go Install:**
```bash
go install github.com/getporter/skeletor/cmd/skeletor@latest
```

### 📥 **Binary Downloads**

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

### 🐳 **Docker**

```bash
# Pull the latest image
docker pull ghcr.io/getporter/skeletor:latest

# Create a mixin using Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin"
```

### 🔨 **Build from Source**

```bash
git clone https://github.com/getporter/skeletor.git
cd skeletor
go run mage.go build install
# Binary will be in ./bin/skeletor
```

## Usage

### 📖 **Command Reference**

```bash
skeletor create [flags]
```

### 🏗️ **Core Flags**

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | **(Required)** Name of the new mixin (lowercase) | - |
| `--author` | **(Required)** Author name for the mixin | - |
| `--module` | Go module path | `github.com/getporter/<name>` |
| `--output` | Output directory | `./<name>` |
| `--compliance-level` | Compliance level (`basic`, `slsa-l1`, `slsa-l3`) | `basic` |
| `--non-interactive` | Run without prompts | `false` |
| `--dry-run` | Simulate generation without writing files | `false` |

### 🏢 **Enterprise Feature Flags**

#### 🔒 **Security Features**
```bash
--enable-security --security-features "feature1,feature2,..."
```

**Available Security Features:**
- `input_validation` - Comprehensive input sanitization and validation
- `rate_limiting` - Request throttling and abuse prevention
- `secure_headers` - HTTP security headers and CORS configuration
- `vulnerability_scanning` - Integrated security scanning with Gosec
- `policy_enforcement` - Role-based access control and security policies

#### 📋 **Compliance Frameworks**
```bash
--enable-compliance --compliance-frameworks "framework1,framework2,..."
```

**Available Compliance Frameworks:**
- `soc2` - System and Organization Controls compliance templates
- `gdpr` - General Data Protection Regulation compliance features
- `hipaa` - Health Insurance Portability and Accountability Act support
- `pci_dss` - Payment Card Industry Data Security Standard templates

#### 🔐 **Authentication & Authorization**
```bash
--enable-auth --auth-features "feature1,feature2,..."
```

**Available Auth Features:**
- `rbac` - Role-Based Access Control implementation
- `ldap` - Enterprise directory service integration
- `sso` - Single Sign-On with SAML/OAuth2/OIDC
- `mfa` - Multi-Factor Authentication implementation
- `vault` - HashiCorp Vault secrets management integration
- `session_management` - Secure session handling and lifecycle management

#### 📊 **Enhanced Observability**
```bash
--enable-observability --observability-features "feature1,feature2,..."
```

**Available Observability Features:**
- `apm` - Application Performance Monitoring setup
- `infrastructure` - System metrics and health checks
- `custom_metrics` - Business-specific metric collection
- `health_checks` - Comprehensive health endpoint implementation
- `opentelemetry` - Distributed tracing and observability
- `audit_logging` - Comprehensive audit trail and compliance logging
- `tracing` - End-to-end request tracing across services

### 🛠️ **Advanced Options**

| Flag | Description |
|------|-------------|
| `--template-url` | URL to a git repository containing a custom template |
| `--template-dir` | Local directory containing the template |
| `--var` | Set template variables in `KEY=VALUE` format (repeatable) |

## Template Variables

The following variables are used by the default template and can be provided during generation (interactively or via `--var` flag):

### 📝 **Core Variables**

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `MixinName` | string | Name of the mixin (lowercase) | **(required)** |
| `AuthorName` | string | Author name | **(required)** |
| `ModulePath` | string | Go module path | `github.com/getporter/<MixinName>` |
| `Description` | string | Short description of the mixin | Auto-generated |
| `License` | string | License (`Apache-2.0`, `MIT`, `GPL-3.0`) | `Apache-2.0` |
| `InitGit` | bool | Initialize git repository | `true` |
| `AuthorEmail` | string | Author's email for security contact | *(optional)* |

### 🔗 **Integration Variables**

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `MixinFeedRepoURL` | string | Git URL for mixin feed repository | *(optional)* |
| `MixinFeedBranch` | string | Branch in mixin feed repository | `main` |

### 🏢 **Enterprise Variables** *(Auto-populated from flags)*

| Variable | Type | Description |
|----------|------|-------------|
| `EnableSecurity` | bool | Security features enabled |
| `SecurityFeatures` | string | Comma-separated security features |
| `EnableCompliance` | bool | Compliance frameworks enabled |
| `ComplianceFrameworks` | string | Comma-separated compliance frameworks |
| `EnableAuth` | bool | Authentication features enabled |
| `AuthFeatures` | string | Comma-separated auth features |
| `EnableObservability` | bool | Observability features enabled |
| `ObservabilityFeatures` | string | Comma-separated observability features |

> **Note:** Enterprise variables are automatically populated based on the enterprise feature flags and don't need to be set manually.

## Generated Project Structure

The generated mixin project follows the standard Porter mixin structure with optional enterprise features:

### 📁 **Core Structure**

```
your-mixin/
├── cmd/your-mixin/          # CLI implementation using Cobra
├── pkg/your-mixin/          # Core mixin logic (build, install, invoke, etc.)
├── ci/main.go               # Dagger pipeline for CI/CD tasks
├── .github/workflows/       # GitHub Actions workflows
├── magefile.go              # Build automation using Mage
├── .goreleaser.yml          # GoReleaser configuration
├── Dockerfile               # Container image build
├── .golangci.yml            # Linter configuration
├── tools.go                 # Go tool dependencies
├── go.mod, go.sum           # Go module files
├── README.md                # Project documentation
├── LICENSE                  # License file
├── CONTRIBUTING.md          # Contribution guidelines
├── SECURITY.md              # Security policy
└── docs/                    # Documentation
    ├── DEVELOPER_GUIDE.md   # Developer guide
    └── OPERATIONS_GUIDE.md  # Operations guide
```

### 🏢 **Enterprise Features** *(Generated when enabled)*

#### 🔒 **Security Features** (`--enable-security`)
```
pkg/security/
├── security.go              # Core security functions
├── middleware.go            # Security middleware
└── validation.go            # Input validation

configs/security.yaml        # Security configuration
```

#### 📋 **Compliance Features** (`--enable-compliance`)
```
pkg/compliance/
└── compliance.go            # Compliance framework support

configs/compliance.yaml      # Compliance configuration
docs/COMPLIANCE_GUIDE.md     # Compliance documentation
```

#### 🔐 **Authentication Features** (`--enable-auth`)
```
pkg/auth/
├── rbac.go                  # Role-based access control
├── ldap.go                  # LDAP integration
├── sso.go                   # Single sign-on
└── vault.go                 # HashiCorp Vault integration

configs/auth.yaml            # Authentication configuration
docs/AUTH_GUIDE.md           # Authentication documentation
```

#### 📊 **Observability Features** (`--enable-observability`)
```
pkg/observability/
├── observability.go         # Enhanced monitoring
├── metrics.go               # Custom metrics
├── tracing.go               # Distributed tracing
└── audit.go                 # Audit logging

configs/observability.yaml   # Observability configuration
docs/OBSERVABILITY_GUIDE.md  # Observability documentation
```

### 🔧 **Configuration Files**

| File | Purpose | Generated When |
|------|---------|----------------|
| `.well-known/security.txt` | Security contact information | Always |
| `configs/security.yaml` | Security feature configuration | `--enable-security` |
| `configs/compliance.yaml` | Compliance framework settings | `--enable-compliance` |
| `configs/auth.yaml` | Authentication configuration | `--enable-auth` |
| `configs/observability.yaml` | Observability settings | `--enable-observability` |

## Development (Porter Mixin Generator)

This section describes how to develop the generator tool itself.

### Prerequisites

* Go 1.23+
* [Mage](https://magefile.org/)
* [Dagger CLI](https://docs.dagger.io/install)

### Building & Testing (using Dagger)

The CI/CD pipeline is defined using the Dagger Go SDK in the `./ci` directory.

To run tests and linters locally:

```bash
go run ./ci -task ci
```

### Releasing (using Dagger & GoReleaser)

Releases are handled automatically by the GitHub Actions workflow (`.github/workflows/skeletor.yml`) on tag pushes. The workflow uses Dagger to execute GoReleaser. This process includes:

* Cross-compiling binaries for Linux, macOS, and Windows (amd64/arm64).
* Generating SLSA L3 provenance attestations.
* Generating SBOMs (CycloneDX and SPDX formats) for binaries and Docker images.
* Calculating SHA256 checksums.
* Signing checksums and archives using Cosign keyless signing (via Sigstore).
* Building and pushing multi-arch Docker images for the generator tool to GHCR.
* Creating a GitHub release with all artifacts, SBOMs, signatures, and attestations attached.

To test the release process locally (requires Docker):

```bash
# Simulate a tag
export GITHUB_REF_NAME=v0.1.0-test
export GITHUB_REF_TYPE=tag

# Run the release task (requires a GITHUB_TOKEN with appropriate permissions)
# Note: This will attempt to create a real release if run outside CI context
GITHUB_TOKEN="YOUR_GITHUB_TOKEN" go run ./ci -task release
```

## Contributing

Contributions to the Porter Mixin Generator are welcome! Please see `CONTRIBUTING.md` in the *generated* mixin template for guidelines applicable to mixin development. For contributing to the generator itself, please open an issue or pull request on this repository.

## License

Apache 2.0 License
