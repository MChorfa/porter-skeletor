//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"get.porter.sh/magefiles/mixins"
	"get.porter.sh/magefiles/releases"
	"github.com/magefile/mage/mg"
)

const (
	mixinName    = "skeletor"
	mixinPackage = "github.com/getporter/skeletor"
	mixinBin     = "bin/mixins/" + mixinName
)

var magefile = mixins.NewMagefile(mixinPackage, mixinName, mixinBin)

// ConfigureAgent sets up the CI server with mage and GO
func ConfigureAgent() {
	magefile.ConfigureAgent()
}

// Build the mixin
func Build() {
	magefile.Build()
}

// XBuildAll cross-compiles the mixin before a release
func XBuildAll() {
	magefile.XBuildAll()
}

// TestUnit runs the unit tests
func TestUnit() {
	magefile.TestUnit()
}

// TestIntegration runs integration tests (requires -tags integration)
func TestIntegration() {
	fmt.Println("Running Integration Tests...")
	testCmd := exec.Command("go", "test", "-v", "-tags=integration", "-timeout", "30m", "./pkg/...") // Run tests in pkg dir
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	if err := testCmd.Run(); err != nil {
		fmt.Printf("Integration tests failed: %v\n", err)
		os.Exit(1) // Fail build on integration test failure
	}
}

// Test runs all types of tests (unit by default from library, plus our integration)
func Test() {
	// magefile.Test() // This likely runs unit tests, let's be explicit
	mg.SerialDeps(TestUnit, TestIntegration) // Run unit then integration tests

	// Run linters and vulnerability checks
	mg.SerialDeps(Lint)
}

// Lint runs linters and vulnerability checks
func Lint() {
	fmt.Println("Running Linters and Security Checks...")

	// Run gosec for security analysis
	fmt.Println("Running gosec...")
	gosecCmd := exec.Command("gosec", "./...")
	gosecCmd.Stdout = os.Stdout
	gosecCmd.Stderr = os.Stderr
	if err := gosecCmd.Run(); err != nil {
		fmt.Printf("gosec failed: %v\n", err)
		// Decide if this should be a hard failure (os.Exit) or just a warning
	}

	// Run govulncheck for vulnerability scanning
	fmt.Println("Running govulncheck...")
	vulnCmd := exec.Command("govulncheck", "./...")
	vulnCmd.Stdout = os.Stdout
	vulnCmd.Stderr = os.Stderr
	if err := vulnCmd.Run(); err != nil {
		fmt.Printf("govulncheck failed: %v\n", err)
		// Decide if this should be a hard failure (os.Exit) or just a warning
	}

	// Add golangci-lint if a config exists (optional, depends if generator itself needs linting)
	// fmt.Println("Running golangci-lint...")
	// lintCmd := exec.Command("golangci-lint", "run", "./...")
	// lintCmd.Stdout = os.Stdout
	// lintCmd.Stderr = os.Stderr
	// if err := lintCmd.Run(); err != nil {
	// 	fmt.Printf("golangci-lint failed: %v\n", err)
	// }
}

// Publish the mixin to GitHub
func Publish() {
	// You can test out publishing locally by overriding PORTER_RELEASE_REPOSITORY and PORTER_PACKAGES_REMOTE
	if _, overridden := os.LookupEnv(releases.ReleaseRepository); !overridden {
		os.Setenv(releases.ReleaseRepository, "github.com/YOURNAME/YOURREPO")
	}
	magefile.PublishBinaries()

	// Publish mixin feed if PORTER_PACKAGES_REMOTE is set (can be set via template variable MixinFeedRepoURL)
	feedRepo := os.Getenv(releases.PackagesRemote)
	if feedRepo != "" {
		fmt.Printf("Publishing mixin feed to %s...\n", feedRepo)
		// Set branch if provided (defaults handled by magefiles library if not set)
		// feedBranch := os.Getenv(releases.PackagesRemote) // TODO: Check updated magefiles lib for branch handling
		// if feedBranch != "" {
		//	 os.Setenv(releases.PackagesBranch, feedBranch) // TODO: Check updated magefiles lib for branch handling
		//	 fmt.Printf("Using branch: %s\n", feedBranch)
		// }
		magefile.PublishMixinFeed()
	} else {
		fmt.Println("Skipping mixin feed publish: PORTER_PACKAGES_REMOTE environment variable not set.")
		fmt.Println("Set the MixinFeedRepoURL variable during generation or set PORTER_PACKAGES_REMOTE manually to enable.")
	}
}

// TestPublish publishes the project to the specified GitHub username.
// If your mixin is official hosted in a repository under your username, you will need to manually
// override PORTER_RELEASE_REPOSITORY and PORTER_PACKAGES_REMOTE to test out publishing safely.
func TestPublish(username string) {
	magefile.TestPublish(username)
}

// Install the mixin
func Install() {
	magefile.Install()
}

// Clean removes generated build files
func Clean() {
	magefile.Clean()
}

// ValidateGenerated runs the Dagger pipeline to validate generated mixin code
func ValidateGenerated() error {
	fmt.Println("Running Dagger pipeline to validate generated mixin...")
	// Assumes the Dagger CLI entrypoint is 'go run ./ci'
	cmd := exec.Command("go", "run", "./ci", "-task", "validate") // Add a 'validate' task to ci/main.go
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
