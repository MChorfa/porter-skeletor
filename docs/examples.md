# Examples

This page provides practical examples of using Skeletor to generate Porter mixins for various use cases.

## Basic Examples

### Simple Mixin

Create a basic mixin with minimal configuration:

```bash
skeletor create \
  --name hello-world \
  --author "John Doe" \
  --module "github.com/johndoe/hello-world-mixin"
```

**Generated Structure:**
```
hello-world/
├── cmd/hello-world/
├── pkg/hello-world/
├── go.mod
├── README.md
├── LICENSE
└── ...
```

### Mixin with Custom Variables

Create a mixin with custom template variables:

```bash
skeletor create \
  --name custom-mixin \
  --author "Jane Smith" \
  --module "github.com/janesmith/custom-mixin" \
  --var "Description=A custom Porter mixin for special operations" \
  --var "License=MIT" \
  --var "AuthorEmail=jane@example.com"
```

### Non-Interactive Mode

Create a mixin in CI/CD environments:

```bash
skeletor create \
  --name ci-mixin \
  --author "CI System" \
  --module "github.com/company/ci-mixin" \
  --output "./generated-mixins/ci-mixin" \
  --non-interactive \
  --compliance-level "slsa-l1"
```

## Security-Focused Examples

### Basic Security Mixin

Create a mixin with essential security features:

```bash
skeletor create \
  --name secure-mixin \
  --author "Security Team" \
  --module "github.com/security-org/secure-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting"
```

**Generated Security Features:**
- Input validation middleware
- Rate limiting configuration
- Security configuration file

### Advanced Security Mixin

Create a mixin with comprehensive security features:

```bash
skeletor create \
  --name fortress-mixin \
  --author "Security Team" \
  --module "github.com/security-org/fortress-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement" \
  --compliance-level "slsa-l3"
```

**Generated Security Features:**
- Comprehensive input validation
- Advanced rate limiting
- Security headers middleware
- Vulnerability scanning integration
- Policy enforcement engine
- SLSA Level 3 compliance

## Compliance Examples

### GDPR Compliant Mixin

Create a mixin for GDPR compliance:

```bash
skeletor create \
  --name gdpr-mixin \
  --author "Compliance Team" \
  --module "github.com/company/gdpr-mixin" \
  --enable-compliance \
  --compliance-frameworks "gdpr" \
  --enable-security \
  --security-features "input_validation,policy_enforcement"
```

### Healthcare Mixin (HIPAA)

Create a mixin for healthcare applications:

```bash
skeletor create \
  --name healthcare-mixin \
  --author "Healthcare IT" \
  --module "github.com/healthcare-org/healthcare-mixin" \
  --enable-compliance \
  --compliance-frameworks "hipaa,soc2" \
  --enable-security \
  --security-features "input_validation,secure_headers,policy_enforcement" \
  --enable-auth \
  --auth-features "rbac,mfa" \
  --compliance-level "slsa-l3"
```

### Financial Services Mixin

Create a mixin for financial services:

```bash
skeletor create \
  --name fintech-mixin \
  --author "FinTech Team" \
  --module "github.com/fintech-org/fintech-mixin" \
  --enable-compliance \
  --compliance-frameworks "soc2,pci_dss" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement" \
  --enable-auth \
  --auth-features "rbac,mfa,vault" \
  --enable-observability \
  --observability-features "audit_logging,apm" \
  --compliance-level "slsa-l3"
```

## Enterprise Examples

### Full Enterprise Mixin

Create a mixin with all enterprise features:

```bash
skeletor create \
  --name enterprise-platform \
  --author "Platform Team" \
  --module "github.com/enterprise/platform-mixin" \
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

### Microservices Platform Mixin

Create a mixin for microservices platforms:

```bash
skeletor create \
  --name microservices-mixin \
  --author "Platform Engineering" \
  --module "github.com/platform/microservices-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers" \
  --enable-auth \
  --auth-features "rbac,sso" \
  --enable-observability \
  --observability-features "apm,opentelemetry,tracing,health_checks" \
  --compliance-level "slsa-l1"
```

### Cloud-Native Mixin

Create a mixin for cloud-native applications:

```bash
skeletor create \
  --name cloud-native-mixin \
  --author "Cloud Team" \
  --module "github.com/cloud-org/cloud-native-mixin" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers" \
  --enable-auth \
  --auth-features "rbac,vault" \
  --enable-observability \
  --observability-features "apm,infrastructure,opentelemetry,health_checks" \
  --compliance-level "slsa-l1"
