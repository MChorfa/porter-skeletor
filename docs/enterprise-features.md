# Enterprise Features Guide

Skeletor provides comprehensive enterprise-grade features that can be enabled through CLI flags. This guide covers all available enterprise features and how to use them effectively.

## Overview

Enterprise features are organized into four main categories:

* **🔒 Security Features** - Input validation, rate limiting, secure headers, vulnerability scanning
* **📋 Compliance Frameworks** - SOC 2, GDPR, HIPAA, PCI DSS compliance templates
* **🔐 Authentication & Authorization** - RBAC, LDAP, SSO, MFA, HashiCorp Vault integration
* **📊 Enhanced Observability** - APM, OpenTelemetry, audit logging, distributed tracing

## Security Features

Enable security features with `--enable-security` and specify features with `--security-features`:

```bash
skeletor create \
  --name secure-mixin \
  --author "Security Team" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement"
```

### Available Security Features

#### Input Validation (`input_validation`)

Generates comprehensive input sanitization and validation:

* **Generated Files:**
  * `pkg/security/validation.go` - Input validation functions
  * `configs/security.yaml` - Validation rules configuration

* **Features:**
  * Schema-based validation
  * Input sanitization
  * Type checking
  * Length limits
  * Pattern matching

#### Rate Limiting (`rate_limiting`)

Implements request throttling and abuse prevention:

* **Generated Files:**
  * `pkg/security/middleware.go` - Rate limiting middleware
  * `configs/security.yaml` - Rate limiting configuration

* **Features:**
  * Configurable rate limits
  * Burst handling
  * IP-based limiting
  * User-based limiting
  * Redis backend support

#### Secure Headers (`secure_headers`)

Adds HTTP security headers and CORS configuration:

* **Generated Files:**
  * `pkg/security/headers.go` - Security headers middleware
  * `configs/security.yaml` - Headers configuration

* **Features:**
  * Content Security Policy (CSP)
  * HTTP Strict Transport Security (HSTS)
  * X-Frame-Options
  * X-Content-Type-Options
  * CORS configuration

#### Vulnerability Scanning (`vulnerability_scanning`)

Integrates security scanning with Gosec:

* **Generated Files:**
  * `.github/workflows/security.yml` - Security scanning workflow
  * `configs/security.yaml` - Scanning configuration

* **Features:**
  * Automated vulnerability scanning
  * Dependency checking
  * Code analysis
  * Security reporting

#### Policy Enforcement (`policy_enforcement`)

Implements role-based access control and security policies:

* **Generated Files:**
  * `pkg/security/policy.go` - Policy enforcement engine
  * `configs/security.yaml` - Policy definitions

* **Features:**
  * Role-based access control
  * Resource-based permissions
  * Policy evaluation engine
  * Audit logging

## Compliance Frameworks

Enable compliance features with `--enable-compliance` and specify frameworks with `--compliance-frameworks`:

```bash
skeletor create \
  --name compliant-mixin \
  --author "Compliance Team" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa,pci_dss"
```

### Available Compliance Frameworks

#### SOC 2 (`soc2`)

System and Organization Controls compliance templates:

* **Generated Files:**
  * `pkg/compliance/soc2.go` - SOC 2 compliance functions
  * `configs/compliance.yaml` - SOC 2 configuration
  * `docs/SOC2_COMPLIANCE.md` - SOC 2 documentation

* **Features:**
  * Security controls
  * Availability controls
  * Processing integrity
  * Confidentiality controls
  * Privacy controls

#### GDPR (`gdpr`)

General Data Protection Regulation compliance features:

* **Generated Files:**
  * `pkg/compliance/gdpr.go` - GDPR compliance functions
  * `configs/compliance.yaml` - GDPR configuration
  * `docs/GDPR_COMPLIANCE.md` - GDPR documentation

* **Features:**
  * Data protection by design
  * Consent management
  * Data subject rights
  * Data breach notification
  * Privacy impact assessments

#### HIPAA (`hipaa`)

Health Insurance Portability and Accountability Act support:

* **Generated Files:**
  * `pkg/compliance/hipaa.go` - HIPAA compliance functions
  * `configs/compliance.yaml` - HIPAA configuration
  * `docs/HIPAA_COMPLIANCE.md` - HIPAA documentation

* **Features:**
  * PHI protection
  * Access controls
  * Audit logging
  * Encryption requirements
  * Business associate agreements

#### PCI DSS (`pci_dss`)

Payment Card Industry Data Security Standard templates:

* **Generated Files:**
  * `pkg/compliance/pci.go` - PCI DSS compliance functions
  * `configs/compliance.yaml` - PCI DSS configuration
  * `docs/PCI_COMPLIANCE.md` - PCI DSS documentation

* **Features:**
  * Cardholder data protection
  * Network security
  * Access controls
  * Monitoring and testing
  * Information security policies

## Authentication & Authorization

Enable authentication features with `--enable-auth` and specify features with `--auth-features`:

```bash
skeletor create \
  --name auth-mixin \
  --author "Auth Team" \
  --enable-auth \
  --auth-features "rbac,ldap,sso,mfa,vault,session_management"
```

### Available Auth Features

#### RBAC (`rbac`)

Role-Based Access Control implementation:

