# Command Reference

Complete reference for all Skeletor commands and options.

## Global Commands

### `skeletor`

Main command for the Porter Mixin Generator.

```bash
skeletor [command]
```

**Available Commands:**
* `create` - Create a new Porter mixin
* `completion` - Generate autocompletion script
* `help` - Help about any command

**Global Flags:**
* `-h, --help` - Help for skeletor

## Create Command

### `skeletor create`

Create a new Porter mixin with optional enterprise features.

```bash
skeletor create [flags]
```

### Core Flags

| Flag | Type | Description | Default | Required |
|------|------|-------------|---------|----------|
| `--name` | string | Name of the new mixin (lowercase) | - | ✅ |
| `--author` | string | Author name for the mixin | - | ✅ |
| `--module` | string | Go module path | `github.com/getporter/<name>` | ❌ |
| `--output` | string | Output directory | `./<name>` | ❌ |
| `--compliance-level` | string | Compliance level (`basic`, `slsa-l1`, `slsa-l3`) | `basic` | ❌ |
| `--non-interactive` | bool | Run without prompts | `false` | ❌ |
| `--dry-run` | bool | Simulate generation without writing files | `false` | ❌ |

### Enterprise Feature Flags

#### Security Features

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--enable-security` | bool | Enable enterprise security features | `false` |
| `--security-features` | string | Comma-separated list of security features | `""` |

**Available Security Features:**
* `input_validation` - Comprehensive input sanitization and validation
* `rate_limiting` - Request throttling and abuse prevention
* `secure_headers` - HTTP security headers and CORS configuration
* `vulnerability_scanning` - Integrated security scanning with Gosec
* `policy_enforcement` - Role-based access control and security policies

**Example:**
```bash
skeletor create \
  --name secure-mixin \
  --author "Security Team" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers"
```

#### Compliance Frameworks

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--enable-compliance` | bool | Enable compliance framework features | `false` |
| `--compliance-frameworks` | string | Comma-separated list of compliance frameworks | `""` |

**Available Compliance Frameworks:**
* `soc2` - System and Organization Controls compliance templates
* `gdpr` - General Data Protection Regulation compliance features
* `hipaa` - Health Insurance Portability and Accountability Act support
* `pci_dss` - Payment Card Industry Data Security Standard templates

**Example:**
```bash
skeletor create \
  --name compliant-mixin \
  --author "Compliance Team" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa"
```

#### Authentication & Authorization

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--enable-auth` | bool | Enable authentication and authorization features | `false` |
| `--auth-features` | string | Comma-separated list of auth features | `""` |

**Available Auth Features:**
* `rbac` - Role-Based Access Control implementation
* `ldap` - Enterprise directory service integration
* `sso` - Single Sign-On with SAML/OAuth2/OIDC
* `mfa` - Multi-Factor Authentication implementation
* `vault` - HashiCorp Vault secrets management integration
* `session_management` - Secure session handling and lifecycle management

**Example:**
```bash
skeletor create \
  --name auth-mixin \
  --author "Auth Team" \
  --enable-auth \
  --auth-features "rbac,ldap,sso,mfa,vault"
```

#### Enhanced Observability

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--enable-observability` | bool | Enable enhanced observability features | `false` |
| `--observability-features` | string | Comma-separated list of observability features | `""` |

**Available Observability Features:**
* `apm` - Application Performance Monitoring setup
* `infrastructure` - System metrics and health checks
* `custom_metrics` - Business-specific metric collection
* `health_checks` - Comprehensive health endpoint implementation
* `opentelemetry` - Distributed tracing and observability
* `audit_logging` - Comprehensive audit trail and compliance logging
* `tracing` - End-to-end request tracing across services

**Example:**
```bash
skeletor create \
  --name observable-mixin \
  --author "SRE Team" \
  --enable-observability \
  --observability-features "apm,opentelemetry,audit_logging,tracing"
```

