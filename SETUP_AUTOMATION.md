# Setup Automation Guide

This guide helps you set up all the automation features for your Porter Skeletor repository.

## 🚀 Quick Setup Checklist

### 1. Enable GitHub Pages

1. Go to your repository: https://github.com/MChorfa/porter-skeletor
2. Click **Settings** → **Pages**
3. Under **Source**, select **Deploy from a branch**
4. Select branch: **main**
5. Select folder: **/ (root)**
6. Click **Save**

The documentation site will be available at: https://mchorfa.github.io/porter-skeletor/

### 2. Set Up Repository Secrets

For automated Homebrew formula updates, add these secrets:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Add the following secrets:

| Secret Name | Description | Required For |
|-------------|-------------|--------------|
| `HOMEBREW_TAP_TOKEN` | GitHub token with access to homebrew-tap repo | Homebrew automation |
| `COSIGN_PRIVATE_KEY` | Cosign private key for artifact signing | Release signing |
| `COSIGN_PASSWORD` | Password for Cosign private key | Release signing |

### 3. Create Homebrew Tap Repository

1. Create a new repository: `MChorfa/homebrew-tap`
2. Initialize with a README
3. This will be used for automated Homebrew formula updates

### 4. Test the Automation

#### Test Documentation Deployment

1. Make any change to a file in the `docs/` directory
2. Commit and push the change
3. Check the **Actions** tab to see the documentation deployment
4. Visit your GitHub Pages site to see the updated documentation

#### Test Release Automation

1. Create and push a git tag:
   ```bash
   git tag v1.0.0
   git push mchorfa v1.0.0
   ```
2. Check the **Actions** tab for the release workflow
3. Check the **Releases** page for the new release

## 📋 Available Workflows

### 1. Documentation Deployment (`.github/workflows/pages.yml`)

**Triggers:**
- Push to `main` branch (when docs files change)
- New releases
- Manual dispatch

**What it does:**
- Builds the MkDocs documentation site
- Generates CLI help documentation
- Deploys to GitHub Pages

### 2. Homebrew Formula Updates (`.github/workflows/homebrew.yml`)

**Triggers:**
- New releases
- Manual dispatch

**What it does:**
- Updates the Homebrew formula with new version
- Tests the formula installation
- Creates PR to homebrew-tap repository

### 3. CI/CD Pipeline (`.github/workflows/skeletor.yml`)

**Triggers:**
- Push to any branch
- Pull requests

**What it does:**
- Runs tests
- Builds binaries
- Validates enterprise features

## 🔧 Configuration Files

### GoReleaser (`.goreleaser.yml`)

Configured for:
- Multi-platform builds (Linux, macOS, Windows)
- Docker image publishing
- SLSA provenance generation
- SBOM generation
- Cosign signing
- Homebrew formula generation

### MkDocs (`mkdocs.yml`)

Configured for:
- Material theme with dark/light mode
- Comprehensive navigation
- Search functionality
- Code highlighting

## 🧪 Testing the Enterprise Features

### Basic Test

```bash
# Clone your repository
git clone https://github.com/MChorfa/porter-skeletor.git
cd porter-skeletor

# Build the binary
go build -o skeletor ./cmd/skeletor

# Test basic functionality
./skeletor --help
./skeletor version
```

### Enterprise Features Test

```bash
# Test enterprise features
./skeletor create \
  --name test-enterprise \
  --author "Test User" \
  --module "github.com/test/test-enterprise" \
  --enable-security \
  --security-features "input_validation,rate_limiting" \
  --enable-compliance \
  --compliance-frameworks "soc2" \
  --enable-auth \
  --auth-features "rbac" \
  --enable-observability \
  --observability-features "apm" \
  --dry-run \
  --non-interactive
```

### Full Enterprise Test

```bash
# Test all enterprise features
./skeletor create \
  --name full-enterprise \
  --author "Enterprise Team" \
  --module "github.com/enterprise/full-test" \
  --enable-security \
  --security-features "input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement" \
  --enable-compliance \
  --compliance-frameworks "soc2,gdpr,hipaa" \
  --enable-auth \
  --auth-features "rbac,ldap,sso,mfa,vault,session_management" \
  --enable-observability \
  --observability-features "apm,infrastructure,custom_metrics,health_checks,opentelemetry,audit_logging,tracing" \
  --compliance-level "slsa-l3" \
  --non-interactive
```

## 📦 Installation Methods

Once automation is set up, users can install Skeletor using:

### Homebrew (after tap is set up)

```bash
brew tap MChorfa/tap
brew install skeletor
```

### Go Install

```bash
go install github.com/MChorfa/porter-skeletor/cmd/skeletor@latest
```

### Binary Downloads

Download from: https://github.com/MChorfa/porter-skeletor/releases

### Docker

```bash
docker pull ghcr.io/mchorfa/porter-skeletor:latest
```

## 🔍 Monitoring

### GitHub Actions

Monitor workflow runs at: https://github.com/MChorfa/porter-skeletor/actions

### Documentation Site

Monitor site deployment at: https://mchorfa.github.io/porter-skeletor/

### Releases

Monitor releases at: https://github.com/MChorfa/porter-skeletor/releases

## 🐛 Troubleshooting

### Documentation Not Deploying

1. Check GitHub Pages settings
2. Verify the workflow ran successfully
3. Check for build errors in Actions tab

### Homebrew Formula Not Updating

1. Verify `HOMEBREW_TAP_TOKEN` secret is set
2. Check that homebrew-tap repository exists
3. Review workflow logs for errors

### Release Workflow Failing

1. Check that the tag follows semantic versioning (v1.0.0)
2. Verify all tests pass before tagging
3. Check for missing secrets or permissions

## 📞 Support

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/MChorfa/porter-skeletor/issues)
2. Review workflow logs in the Actions tab
3. Verify all secrets and settings are configured correctly

## 🎯 Next Steps

1. **Enable GitHub Pages** (most important)
2. **Set up repository secrets** for full automation
3. **Create your first release** to test the pipeline
4. **Customize the documentation** for your specific needs
5. **Share with your team** and start generating enterprise mixins!

Your Porter Skeletor is now ready for enterprise use with full automation! 🚀
