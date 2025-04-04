# Porter Mixin Generator

[![Build Status](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml/badge.svg)](https://github.com/getporter/skeletor/actions/workflows/skeletor.yml)
[![GitHub Pages](https://github.com/getporter/skeletor/actions/workflows/pages.yml/badge.svg)](https://getporter.github.io/skeletor/) <!-- Add Pages badge/link -->

This repository contains the `skeletor`, a command-line tool designed to streamline the creation of new Porter mixins. It scaffolds a new mixin project based on an enhanced template (sourced from `templates/` in this repository or an external URL/directory), providing a solid foundation with enterprise-grade features built-in.

**[View Documentation Site](https://getporter.github.io/skeletor/)** <!-- Add prominent link -->

## Features

* **Rapid Scaffolding:** Quickly generate a new Porter mixin project structure.
* **Enterprise-Ready Template:** The generated mixin includes:
  * **Observability:** Built-in OpenTelemetry tracing and structured logging, configurable via environment variables.
  * **CI/CD Pipeline:** A default GitHub Actions workflow (`.github/workflows/mixin-ci.yml`) for building, testing, linting, and vulnerability scanning.
  * **Dockerfile:** A multi-stage Dockerfile for creating optimized and secure mixin images.
  * **Security & Contribution Docs:** Standard `SECURITY.md` and `CONTRIBUTING.md` files.
  * **Linting:** Default `golangci-lint` configuration (`.golangci.yml`).
* **Customizable:** Supports flags for non-interactive generation and specifying template variables.
* **Validation:** Performs basic post-generation checks (`go mod tidy`, `go build`, `go test`).
* **Filename Templating:** Supports Go template syntax in template filenames and directory names (e.g., `cmd/{{ .MixinName }}/main.go.tmpl`).
* **Onboarding Guides:** Generates `docs/DEVELOPER_GUIDE.md` and `docs/OPERATIONS_GUIDE.md` within the new mixin project.

## Installation

You can install the generator using `go install`:

```bash
go install github.com/getporter/skeletor/cmd/skeletor@latest
```

Alternatively, build from source:

```bash
git clone https://github.com/getporter/skeletor.git
cd skeletor
go run mage.go build install
# The binary will be in ./bin/skeletor
```

**Using Docker:**

Pull the latest image from GitHub Container Registry:

```bash
docker pull ghcr.io/getporter/skeletor:latest
```

Run the generator using Docker:

```bash
# Example: Create mixin in the current directory, mounting it as /work
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/getporter/skeletor:latest \
  create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin" --output ./my-mixin
```

## Usage

To create a new mixin named `my-mixin`:

```bash
skeletor create --name my-mixin --author "Your Name" --module "github.com/your-org/my-mixin"
```

This will create a new directory `./my-mixin` containing the scaffolded project.

**Flags:**

* `--name`: (Required) Name of the new mixin (lowercase).
* `--author`: (Required) Author name for the mixin.
* `--module`: Go module path (default: `github.com/getporter/<name>`).
* `--output`: Output directory (default: `./<name>`).
* `--non-interactive`: Run without prompts, using defaults or provided flags.
* `--template-url`: URL to a git repository containing a custom template (overrides default).
* `--template-dir`: Local directory containing the template (e.g., `--template-dir templates` to use the one in this repo).
* `--var`: Set template variables in `KEY=VALUE` format (can be used multiple times).
* `--compliance-level`: Desired compliance level ("basic", "slsa-l1", "slsa-l3"; default: "basic"). Affects generated files like Dockerfile, .goreleaser.yml, .golangci.yml, SECURITY.md.
* `--dry-run`: Simulate generation without writing files or running hooks.

## Template Variables

The following variables are used by the default template (`templates/template.json`) and can be provided during generation (interactively or via `--var` flag):

* `MixinName` (string, required): Name of the mixin (lowercase).
* `AuthorName` (string, required): Author name.
* `ModulePath` (string): Go module path (defaults based on `MixinName`).
* `Description` (string): Short description of the mixin (defaults based on `MixinName`).
* `License` (string): License for the mixin (choices: "Apache-2.0", "MIT", "GPL-3.0"; default: "Apache-2.0").
* `InitGit` (bool): Initialize a git repository in the output directory? (default: true).
* `MixinFeedRepoURL` (string, optional): Git URL for the mixin feed repository. If provided, `mage Publish` will attempt to update the feed.
* `MixinFeedBranch` (string): Branch in the mixin feed repository (default: "main").
* `AuthorEmail` (string, optional): Author's email for security contact (used in `.well-known/security.txt`).

# Note: ComplianceLevel is now a direct flag (--compliance-level), not a template variable

## Generated Project Structure

The generated mixin project follows the standard Porter mixin structure:

* `cmd/YOURMIXIN/`: CLI implementation using Cobra. Sourced from `templates/cmd/mixin/*.go.tmpl`.
* `pkg/YOURMIXIN/`: Core mixin logic (implement build, install, invoke, etc.). Sourced from `templates/pkg/mixin/`. *(Note: `pkg/` directory structure might need creation/templating if not already present)*
* `ci/main.go`: Dagger pipeline definition for CI/CD tasks (test, build, release). Sourced from `templates/ci/main.go.tmpl`.
* `.github/workflows/mixin-ci.yml`: GitHub Actions workflow that executes the Dagger pipeline. Sourced from `templates/.github/workflows/mixin-ci.yml.tmpl`.
* `magefile.go`: Build automation using Mage (can be invoked by Dagger). Sourced from `magefile.go` at the root.
* `.goreleaser.yml`: Configuration for GoReleaser (used by Dagger release task).
* `tools.go`: Go tool dependencies.
* `go.mod`, `go.sum`: Go module files.
* `Dockerfile`: For building the mixin container image (can be used by GoReleaser). Sourced from `templates/Dockerfile.tmpl`.
* `.golangci.yml`: Default or strict linter configuration (conditionally generated). Sourced from `templates/.golangci.yml.tmpl` or `templates/.golangci-strict.yml.tmpl`.
* `.well-known/security.txt`: Standard security contact file. Sourced from `templates/.well-known/security.txt.tmpl`.
* `docs/DEVELOPER_GUIDE.md`: Guide for developers extending the mixin. Sourced from `templates/docs/DEVELOPER_GUIDE.md.tmpl`.
* `docs/OPERATIONS_GUIDE.md`: Guide for users operating the mixin. Sourced from `templates/docs/OPERATIONS_GUIDE.md.tmpl`.
* `README.md`, `LICENSE`, `CONTRIBUTING.md`, `SECURITY.md`: Documentation and policy files. Sourced from templates.

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
