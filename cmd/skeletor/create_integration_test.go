package main

import (
	"bytes" // Ensure bytes is imported
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function to get the path to the built generator binary
func getGeneratorBinaryPath(t *testing.T) string {
	// Assume the binary is built in the project root for testing purposes
	// Adjust this path if the build process places it elsewhere
	binaryName := "skeletor"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	// Assuming tests run from the cmd/skeletor directory
	projectRoot, err := filepath.Abs("../..") // Go up two levels
	require.NoError(t, err, "Failed to get project root")
	binaryPath := filepath.Join(projectRoot, binaryName)

	// Basic check if the binary exists, maybe build it if not?
	// For now, assume it's pre-built by a 'mage build' or similar
	_, err = os.Stat(binaryPath)
	if os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s. Build the project first (e.g., 'mage build').", binaryPath)
	}
	require.NoError(t, err, "Error checking generator binary")

	return binaryPath
}

// Helper function to run the generator create command
func runGeneratorCreate(t *testing.T, binaryPath string, args ...string) (string, error) {
	outputDir := t.TempDir()
	baseArgs := []string{"create", "--output", outputDir, "--non-interactive"} // Non-interactive for tests
	fullArgs := append(baseArgs, args...)

	cmd := exec.Command(binaryPath, fullArgs...)
	cmd.Stdout = os.Stdout // Or capture if needed
	cmd.Stderr = os.Stderr // Or capture if needed

	err := cmd.Run()
	return outputDir, err
}

func TestCreateMixin_Integration_BasicCompliance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	binaryPath := getGeneratorBinaryPath(t)

	// Define common args for a basic mixin
	mixinName := "basic-test-mixin"
	authorName := "Test Author"
	modulePath := "example.com/getporter/" + mixinName
	args := []string{
		"--name", mixinName,
		"--author", authorName,
		"--module", modulePath,
		"--compliance-level", "basic", // Explicitly test basic level
		// Add other necessary flags/vars if defaults aren't sufficient
	}

	// Run the generator
	outputDir, err := runGeneratorCreate(t, binaryPath, args...)
	require.NoError(t, err, "Generator create command failed for basic compliance")
	defer os.RemoveAll(outputDir) // Clean up generated files

	// --- Assertions for Basic Compliance ---

	// 1. Check for essential files existence
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"porter.yaml",
		"Dockerfile",
		".goreleaser.yml",
		"cmd/" + mixinName + "/main.go",
		"pkg/" + mixinName + "/mixin.go", // Assuming pkg structure exists
		".golangci.yml",                  // Expecting the non-strict version
		// Add other core files expected for all levels
	}
	for _, file := range expectedFiles {
		_, err := os.Stat(filepath.Join(outputDir, file))
		require.NoError(t, err, "Expected file %s not found in output for basic compliance", file)
	}

	// 2. Check that SLSA L3 specific files DO NOT exist
	notExpectedFiles := []string{
		".golangci-strict.yml", // Should not exist for basic
	}
	for _, file := range notExpectedFiles {
		_, err := os.Stat(filepath.Join(outputDir, file))
		require.Error(t, err, "File %s should NOT exist for basic compliance", file)
		require.True(t, os.IsNotExist(err), "Error for %s should be os.IsNotExist", file)
	}

	// 3. TODO: Add content checks for specific files (e.g., Dockerfile should match basic template section)
	//    - Read Dockerfile content
	//    - Assert it contains markers for basic level and not L1/L3
	//    - Read .goreleaser.yml content
	//    - Assert it does NOT contain 'slsa:' or 'signs:' blocks
	dockerfilePath := filepath.Join(outputDir, "Dockerfile")
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	require.NoError(t, err, "Failed to read generated Dockerfile")
	require.Contains(t, string(dockerfileContent), "# --- Basic Compliance Level ---", "Dockerfile should contain basic compliance marker")
	require.NotContains(t, string(dockerfileContent), "# --- SLSA Level 1 Compliance ---", "Dockerfile should NOT contain L1 marker for basic")
	require.NotContains(t, string(dockerfileContent), "# --- SLSA Level 3 Compliance ---", "Dockerfile should NOT contain L3 marker for basic")

	goreleaserPath := filepath.Join(outputDir, ".goreleaser.yml")
	goreleaserContent, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err, "Failed to read generated .goreleaser.yml")
	require.NotContains(t, string(goreleaserContent), "slsa:", ".goreleaser.yml should not contain slsa block for basic")
	require.NotContains(t, string(goreleaserContent), "signs:", ".goreleaser.yml should not contain signs block for basic")

	// Check .golangci.yml content (should be basic)
	golangciPath := filepath.Join(outputDir, ".golangci.yml")
	golangciContent, err := os.ReadFile(golangciPath)
	require.NoError(t, err, "Failed to read generated .golangci.yml for basic")
	require.Contains(t, string(golangciContent), "basic lint config for "+mixinName, "Expected basic lint config content") // Check marker from mock file
	require.NotContains(t, string(golangciContent), "strict lint config for "+mixinName, "Should not contain strict lint config content")

	// 4. Check if post-gen hooks ran (e.g., check for go.sum)
	_, err = os.Stat(filepath.Join(outputDir, "go.sum"))
	require.NoError(t, err, "go.sum not found, post-gen hook 'go mod tidy' might not have run")

	// 5. Lint the generated code
	t.Logf("Linting generated code in %s...", outputDir)
	lintCmd := exec.Command("golangci-lint", "run", "./...")
	lintCmd.Dir = outputDir
	lintOutput, lintErr := lintCmd.CombinedOutput() // Capture output for debugging
	require.NoError(t, lintErr, "Linting generated code failed. Output:\n%s", string(lintOutput))
	t.Logf("Linting successful.")

	t.Logf("Successfully generated and validated mixin in %s for basic compliance test", outputDir)
}