### Advanced Options

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--template-url` | string | URL to a git repository containing a custom template | `""` |
| `--template-dir` | string | Local directory containing the template | `""` |
| `--var` | stringArray | Extra variables in `KEY=VALUE` format (repeatable) | `[]` |

**Example:**
```bash
skeletor create \
  --name custom-mixin \
  --author "Custom Team" \
  --template-url "https://github.com/your-org/custom-template.git" \
  --var "CustomVar=value" \
  --var "AnotherVar=another-value"
```

## Usage Examples

### Basic Mixin

Create a simple mixin with default settings:

```bash
skeletor create --name my-mixin --author "John Doe"
```

### Security-Focused Mixin

Create a mixin with comprehensive security features:

```bash
skeletor create \
  --name secure-mixin \
  --author "Security Team" \
  --module "github.com/security-org/secure-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement" \
  --compliance-level "slsa-l3"
```

### Compliance-Ready Mixin

Create a mixin with multiple compliance frameworks:

```bash
skeletor create \
  --name compliant-mixin \
  --author "Compliance Team" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa,pci_dss" \
  --enable-security \
  --security-features "input_validation,policy_enforcement"
```

### Full Enterprise Mixin

Create a mixin with all enterprise features enabled:

```bash
skeletor create \
  --name enterprise-mixin \
  --author "Enterprise Team" \
  --module "github.com/enterprise-org/enterprise-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa" \
  --enable-auth \
  --auth-features "rbac,ldap,sso,mfa,vault,session_management" \
  --enable-observability \
  --observability-features "apm,infrastructure,custom_metrics,health_checks,opentelemetry,audit_logging,tracing" \
  --compliance-level "slsa-l3"
```

### Non-Interactive Mode

Create a mixin without prompts (useful for CI/CD):

```bash
skeletor create \
  --name ci-mixin \
  --author "CI System" \
  --module "github.com/ci-org/ci-mixin" \
  --output "./generated-mixin" \
  --non-interactive \
  --enable-security \
  --security-features "input_validation,rate_limiting"
```

### Dry Run

Preview what would be generated without creating files:

```bash
skeletor create \
  --name preview-mixin \
  --author "Preview User" \
  --enable-security \
  --security-features "input_validation" \
  --dry-run
```

## Completion Command

### `skeletor completion`

Generate autocompletion scripts for various shells.

```bash
skeletor completion [bash|zsh|fish|powershell]
```

**Examples:**

```bash
# Bash
source <(skeletor completion bash)

# Zsh
source <(skeletor completion zsh)

# Fish
skeletor completion fish | source

# PowerShell
skeletor completion powershell | Out-String | Invoke-Expression
```

## Help Command

### `skeletor help`

Get help for any command.

```bash
skeletor help [command]
```

**Examples:**

```bash
# General help
skeletor help

# Help for create command
skeletor help create

# Help for completion command
skeletor help completion
```

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | General error |
| `2` | Invalid arguments |
| `3` | Template processing error |
| `4` | File system error |
| `5` | Validation error |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SKELETOR_TEMPLATE_URL` | Default template URL | Built-in templates |
| `SKELETOR_OUTPUT_DIR` | Default output directory | Current directory |
| `SKELETOR_NON_INTERACTIVE` | Run in non-interactive mode | `false` |

## Configuration Files

Skeletor doesn't use configuration files by default, but you can create shell aliases or scripts for common configurations:

```bash
# ~/.bashrc or ~/.zshrc
alias skeletor-enterprise='skeletor create --enable-security --security-features "input_validation,rate_limiting,secure_headers" --enable-compliance --compliance-frameworks "soc2,gdpr" --enable-auth --auth-features "rbac,sso" --enable-observability --observability-features "apm,opentelemetry,audit_logging"'

# Usage
skeletor-enterprise --name my-enterprise-mixin --author "Enterprise Team"
```

## Next Steps

* [Enterprise Features Guide](enterprise-features.md) - Detailed enterprise features documentation
* [Template Customization](template-customization.md) - How to customize templates
* [Examples](examples.md) - Real-world usage examples
