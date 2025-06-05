# CLI Help

This page contains the command-line help for Skeletor.

## Skeletor Help

```
skeletor - Create new Porter mixins easily

Usage:
  skeletor [command]

Available Commands:
  create      Create a new Porter mixin
  completion  Generate autocompletion script
  help        Help about any command
  version     Show version information

Flags:
  -h, --help   help for skeletor

Use "skeletor [command] --help" for more information about a command.
```

## Create Command Help

```
Create a new Porter mixin with optional enterprise features.

Usage:
  skeletor create [flags]

Core Flags:
      --name string                    Name of the new mixin (lowercase) (required)
      --author string                  Author name for the mixin (required)
      --module string                  Go module path (default "github.com/getporter/<name>")
      --output string                  Output directory (default "./<name>")
      --compliance-level string        Compliance level (basic, slsa-l1, slsa-l3) (default "basic")
      --non-interactive               Run without prompts
      --dry-run                       Simulate generation without writing files

Enterprise Feature Flags:
      --enable-security               Enable enterprise security features
      --security-features string     Comma-separated list of security features
      --enable-compliance             Enable compliance framework features
      --compliance-frameworks string Comma-separated list of compliance frameworks
      --enable-auth                   Enable authentication and authorization features
      --auth-features string          Comma-separated list of auth features
      --enable-observability          Enable enhanced observability features
      --observability-features string Comma-separated list of observability features

Advanced Options:
      --template-url string           URL to a git repository containing a custom template
      --template-dir string           Local directory containing the template
      --var stringArray               Extra variables in KEY=VALUE format (repeatable)

Global Flags:
  -h, --help   help for skeletor
```

## Version Information

```
Skeletor dev
Commit: unknown
Built: unknown
```

## Enterprise Features

### Security Features
- `input_validation` - Comprehensive input sanitization and validation
- `rate_limiting` - Request throttling and abuse prevention
- `secure_headers` - HTTP security headers and CORS configuration
- `vulnerability_scanning` - Integrated security scanning with Gosec
- `policy_enforcement` - Role-based access control and security policies

### Compliance Frameworks
- `soc2` - System and Organization Controls compliance templates
- `gdpr` - General Data Protection Regulation compliance features
- `hipaa` - Health Insurance Portability and Accountability Act support
- `pci_dss` - Payment Card Industry Data Security Standard templates

### Authentication Features
- `rbac` - Role-Based Access Control implementation
- `ldap` - Enterprise directory service integration
- `sso` - Single Sign-On with SAML/OAuth2/OIDC
- `mfa` - Multi-Factor Authentication implementation
- `vault` - HashiCorp Vault secrets management integration
- `session_management` - Secure session handling and lifecycle management

### Observability Features
- `apm` - Application Performance Monitoring setup
- `infrastructure` - System metrics and health checks
- `custom_metrics` - Business-specific metric collection
- `health_checks` - Comprehensive health endpoint implementation
- `opentelemetry` - Distributed tracing and observability
- `audit_logging` - Comprehensive audit trail and compliance logging
- `tracing` - End-to-end request tracing across services

## Examples

### Basic Mixin
```bash
skeletor create --name my-mixin --author "John Doe"
```

### Enterprise Mixin
```bash
skeletor create \
  --name enterprise-mixin \
  --author "Enterprise Team" \
  --enable-security \
  --security-features "input_validation,rate_limiting" \
  --enable-observability \
  --observability-features "apm,opentelemetry"
```

### Full Enterprise Mixin
```bash
skeletor create \
  --name full-enterprise \
  --author "Enterprise Team" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,policy_enforcement" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr" \
  --enable-auth \
  --auth-features "rbac,sso,mfa" \
  --enable-observability \
  --observability-features "apm,opentelemetry,audit_logging,tracing"
```