func TestCreateMixin_Integration_SlsaL1Compliance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	binaryPath := getGeneratorBinaryPath(t)

	// Define common args for an L1 mixin
	mixinName := "l1-test-mixin"
	authorName := "Test Author L1"
	modulePath := "example.com/getporter/" + mixinName
	args := []string{
		"--name", mixinName,
		"--author", authorName,
		"--module", modulePath,
		"--compliance-level", "slsa-l1", // Explicitly test L1 level
	}

	// Run the generator
	outputDir, err := runGeneratorCreate(t, binaryPath, args...)
	require.NoError(t, err, "Generator create command failed for slsa-l1 compliance")
	defer os.RemoveAll(outputDir)

	// --- Assertions for SLSA L1 Compliance ---

	// 1. Check essential files (similar to basic, assuming L1 doesn't remove core files)
	expectedFiles := []string{
		"go.mod", "README.md", "porter.yaml", "Dockerfile", ".goreleaser.yml",
		"cmd/" + mixinName + "/main.go", "pkg/" + mixinName + "/mixin.go",
		".golangci.yml", // Expecting non-strict for L1
	}
	for _, file := range expectedFiles {
		_, err := os.Stat(filepath.Join(outputDir, file))
		require.NoError(t, err, "Expected file %s not found in output for slsa-l1 compliance", file)
	}

	// 2. Check that SLSA L3 specific files DO NOT exist
	notExpectedFiles := []string{".golangci-strict.yml"}
	for _, file := range notExpectedFiles {
		_, err := os.Stat(filepath.Join(outputDir, file))
		require.Error(t, err, "File %s should NOT exist for slsa-l1 compliance", file)
		require.True(t, os.IsNotExist(err), "Error for %s should be os.IsNotExist", file)
	}

	// 3. Check Dockerfile content for L1 marker
	dockerfilePath := filepath.Join(outputDir, "Dockerfile")
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	require.NoError(t, err, "Failed to read generated Dockerfile for L1")
	require.Contains(t, string(dockerfileContent), "# --- SLSA Level 1 Compliance ---", "Dockerfile should contain L1 compliance section marker")
	require.Contains(t, string(dockerfileContent), "# Placeholder: Add SLSA L1 specific steps here.", "Dockerfile should contain L1 placeholder comment") // More specific check
	require.NotContains(t, string(dockerfileContent), "# --- Basic Compliance Level ---", "Dockerfile should NOT contain basic marker for L1")
	require.NotContains(t, string(dockerfileContent), "# --- SLSA Level 3 Compliance ---", "Dockerfile should NOT contain L3 marker for L1")

	// 4. Check .goreleaser.yml content (expect no slsa or signs blocks for L1)
	goreleaserPath := filepath.Join(outputDir, ".goreleaser.yml")
	goreleaserContent, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err, "Failed to read generated .goreleaser.yml for L1")
	require.NotContains(t, string(goreleaserContent), "slsa:", ".goreleaser.yml should not contain slsa block for L1")
	require.NotContains(t, string(goreleaserContent), "signs:", ".goreleaser.yml should not contain signs block for L1")

	// 5. Check hooks ran
	_, err = os.Stat(filepath.Join(outputDir, "go.sum"))
	require.NoError(t, err, "go.sum not found, post-gen hook 'go mod tidy' might not have run for L1")

	// 6. Lint the generated code
	t.Logf("Linting generated code in %s for L1...", outputDir)
	lintCmd := exec.Command("golangci-lint", "run", "./...")
	lintCmd.Dir = outputDir
	lintOutput, lintErr := lintCmd.CombinedOutput()
	require.NoError(t, lintErr, "Linting generated code failed for L1. Output:\n%s", string(lintOutput))
	t.Logf("Linting successful for L1.")

	t.Logf("Successfully generated and validated mixin in %s for slsa-l1 compliance test", outputDir)
}

