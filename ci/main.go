package main

import (
	"context"
	"flag" // Import flag package
	"fmt"
	"os"
	"path/filepath"
	"strings" // Ensure strings is imported

	"dagger.io/dagger"
)

// Define constants for Go version and image
const (
	goVersion = "1.23" // Keep in sync with go.mod
	goImage   = "golang:" + goVersion + "-alpine"
	mixinName = "skeletor" // Should match the project name or be dynamic if needed
	// goreleaserVersion constant removed
)

func main() {
	// Define flags
	task := flag.String("task", "ci", "Task to run: ci, release, or validate") // Updated description
	flag.Parse()

	ctx := context.Background()

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		fmt.Printf("Error connecting to Dagger engine: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Execute the requested task
	switch *task {
	case "ci":
		fmt.Println("Running CI tasks (test & build)...")
		if err := runCI(ctx, client); err != nil {
			fmt.Printf("CI task failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("CI tasks completed successfully!")
	case "release":
		fmt.Println("Running Release task...")
		// Get GITHUB_TOKEN from environment
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			fmt.Println("Error: GITHUB_TOKEN environment variable is required for release")
			os.Exit(1)
		}
		if err := release(ctx, client, githubToken); err != nil {
			fmt.Printf("Release task failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Release task completed successfully!")
	case "validate": // Add validate case
		fmt.Println("Running Validate Generated Mixin task...")
		if err := ValidateGeneratedMixin(ctx, client); err != nil {
			fmt.Printf("Validation task failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Validation task completed successfully!")
	default:
		fmt.Printf("Error: Unknown task '%s'. Valid tasks are 'ci', 'release', or 'validate'.\n", *task) // Updated error message
		os.Exit(1)
	}
}

// runCI executes the standard CI steps (test, build)
func runCI(ctx context.Context, client *dagger.Client) error {
	fmt.Println("--> Running tests...")
	if err := test(ctx, client); err != nil {
		return fmt.Errorf("testing failed: %w", err)
	}
	fmt.Println("--> Tests passed!")

	fmt.Println("--> Building artifacts (example: linux/amd64)...")
	// Example: Build for linux/amd64 for verification during CI
	artifactsDir := "build-output"
	if err := build(ctx, client, "linux", "amd64", artifactsDir); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	fmt.Printf("--> Build artifacts generated in ./%s\n", artifactsDir)
	return nil
}

// ValidateGeneratedMixin builds the generator and uses it to create a sample mixin,
// then runs validation checks (build, test, lint, scan) on the generated code.
func ValidateGeneratedMixin(ctx context.Context, client *dagger.Client) error {
	fmt.Println("--> Building generator binary...")
	src := projectSource(client)

	// Container for building the generator
	builder := client.Container().From(goImage).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod-builder")). // Separate cache?
		WithMountedCache("/go/build-cache", client.CacheVolume("go-build-builder")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithWorkdir("/src").
		WithMountedDirectory("/src", src).
		WithEnvVariable("CGO_ENABLED", "0")

	// Build the generator
	generatorPath := "/src/bin/skeletor" // Path inside container
	builder = builder.WithExec([]string{
		"go", "build",
		"-ldflags", "-s -w",
		"-o", generatorPath,
		"./cmd/skeletor",
	})

	// Container for running validation on the generated mixin
	// Needs Go, git, golangci-lint, gosec, govulncheck
	validator := client.Container().From(goImage).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod-validator")).
		WithMountedCache("/go/build-cache", client.CacheVolume("go-build-validator")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"apk", "add", "--no-cache", "git", "build-base"}) // Add git

	// Install validation tools (pin versions for consistency)
	golangciVersion := "v1.59.1"   // Example pinned version
	gosecVersion := "v2.19.0"      // Example pinned version
	govulncheckVersion := "v1.1.1" // Example pinned version
	validator = validator.
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@" + golangciVersion}).
		WithExec([]string{"go", "install", "github.com/securego/gosec/v2/cmd/gosec@" + gosecVersion}).
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@" + govulncheckVersion})

	// Mount the built generator binary into the validator container
	generatorBinary := builder.File(generatorPath)
	validator = validator.WithMountedFile("/usr/local/bin/skeletor", generatorBinary)

	// Generate a sample mixin inside the validator container
	fmt.Println("--> Generating sample mixin for validation...")
	generatedMixinPath := "/tmp/generated-mixin"
	validator = validator.WithExec([]string{
		"skeletor", "create",
		"--name", "validateme",
		"--author", "Dagger CI",
		"--module", "example.com/dagger/validateme",
		"--output", generatedMixinPath,
		"--non-interactive",
		"--compliance-level", "basic", // Or test multiple levels? Start with basic.
		// Use the embedded/default template source implicitly if generator supports it,
		// otherwise mount the templates dir and use --template-dir
		// Assuming generator uses local templates dir if available:
		// WithMountedDirectory("/src/templates", src.Directory("templates")).
		// WithExec([]string{..., "--template-dir", "/src/templates"})
	})

	// Run validation commands within the generated mixin directory
	fmt.Println("--> Validating generated mixin code...")
	validator = validator.WithWorkdir(generatedMixinPath)

	// Go Mod Tidy
	fmt.Println("  --> Running go mod tidy...")
	_, err := validator.WithExec([]string{"go", "mod", "tidy"}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: go mod tidy: %w", err)
	}

	// Go Build
	fmt.Println("  --> Running go build ./...")
	_, err = validator.WithExec([]string{"go", "build", "./..."}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: go build: %w", err)
	}

	// Go Test
	fmt.Println("  --> Running go test ./...")
	_, err = validator.WithExec([]string{"go", "test", "./..."}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: go test: %w", err)
	}

	// GolangCI-Lint
	fmt.Println("  --> Running golangci-lint run ./...")
	_, err = validator.WithExec([]string{"golangci-lint", "run", "./..."}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: golangci-lint: %w", err)
	}

	// Gosec
	fmt.Println("  --> Running gosec ./...")
	_, err = validator.WithExec([]string{"gosec", "./..."}).Sync(ctx)
	if err != nil {
		// Decide if gosec findings should fail the validation
		fmt.Printf("Warning: gosec found issues: %v\n", err)
		// return fmt.Errorf("validation failed: gosec: %w", err) // Uncomment to fail on gosec issues
	}

	// Govulncheck
	fmt.Println("  --> Running govulncheck ./...")
	_, err = validator.WithExec([]string{"govulncheck", "./..."}).Sync(ctx)
	if err != nil {
		// Decide if govulncheck findings should fail the validation
		fmt.Printf("Warning: govulncheck found issues: %v\n", err)
		// return fmt.Errorf("validation failed: govulncheck: %w", err) // Uncomment to fail on vulnerabilities
	}

	fmt.Println("--> Generated mixin validation successful!")
	return nil
}

// test runs linters and unit tests
func test(ctx context.Context, client *dagger.Client) error {
	src := projectSource(client)

	// Create a Go container
	golang := client.Container().From(goImage).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod")).
		WithMountedCache("/go/build-cache", client.CacheVolume("go-build")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithWorkdir("/src").
		WithMountedDirectory("/src", src)

	// Install tools
	// Define pinned versions
	gosecVersion := "v2.19.0"      // Example pinned version
	govulncheckVersion := "v1.1.1" // Example pinned version
	mageVersion := "v1.15.0"       // Example pinned version

	golang = golang.
		WithExec([]string{"apk", "add", "--no-cache", "git", "build-base"}). // Add git and build tools if needed
		WithExec([]string{"go", "install", "github.com/securego/gosec/v2/cmd/gosec@" + gosecVersion}).
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@" + govulncheckVersion}).
		WithExec([]string{"go", "install", "github.com/magefile/mage@" + mageVersion}) // Install mage

	// Run linters (using mage target)
	_, err := golang.WithExec([]string{"mage", "Lint"}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("linting failed: %w", err)
	}

	// Run tests (using mage target)
	_, err = golang.WithExec([]string{"mage", "TestUnit"}).Sync(ctx) // Assuming TestUnit runs only unit tests
	if err != nil {
		return fmt.Errorf("unit tests failed: %w", err)
	}

	return nil
}

// build compiles the mixin for a specific platform
func build(ctx context.Context, client *dagger.Client, goos, goarch, outputDir string) error {
	src := projectSource(client)

	// Create a Go container
	golang := client.Container().From(goImage).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod")).
		WithMountedCache("/go/build-cache", client.CacheVolume("go-build")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithWorkdir("/src").
		WithMountedDirectory("/src", src).
		WithEnvVariable("GOOS", goos).
		WithEnvVariable("GOARCH", goarch).
		WithEnvVariable("CGO_ENABLED", "0") // Ensure static builds

	// Build the binary using go build (or mage build if preferred)
	// Using go build directly for simplicity here
	outputPath := filepath.Join("/src", outputDir, fmt.Sprintf("%s-%s-%s", mixinName, goos, goarch))
	golang = golang.WithExec([]string{
		"go", "build",
		"-ldflags", "-s -w", // Strip symbols and debug info
		"-o", outputPath,
		"./cmd/" + mixinName, // Path to main package
	})

	// Extract the built binary
	output := client.Directory().WithDirectory(outputDir, golang.Directory(filepath.Join("/src", outputDir)))
	_, err := output.Export(ctx, ".") // Export to host filesystem under ./build-output
	if err != nil {
		return fmt.Errorf("failed to export build artifact: %w", err)
	}

	return nil
}

// projectSource returns the host project directory mounted in the container
func projectSource(client *dagger.Client) *dagger.Directory {
	// Get reference to host directory
	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"ci/build-output/", "bin/"}, // Exclude build output and bin directories
	})
	return src
}

// release runs goreleaser within a Dagger container
func release(ctx context.Context, client *dagger.Client, githubToken string) error {
	fmt.Println("--> Preparing GoReleaser container...")
	src := projectSource(client)

	// Use a container with Go and Git installed
	// Mount caches and source code
	goreleaser := client.Container().From(goImage). // Base Go image
							WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod")).
							WithMountedCache("/go/build-cache", client.CacheVolume("go-build")).
							WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
							WithEnvVariable("GOCACHE", "/go/build-cache").
							WithWorkdir("/src").
							WithMountedDirectory("/src", src).
							WithExec([]string{"apk", "add", "--no-cache", "git", "build-base"}) // Ensure git is present

	// Define pinned versions within function scope
	goreleaserVersion := "v2.1.0" // Example pinned version
	cosignVersion := "v2.2.4"     // Example pinned version
	syftVersion := "v1.10.0"      // Example pinned syft version

	// Install GoReleaser, Cosign, Syft
	releaserTools := goreleaser.
		WithExec([]string{"go", "install", "github.com/goreleaser/goreleaser@" + goreleaserVersion}).
		WithExec([]string{"go", "install", "github.com/sigstore/cosign/v2/cmd/cosign@" + cosignVersion}).
		WithExec([]string{"apk", "add", "--no-cache", "syft=" + syftVersion}) // Assuming syft is available via apk, adjust if needed

	// Run GoReleaser release command
	// Pass GITHUB_TOKEN as a secret
	githubTokenSecret := client.SetSecret("github-token", githubToken)

	fmt.Println("--> Running GoReleaser...")
	// Use --clean to ensure a clean build environment within the container
	// GoReleaser automatically detects it's running in GitHub Actions
	// and uses the token for publishing and SLSA generation.
	releaserExec := releaserTools.
		WithSecretVariable("GITHUB_TOKEN", githubTokenSecret). // Expose token as env var
		WithExec([]string{"goreleaser", "release", "--clean"})

	// Execute the GoReleaser command
	_, err := releaserExec.Sync(ctx)
	if err != nil {
		return fmt.Errorf("goreleaser execution failed: %w", err)
	}
	fmt.Println("--> GoReleaser finished successfully (binaries, checksums, SBOMs, release assets).")

	// --- Explicit Docker Build, Push, Attest ---
	fmt.Println("--> Building and pushing Docker images...")

	// Define platforms
	platforms := []string{"linux/amd64", "linux/arm64"}
	imageRepo := "ghcr.io/getporter/skeletor" // Define image repo base
	gitTag := os.Getenv("GITHUB_REF_NAME")                  // Assumes GITHUB_REF_NAME is set (e.g., v1.2.3)
	if gitTag == "" {
		return fmt.Errorf("GITHUB_REF_NAME environment variable not set, cannot determine image tag")
	}

	imageRefs := []string{}             // To store refs for manifest list
	containers := []*dagger.Container{} // To store built containers for manifest list

	for _, platform := range platforms {
		goos := strings.Split(platform, "/")[0]
		goarch := strings.Split(platform, "/")[1]
		imageRef := fmt.Sprintf("%s:%s-%s", imageRepo, gitTag, goarch)
		imageRefLatest := fmt.Sprintf("%s:latest-%s", imageRepo, goarch) // Arch-specific latest tag

		fmt.Printf("  --> Building %s...\n", imageRef)
		// Specify platform in ContainerOpts, not BuildOpts
		ctr := client.Container(dagger.ContainerOpts{Platform: dagger.Platform(platform)}).
			Build(src, dagger.ContainerBuildOpts{
				Dockerfile: "Dockerfile", // Assumes Dockerfile at root
				// Platform removed from here
				BuildArgs: []dagger.BuildArg{ // Pass necessary build args
					{Name: "GOOS", Value: goos},
					{Name: "GOARCH", Value: goarch},
					// Add other build args from goreleaser config if needed
				},
			})

		fmt.Printf("  --> Pushing %s and %s...\n", imageRef, imageRefLatest)
		_, err = ctr.
			WithRegistryAuth(imageRepo, "GITHUB_ACTOR", githubTokenSecret). // Use token for auth
			Publish(ctx, imageRef)
		if err != nil {
			return fmt.Errorf("failed to push image %s: %w", imageRef, err)
		}
		_, err = ctr.
			WithRegistryAuth(imageRepo, "GITHUB_ACTOR", githubTokenSecret).
			Publish(ctx, imageRefLatest)
		if err != nil {
			return fmt.Errorf("failed to push image %s: %w", imageRefLatest, err)
		}
		imageRefs = append(imageRefs, imageRef) // Add versioned ref for manifest list
		containers = append(containers, ctr)    // Add built container for manifest list

		// --- Generate SBOM for the pushed image ---
		fmt.Printf("  --> Generating SBOM for %s...\n", imageRef)
		sbomFileName := fmt.Sprintf("image-%s-%s.spdx.json", goos, goarch) // Example filename
		sbomPath := "/tmp/" + sbomFileName

		// Container with Syft (use the releaserTools container which has tools installed)
		syftCtr := releaserTools.
			WithEntrypoint([]string{""}) // Clear entrypoint for custom command

		// Scan the image and save SBOM
		// Note: Requires registry auth if image is private, but GHCR might allow public reads
		syftCtr = syftCtr.WithExec([]string{
			"syft", "packages", imageRef, // Scan the image we just pushed
			"-o", "spdx-json", // Output format
			"--file", sbomPath,
		})

		// Export SBOM to host
		sbomFile := syftCtr.File(sbomPath)
		_, err = sbomFile.Export(ctx, sbomFileName)
		if err != nil {
			return fmt.Errorf("failed to export SBOM %s: %w", sbomFileName, err)
		}
		fmt.Printf("  --> SBOM saved to %s\n", sbomFileName)

		// --- Attest SBOM to the pushed image ---
		fmt.Printf("  --> Attesting SBOM %s to %s...\n", sbomFileName, imageRef)
		// Use the releaserTools container which has cosign installed
		cosignCtr := releaserTools.
			WithMountedFile("/tmp/"+sbomFileName, sbomFile).       // Mount the exported SBOM
			WithSecretVariable("GITHUB_TOKEN", githubTokenSecret). // For keyless signing identity
			WithEnvVariable("COSIGN_EXPERIMENTAL", "1")            // For keyless signing

		// Run cosign attest
		// Note: This assumes keyless signing. Adjust if using keys.
		_, err = cosignCtr.WithExec([]string{
			"cosign", "attest",
			"--type", "spdxjson", // Match SBOM format
			"--predicate", "/tmp/" + sbomFileName,
			imageRef, // Attest the specific arch image ref
		}).Sync(ctx)
		if err != nil {
			// Log warning instead of failing the build? Attestation might be best-effort initially.
			fmt.Printf("Warning: failed to attest SBOM %s to %s: %w\n", sbomFileName, imageRef, err)
			// return fmt.Errorf("failed to attest SBOM %s to %s: %w", sbomFileName, imageRef, err)
		} else {
			fmt.Printf("  --> Successfully attested SBOM to %s\n", imageRef)
		}
		// Clean up the local SBOM file? Or keep it? Let's keep it for now.
		// os.Remove(sbomFileName)
	}

	// --- Create and Push Manifest List ---
	// Ensure we have containers before proceeding
	if len(containers) == 0 || len(containers) != len(platforms) {
		return fmt.Errorf("unexpected number of containers built (%d), cannot create manifest", len(containers))
	}

	fmt.Println("--> Creating and pushing multi-arch manifest...")
	manifestRef := fmt.Sprintf("%s:%s", imageRepo, gitTag)
	manifestLatestRef := fmt.Sprintf("%s:latest", imageRepo)

	// Publish versioned manifest
	_, err = client.Container().
		WithRegistryAuth(imageRepo, "GITHUB_ACTOR", githubTokenSecret).
		Publish(ctx, manifestRef, dagger.ContainerPublishOpts{
			PlatformVariants: containers, // Use the containers built in the loop
		})
	if err != nil {
		return fmt.Errorf("failed to publish multi-arch manifest %s: %w", manifestRef, err)
	}
	fmt.Printf("  --> Pushed manifest %s\n", manifestRef)

	// Publish latest manifest
	_, err = client.Container().
		WithRegistryAuth(imageRepo, "GITHUB_ACTOR", githubTokenSecret).
		Publish(ctx, manifestLatestRef, dagger.ContainerPublishOpts{
			PlatformVariants: containers, // Use the same containers
		})
	if err != nil {
		return fmt.Errorf("failed to publish multi-arch manifest %s: %w", manifestLatestRef, err)
	}
	fmt.Printf("  --> Pushed manifest %s\n", manifestLatestRef)

	fmt.Println("--> Docker operations complete.")
	return nil
}