* **Generated Files:**
  * `pkg/auth/rbac.go` - RBAC implementation
  * `configs/auth.yaml` - RBAC configuration

* **Features:**
  * Role definitions
  * Permission management
  * User-role assignments
  * Resource-based access control

#### LDAP Integration (`ldap`)

Enterprise directory service integration:

* **Generated Files:**
  * `pkg/auth/ldap.go` - LDAP integration
  * `configs/auth.yaml` - LDAP configuration

* **Features:**
  * LDAP authentication
  * User synchronization
  * Group mapping
  * Directory queries

#### Single Sign-On (`sso`)

SSO with SAML/OAuth2/OIDC:

* **Generated Files:**
  * `pkg/auth/sso.go` - SSO implementation
  * `configs/auth.yaml` - SSO configuration

* **Features:**
  * SAML 2.0 support
  * OAuth2/OIDC integration
  * Identity provider integration
  * Token management

#### Multi-Factor Authentication (`mfa`)

MFA implementation:

* **Generated Files:**
  * `pkg/auth/mfa.go` - MFA implementation
  * `configs/auth.yaml` - MFA configuration

* **Features:**
  * TOTP support
  * SMS verification
  * Email verification
  * Hardware token support

#### HashiCorp Vault (`vault`)

Secrets management integration:

* **Generated Files:**
  * `pkg/auth/vault.go` - Vault integration
  * `configs/auth.yaml` - Vault configuration

* **Features:**
  * Secret storage
  * Dynamic secrets
  * Encryption as a service
  * PKI management

#### Session Management (`session_management`)

Secure session handling:

* **Generated Files:**
  * `pkg/auth/session.go` - Session management
  * `configs/auth.yaml` - Session configuration

* **Features:**
  * Session creation
  * Session validation
  * Session expiration
  * Session storage

## Enhanced Observability

Enable observability features with `--enable-observability` and specify features with `--observability-features`:

```bash
skeletor create \
  --name observable-mixin \
  --author "SRE Team" \
  --enable-observability \
  --observability-features "apm,infrastructure,custom_metrics,health_checks,opentelemetry,audit_logging,tracing"
```

### Available Observability Features

#### APM Integration (`apm`)

Application Performance Monitoring setup:

* **Generated Files:**
  * `pkg/observability/apm.go` - APM integration
  * `configs/observability.yaml` - APM configuration

* **Features:**
  * Performance monitoring
  * Error tracking
  * Transaction tracing
  * Service maps

#### Infrastructure Monitoring (`infrastructure`)

System metrics and health checks:

* **Generated Files:**
  * `pkg/observability/infrastructure.go` - Infrastructure monitoring
  * `configs/observability.yaml` - Infrastructure configuration

* **Features:**
  * System metrics
  * Resource monitoring
  * Health checks
  * Alerting

#### Custom Metrics (`custom_metrics`)

Business-specific metric collection:

* **Generated Files:**
  * `pkg/observability/metrics.go` - Custom metrics
  * `configs/observability.yaml` - Metrics configuration

* **Features:**
  * Custom metric definitions
  * Metric collection
  * Metric aggregation
  * Metric export

#### Health Checks (`health_checks`)

Comprehensive health endpoint implementation:

* **Generated Files:**
  * `pkg/observability/health.go` - Health checks
  * `configs/observability.yaml` - Health configuration

* **Features:**
  * Liveness probes
  * Readiness probes
  * Dependency checks
  * Health reporting

#### OpenTelemetry (`opentelemetry`)

Distributed tracing and observability:

* **Generated Files:**
  * `pkg/observability/otel.go` - OpenTelemetry integration
  * `configs/observability.yaml` - OpenTelemetry configuration

* **Features:**
  * Distributed tracing
  * Metrics collection
  * Log correlation
  * Span management

#### Audit Logging (`audit_logging`)

Comprehensive audit trail:

* **Generated Files:**
  * `pkg/observability/audit.go` - Audit logging
  * `configs/observability.yaml` - Audit configuration

* **Features:**
  * Event logging
  * User activity tracking
  * Compliance logging
  * Log retention

#### Distributed Tracing (`tracing`)

End-to-end request tracing:

* **Generated Files:**
  * `pkg/observability/tracing.go` - Distributed tracing
  * `configs/observability.yaml` - Tracing configuration

* **Features:**
  * Request tracing
  * Service correlation
  * Performance analysis
  * Error tracking

## Configuration

All enterprise features generate configuration files in the `configs/` directory:

* `configs/security.yaml` - Security feature configuration
* `configs/compliance.yaml` - Compliance framework settings
* `configs/auth.yaml` - Authentication configuration
* `configs/observability.yaml` - Observability settings

These configuration files provide extensive customization options for each feature and can be modified after generation to meet specific requirements.

## Best Practices

1. **Start Small:** Begin with basic features and add more as needed
2. **Review Generated Code:** Customize the generated templates for your specific use case
3. **Test Thoroughly:** Validate all enterprise features in your environment
4. **Monitor Performance:** Ensure enterprise features don't impact performance
5. **Keep Updated:** Regularly update Skeletor to get the latest enterprise features

## Next Steps

* [Command Reference](command-reference.md) - Complete CLI documentation
* [Template Customization](template-customization.md) - How to customize templates
* [Examples](examples.md) - Real-world usage examples