func TestCreateMixin_Integration_SlsaL3Compliance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	binaryPath := getGeneratorBinaryPath(t)

	// Define common args for an L3 mixin
	mixinName := "l3-test-mixin"
	authorName := "Test Author L3"
	modulePath := "example.com/getporter/" + mixinName
	authorEmail := "l3-security@example.com" // For SECURITY.md check
	args := []string{
		"--name", mixinName,
		"--author", authorName,
		"--module", modulePath,
		"--compliance-level", "slsa-l3", // Explicitly test L3 level
		"--var", "AuthorEmail=" + authorEmail,
	}

	// Run the generator
	outputDir, err := runGeneratorCreate(t, binaryPath, args...)
	require.NoError(t, err, "Generator create command failed for slsa-l3 compliance")
	defer os.RemoveAll(outputDir)

	// --- Assertions for SLSA L3 Compliance ---

	// 1. Check essential files (including L3 specific ones)
	expectedFiles := []string{
		"go.mod", "README.md", "porter.yaml", "Dockerfile", ".goreleaser.yml",
		"cmd/" + mixinName + "/main.go", "pkg/" + mixinName + "/mixin.go",
		".golangci.yml", // This is the DESTINATION file name
		"SECURITY.md",
	}
	for _, file := range expectedFiles {
		_, err := os.Stat(filepath.Join(outputDir, file))
		require.NoError(t, err, "Expected file %s not found in output for slsa-l3 compliance", file)
	}

	// 2. Check that non-L3 specific files DO NOT exist (if any)
	//    (In this case, the non-strict .golangci.yml source template shouldn't be copied directly)
	//    We check the *content* of the destination .golangci.yml below.

	// 3. Check Dockerfile content for L3 marker
	dockerfilePath := filepath.Join(outputDir, "Dockerfile")
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	require.NoError(t, err, "Failed to read generated Dockerfile for L3")
	require.Contains(t, string(dockerfileContent), "# --- SLSA Level 3 Compliance ---", "Dockerfile should contain L3 compliance marker")
	require.NotContains(t, string(dockerfileContent), "# --- Basic Compliance Level ---", "Dockerfile should NOT contain basic marker for L3")
	require.NotContains(t, string(dockerfileContent), "# --- SLSA Level 1 Compliance ---", "Dockerfile should NOT contain L1 marker for L3")

	// 4. Check .goreleaser.yml content (expect slsa and signs blocks for L3)
	goreleaserPath := filepath.Join(outputDir, ".goreleaser.yml")
	goreleaserContent, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err, "Failed to read generated .goreleaser.yml for L3")
	require.Contains(t, string(goreleaserContent), "slsa:", ".goreleaser.yml should contain slsa block for L3")
	require.Contains(t, string(goreleaserContent), "signs:", ".goreleaser.yml should contain signs block for L3")

	// 5. Check .golangci.yml content (should be from the strict template)
	golangciPath := filepath.Join(outputDir, ".golangci.yml")
	golangciContent, err := os.ReadFile(golangciPath)
	require.NoError(t, err, "Failed to read generated .golangci.yml for L3")
	require.Contains(t, string(golangciContent), "strict lint config for "+mixinName, "Expected strict lint config content") // Check content marker
	require.NotContains(t, string(golangciContent), "basic lint config for "+mixinName, "Should not contain basic lint config content")

	// 6. Check SECURITY.md content for L3 sections and email
	securityPath := filepath.Join(outputDir, "SECURITY.md")
	securityContent, err := os.ReadFile(securityPath)
	require.NoError(t, err, "Failed to read generated SECURITY.md for L3")
	require.Contains(t, string(securityContent), "## Build Integrity & Provenance (SLSA Level 3)", "SECURITY.md should contain L3 provenance section")
	require.Contains(t, string(securityContent), "## Binary Signing (SLSA Level 3)", "SECURITY.md should contain L3 signing section")
	require.Contains(t, string(securityContent), "at "+authorEmail, "SECURITY.md should contain author email")

	// 7. Check hooks ran
	_, err = os.Stat(filepath.Join(outputDir, "go.sum"))
	require.NoError(t, err, "go.sum not found, post-gen hook 'go mod tidy' might not have run for L3")

	// 8. Lint the generated code (using the generated strict config)
	t.Logf("Linting generated code in %s for L3...", outputDir)
	lintCmd := exec.Command("golangci-lint", "run", "./...") // Assumes golangci-lint respects the .golangci.yml in the dir
	lintCmd.Dir = outputDir
	lintOutput, lintErr := lintCmd.CombinedOutput()
	require.NoError(t, lintErr, "Linting generated code failed for L3. Output:\n%s", string(lintOutput))
	t.Logf("Linting successful for L3.")

	t.Logf("Successfully generated and validated mixin in %s for slsa-l3 compliance test", outputDir)
}

