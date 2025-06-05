# Installation Guide

This guide covers all the ways to install Skeletor on your system.

## Package Managers

### Homebrew (macOS/Linux)

The easiest way to install Skeletor on macOS and Linux:

```bash
# Add the Porter tap
brew tap getporter/tap

# Install Skeletor
brew install skeletor

# Verify installation
skeletor --help
```

**Update Skeletor:**
```bash
brew update
brew upgrade skeletor
```

### Go Install

If you have Go installed, you can install Skeletor directly:

```bash
# Install latest version
go install github.com/getporter/skeletor/cmd/skeletor@latest

# Install specific version
go install github.com/getporter/skeletor/cmd/skeletor@v1.0.0

# Verify installation
skeletor --help
```

**Requirements:**
- Go 1.21 or later
- Git (for template cloning)

## Binary Downloads

Download pre-built binaries from the [GitHub releases page](https://github.com/getporter/skeletor/releases).

### Linux

```bash
# Download and install (amd64)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_linux_amd64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# Download and install (arm64)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_linux_arm64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# Verify installation
skeletor --help
```

### macOS

```bash
# Download and install (amd64)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_darwin_amd64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# Download and install (arm64 - Apple Silicon)
curl -L https://github.com/getporter/skeletor/releases/latest/download/skeletor_darwin_arm64.tar.gz | tar xz
sudo mv skeletor /usr/local/bin/

# Verify installation
skeletor --help
```

### Windows

1. Download the appropriate archive from the [releases page](https://github.com/getporter/skeletor/releases):
   - `skeletor_windows_amd64.zip` (64-bit)
   - `skeletor_windows_arm64.zip` (ARM64)

2. Extract the archive to a directory in your PATH (e.g., `C:\Program Files\Skeletor\`)

3. Add the directory to your PATH environment variable

4. Open a new command prompt and verify:
   ```cmd
   skeletor --help
   ```

## Docker

Use the official Docker image for containerized environments:

```bash
# Pull the latest image
docker pull ghcr.io/getporter/skeletor:latest

# Verify the image
docker run --rm ghcr.io/getporter/skeletor:latest --help

# Create a mixin using Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create --name my-mixin --author "Your Name" --non-interactive
```

**Available Tags:**
- `latest` - Latest stable release
- `v1.0.0` - Specific version
- `main` - Latest development build

## Build from Source

For development or custom builds:

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for convenience)

### Build Steps

```bash
# Clone the repository
git clone https://github.com/getporter/skeletor.git
cd skeletor

# Build using Go
go build -o bin/skeletor ./cmd/skeletor

# Or build using Mage (if available)
go run mage.go build

# Install to GOPATH/bin
go install ./cmd/skeletor

# Verify installation
./bin/skeletor --help
```

### Development Build

For development with additional debugging:

```bash
# Build with debug information
go build -gcflags="all=-N -l" -o bin/skeletor-debug ./cmd/skeletor

# Build with race detection
go build -race -o bin/skeletor-race ./cmd/skeletor
```

## Verification

After installation, verify Skeletor is working correctly:

```bash
# Check version
skeletor version

# Check help
skeletor --help

# Test basic functionality
skeletor create --help

# Test enterprise features
skeletor create \
  --name test-mixin \
  --author "Test User" \
  --enable-security \
  --security-features "input_validation" \
  --dry-run \
  --non-interactive
```

## Shell Completion

Enable shell completion for better CLI experience:

### Bash

```bash
# Add to ~/.bashrc
echo 'source <(skeletor completion bash)' >> ~/.bashrc

# Or install system-wide
skeletor completion bash | sudo tee /etc/bash_completion.d/skeletor
```

### Zsh

```bash
# Add to ~/.zshrc
echo 'source <(skeletor completion zsh)' >> ~/.zshrc

# Or for Oh My Zsh
skeletor completion zsh > ~/.oh-my-zsh/completions/_skeletor
```

### Fish

```bash
skeletor completion fish | source

# Or save permanently
skeletor completion fish > ~/.config/fish/completions/skeletor.fish
```

### PowerShell

```powershell
# Add to PowerShell profile
skeletor completion powershell | Out-String | Invoke-Expression

# Or save to profile
skeletor completion powershell >> $PROFILE
```

## Troubleshooting

### Common Issues

**Command not found:**
- Ensure the binary is in your PATH
- Check installation directory permissions
- Verify the binary is executable (`chmod +x skeletor`)

**Permission denied:**
- Use `sudo` for system-wide installation
- Install to user directory instead: `~/bin/skeletor`
- Check directory permissions

**Template errors:**
- Ensure Git is installed (required for template cloning)
- Check network connectivity for remote templates
- Verify template repository access

**Go installation issues:**
- Ensure Go 1.21 or later is installed
- Check GOPATH and GOBIN environment variables
- Clear Go module cache: `go clean -modcache`

### Getting Help

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/getporter/skeletor/issues)
2. Review the [documentation](https://getporter.github.io/skeletor/)
3. Join the [Porter community](https://porter.sh/community/)

## Uninstallation

### Homebrew

```bash
brew uninstall skeletor
brew untap getporter/tap  # Optional
```

### Manual Installation

```bash
# Remove binary
sudo rm /usr/local/bin/skeletor

# Remove shell completions
sudo rm /etc/bash_completion.d/skeletor  # Bash
rm ~/.oh-my-zsh/completions/_skeletor    # Zsh
rm ~/.config/fish/completions/skeletor.fish  # Fish
```

### Go Install

```bash
# Remove from GOPATH/bin
rm $(go env GOPATH)/bin/skeletor
```

## Next Steps

After installation:

1. [Quick Start Guide](quick-start.md) - Get started with basic usage
2. [Command Reference](command-reference.md) - Complete CLI documentation
3. [Enterprise Features](enterprise-features.md) - Learn about enterprise capabilities
4. [Examples](examples.md) - See real-world usage examples