```

## Custom Template Examples

### Using Custom Template Repository

Create a mixin using a custom template:

```bash
skeletor create \
  --name custom-template-mixin \
  --author "Custom Team" \
  --template-url "https://github.com/your-org/custom-mixin-template.git" \
  --var "CustomFeature=enabled"
```

### Using Local Template Directory

Create a mixin using a local template:

```bash
skeletor create \
  --name local-template-mixin \
  --author "Local Team" \
  --template-dir "./my-custom-template" \
  --var "LocalFeature=enabled"
```

## Testing Examples

### Dry Run

Preview what would be generated without creating files:

```bash
skeletor create \
  --name preview-mixin \
  --author "Preview User" \
  --enable-security \
  --security-features "input_validation,rate_limiting" \
  --enable-observability \
  --observability-features "apm,health_checks" \
  --dry-run
```

### Validation Testing

Create a mixin and validate it builds correctly:

```bash
# Create the mixin
skeletor create \
  --name test-mixin \
  --author "Test User" \
  --module "github.com/test/test-mixin" \
  --enable-security \
  --security-features "input_validation" \
  --non-interactive

# Navigate to the generated directory
cd test-mixin

# Validate the mixin builds
go mod tidy
go build ./...
go test ./...

# Run mixin-specific tests if available
if [ -f magefile.go ]; then
  go run mage.go test
fi
```

## Docker Examples

### Using Docker

Create a mixin using the Docker image:

```bash
# Create a basic mixin
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create --name docker-mixin --author "Docker User" --non-interactive

# Create an enterprise mixin
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create \
  --name enterprise-docker-mixin \
  --author "Enterprise Docker User" \
  --enable-security \
  --security-features "input_validation,rate_limiting" \
  --enable-observability \
  --observability-features "apm,health_checks" \
  --non-interactive
```

### Docker Compose for Development

Create a `docker-compose.yml` for mixin development:

```yaml
version: '3.8'
services:
  skeletor:
    image: ghcr.io/getporter/skeletor:latest
    volumes:
      - .:/work
    working_dir: /work
    command: >
      create
      --name dev-mixin
      --author "Dev Team"
      --enable-security
      --security-features "input_validation,rate_limiting"
      --enable-observability
      --observability-features "apm,health_checks"
      --non-interactive
```

## CI/CD Integration Examples

### GitHub Actions

Create a GitHub Actions workflow to generate mixins:

```yaml
name: Generate Mixin
on:
  workflow_dispatch:
    inputs:
      mixin_name:
        description: 'Name of the mixin to generate'
        required: true
      author:
        description: 'Author name'
        required: true

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Install Skeletor
        run: |
          curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_linux_amd64.tar.gz | tar xz
          sudo mv skeletor /usr/local/bin/
      
      - name: Generate Mixin
        run: |
          skeletor create \
            --name "${{ github.event.inputs.mixin_name }}" \
            --author "${{ github.event.inputs.author }}" \
            --module "github.com/${{ github.repository_owner }}/${{ github.event.inputs.mixin_name }}-mixin" \
            --enable-security \
            --security-features "input_validation,rate_limiting" \
            --enable-observability \
            --observability-features "apm,health_checks" \
            --non-interactive
      
      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v6
        with:
          title: "Add ${{ github.event.inputs.mixin_name }} mixin"
          body: "Generated mixin using Skeletor"
```

### GitLab CI

Create a GitLab CI pipeline:

```yaml
generate_mixin:
  image: ghcr.io/getporter/skeletor:latest
  script:
    - |
      skeletor create \
        --name "$MIXIN_NAME" \
        --author "$AUTHOR_NAME" \
        --module "gitlab.com/$CI_PROJECT_NAMESPACE/$MIXIN_NAME-mixin" \
        --enable-security \
        --security-features "input_validation,rate_limiting" \
        --non-interactive
  artifacts:
    paths:
      - "$MIXIN_NAME/"
  only:
    - web
```

## Next Steps

After generating your mixin:

1. **Review Generated Code**: Examine the generated files and customize them for your specific use case
2. **Configure Features**: Update the configuration files in the `configs/` directory
3. **Implement Business Logic**: Add your mixin-specific functionality to the generated skeleton
4. **Test Thoroughly**: Run tests and validate the mixin works as expected
5. **Deploy**: Use the generated CI/CD pipeline to deploy your mixin

For more detailed information, see:
- [Enterprise Features Guide](enterprise-features.md)
- [Command Reference](command-reference.md)
- [Template Customization](template-customization.md)