func TestCreateMixin_Integration_DryRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	binaryPath := getGeneratorBinaryPath(t)

	// Define common args
	mixinName := "dryrun-test-mixin"
	authorName := "Test Author DryRun"
	modulePath := "example.com/getporter/" + mixinName
	args := []string{
		"--name", mixinName,
		"--author", authorName,
		"--module", modulePath,
		"--compliance-level", "basic", // Use any level for dry run check
		"--dry-run", // Enable dry run
	}

	// Run the generator and capture output
	outputDir := t.TempDir() // Need a path, though it shouldn't be used
	defer os.RemoveAll(outputDir)

	baseArgs := []string{"create", "--output", outputDir, "--non-interactive"}
	fullArgs := append(baseArgs, args...)
	cmd := exec.Command(binaryPath, fullArgs...)

	// Capture combined stdout and stderr
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb

	err := cmd.Run()
	output := outb.String()
	t.Logf("Dry run output:\n%s", output) // Log output for debugging

	// Dry run should exit successfully
	require.NoError(t, err, "Generator create command with --dry-run failed")

	// Assert that dry run simulation messages are present
	require.Contains(t, output, "[Dry Run] Simulating file generation...", "Expected dry run start message")
	require.Contains(t, output, "[Dry Run] Would create directory:", "Expected dry run message for directory creation")
	require.Contains(t, output, "[Dry Run] Would write file:", "Expected dry run message for file writing")
	require.Contains(t, output, "[Dry Run] Skipping post-generation validation.", "Expected dry run message for skipping validation")
	require.Contains(t, output, "[Dry Run] Skipping post-generation hooks.", "Expected dry run message for skipping hooks")
	require.Contains(t, output, "[Dry Run] Simulation complete.", "Expected dry run completion message")

	// Assert that no files were actually created in the temp dir
	files, readErr := os.ReadDir(outputDir)
	require.NoError(t, readErr, "Failed to read temporary output directory")
	require.Empty(t, files, "No files should be created in the output directory during dry run")

	t.Logf("Successfully verified dry run simulation for mixin %s", mixinName)
}

func TestCreateMixin_Integration_TemplateDir(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	binaryPath := getGeneratorBinaryPath(t)

	// Create a temporary local template directory structure
	localTemplateDir := t.TempDir()
	templateJsonContent := `{
		"name": "Local Template Test",
		"variables": {
			"MixinName": {"type": "string", "required": true},
			"LocalVar": {"type": "string", "default": "localDefault"}
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(localTemplateDir, "template.json"), []byte(templateJsonContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(localTemplateDir, "test.txt.tmpl"), []byte("Local template var: {{ .LocalVar }}"), 0644))

	// Define args
	mixinName := "local-template-mixin"
	args := []string{
		"--name", mixinName,
		"--author", "Local Tester",
		"--module", "example.com/local/" + mixinName,
		"--template-dir", localTemplateDir, // Use the local template dir flag
	}

	// Run the generator
	outputDir, err := runGeneratorCreate(t, binaryPath, args...)
	require.NoError(t, err, "Generator create command failed using --template-dir")
	defer os.RemoveAll(outputDir)

	// Assertions
	// Check if a file from the local template was generated
	expectedFilePath := filepath.Join(outputDir, "test.txt")
	_, err = os.Stat(expectedFilePath)
	require.NoError(t, err, "Expected file test.txt not found in output using --template-dir")

	// Check content
	contentBytes, err := os.ReadFile(expectedFilePath)
	require.NoError(t, err, "Failed to read generated test.txt")
	require.Equal(t, "Local template var: localDefault", string(contentBytes), "Content mismatch for file generated from local template")

	t.Logf("Successfully generated mixin from local template dir: %s", localTemplateDir)
}

// TODO: Add TestCreateMixin_Integration_TemplateUrl (requires a test git repo URL)
