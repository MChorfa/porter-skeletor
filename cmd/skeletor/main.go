package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp" // Ensure regexp is imported
	"strconv"
	"strings"
	"text/template" // Import text/template
	"time"          // Import time package

	"github.com/spf13/cobra"

	"github.com/getporter/skeletor/pkg" // Import the local pkg
)

// Version information (set by build flags)
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

var workingDir string

func init() {
	// Get the executable's directory
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workingDir = filepath.Dir(filepath.Dir(ex)) // Go up two levels from bin/porter-mixin-generator
}

// TemplateData holds the variables to be replaced in templates
type TemplateData struct {
	MixinName       string // Name of the mixin (lowercase)
	MixinNameCap    string // Capitalized mixin name
	AuthorName      string // Author's name
	ModulePath      string // Go module path
	ComplianceLevel string // Compliance level (basic, standard, advanced)
}

func main() {
	cmd := buildRootCommand()
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skeletor",
		Short: "Create new Porter mixins easily",
	}

	cmd.AddCommand(buildCreateCommand())
	cmd.AddCommand(buildVersionCommand())
	return cmd
}

func buildVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Skeletor %s\n", Version)
			fmt.Printf("Commit: %s\n", Commit)
			fmt.Printf("Built: %s\n", Date)
		},
	}
}

func buildCreateCommand() *cobra.Command {
	var (
		name                  string
		author                string
		modulePath            string
		outputDir             string
		nonInteractive        bool
		templateUrl           string
		templateDir           string
		extraVars             []string
		dryRun                bool   // Add dryRun variable
		complianceLevel       string // Declare complianceLevel
		enableSecurity        bool
		enableCompliance      bool
		enableAuth            bool
		enableObservability   bool
		securityFeatures      string
		complianceFrameworks  string
		authFeatures          string
		observabilityFeatures string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Porter mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine the template source (embed, url, local dir)
			tmplFS, rootDirForWalk, cleanupDir, err := getTemplateSource(templateUrl, templateDir)
			if err != nil {
				return err
			}
			// If cloned from URL, ensure cleanup
			if cleanupDir != "" {
				defer os.RemoveAll(cleanupDir)
			}

			// Load template configuration from the determined source
			config, err := LoadTemplateConfig(tmplFS, rootDirForWalk) // Use rootDirForWalk here
			if err != nil {
				return fmt.Errorf("failed to load template config from %s: %w", rootDirForWalk, err)
			}

			// Create template data from config and flags
			// Pass complianceLevel and feature toggles to buildTemplateData
			data, err := buildTemplateDataWithFeatures(config, name, author, modulePath, outputDir, complianceLevel, nonInteractive, extraVars,
				enableSecurity, enableCompliance, enableAuth, enableObservability,
				securityFeatures, complianceFrameworks, authFeatures, observabilityFeatures)
			if err != nil {
				return err
			}

			// Get the final output directory path from the template data
			finalOutputDir, ok := data["OutputDir"].(string)
			if !ok || finalOutputDir == "" {
				// Fallback or error if OutputDir isn't set correctly in data
				// This shouldn't happen with current logic but good practice to check
				mixinName, _ := data["MixinName"].(string) // Assume MixinName exists
				if mixinName == "" {
					return fmt.Errorf("output directory could not be determined: MixinName is missing")
				}
				finalOutputDir = "./" + mixinName
				fmt.Printf("Warning: OutputDir not found in template data, defaulting to %s\n", finalOutputDir)
			}

			// Create mixin from template using the determined source FS and root
			// Pass dryRun variable and the finalOutputDir
			if err := createMixin(data, tmplFS, rootDirForWalk, finalOutputDir, config, dryRun); err != nil {
				return err
			}

			// Run post-generation hooks or simulate if dry run
			if dryRun {
				fmt.Println("\n[Dry Run] Skipping post-generation hooks.")
				if hooks, exists := config.Hooks["post_gen"]; exists && len(hooks) > 0 {
					fmt.Println("[Dry Run] Would run the following hooks:")
					for _, hookCmd := range hooks {
						// Attempt to render hook command for better output, ignore errors
						tmpl, err := template.New("hook-dry-run").Parse(hookCmd)
						renderedCmd := hookCmd // Default to raw if template fails
						if err == nil {
							var buf bytes.Buffer
							if tmpl.Execute(&buf, data) == nil {
								renderedCmd = buf.String()
							}
						}
						fmt.Printf("  - %s\n", renderedCmd)
					}
				}
				// Final dry run message moved to end of createMixin when dryRun is true
				return nil // Exit successfully after dry run simulation in createMixin
			} else {
				// Only run hooks if not a dry run
				// Use finalOutputDir for running hooks
				if err := RunHooks(config, "post_gen", finalOutputDir, data); err != nil {
					return err // Return hook errors if they occur
				}
				// Print next steps only on successful non-dry run
				fmt.Println("\nNext steps:")
				fmt.Println("1. cd", finalOutputDir) // Use finalOutputDir here too
				fmt.Println("2. Review the generated code and customize as needed.")
				fmt.Println("3. Run 'mage build test' or 'go run ./ci' for further verification.")
				return nil
			}
		},
	}

	// Add flags (same as before)
	cmd.Flags().StringVar(&name, "name", "", "Name of the new mixin (lowercase)")
	cmd.Flags().StringVar(&author, "author", "", "Author name for the mixin")
	cmd.Flags().StringVar(&modulePath, "module", "", "Go module path (default: github.com/getporter/<name>)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory (default: ./<name>)")
	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Run in non-interactive mode")
	cmd.Flags().StringVar(&templateUrl, "template-url", "", "URL to a git repository containing the template")
	cmd.Flags().StringVar(&templateDir, "template-dir", "", "Local directory containing the template")
	cmd.Flags().StringArrayVar(&extraVars, "var", []string{}, "Extra variables in KEY=VALUE format")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate generation without writing files") // Add dry-run flag
	// Use "basic" as default to match template.json, ensure choices match template.json
	cmd.Flags().StringVar(&complianceLevel, "compliance-level", "basic", "Compliance level (basic, slsa-l1, slsa-l3)")

	// Enterprise feature toggle flags
	cmd.Flags().BoolVar(&enableSecurity, "enable-security", false, "Enable enterprise security features")
	cmd.Flags().BoolVar(&enableCompliance, "enable-compliance", false, "Enable compliance framework features")
	cmd.Flags().BoolVar(&enableAuth, "enable-auth", false, "Enable authentication and authorization features")
	cmd.Flags().BoolVar(&enableObservability, "enable-observability", false, "Enable enhanced observability features")

	cmd.Flags().StringVar(&securityFeatures, "security-features", "", "Comma-separated list of security features (input_validation,rate_limiting,secure_headers,vulnerability_scanning,policy_enforcement)")
	cmd.Flags().StringVar(&complianceFrameworks, "compliance-frameworks", "", "Comma-separated list of compliance frameworks (soc2,gdpr,hipaa,pci_dss)")
	cmd.Flags().StringVar(&authFeatures, "auth-features", "", "Comma-separated list of auth features (rbac,ldap,sso,mfa,vault,session_management)")
	cmd.Flags().StringVar(&observabilityFeatures, "observability-features", "", "Comma-separated list of observability features (apm,infrastructure,custom_metrics,health_checks,opentelemetry,audit_logging,tracing)")

	return cmd
}

// getTemplateSource determines the source filesystem and the root path for walking.
// Returns fs.FS, root path for WalkDir, path to cleanup (if any), error.
func getTemplateSource(templateUrl, templateDir string) (fs.FS, string, string, error) {
	// Priority: Local Directory > URL > Embedded
	if templateDir != "" {
		fileInfo, err := os.Stat(templateDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, "", "", fmt.Errorf("template directory does not exist: %s", templateDir)
			}
			return nil, "", "", fmt.Errorf("failed to stat template directory %s: %w", templateDir, err)
		}
		if !fileInfo.IsDir() {
			return nil, "", "", fmt.Errorf("template path is not a directory: %s", templateDir)
		}
		fmt.Printf("Using local template directory: %s\n", templateDir)
		return os.DirFS(templateDir), ".", "", nil // Root is "." relative to the DirFS, no cleanup needed
	}

	if templateUrl != "" {
		tempDir, err := os.MkdirTemp("", "porter-template-*")
		if err != nil {
			return nil, "", "", fmt.Errorf("failed to create temp directory: %w", err)
		}

		fmt.Printf("Fetching template from %s...\n", templateUrl)
		// #nosec G204 -- URL is from user flag, tempDir is generated, command is allow-listed
		cmd := createCommand("git", "clone", "--depth=1", templateUrl, tempDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// Attempt to clean up temp dir even on clone failure, but ignore cleanup error
			_ = os.RemoveAll(tempDir)
			return nil, "", "", fmt.Errorf("failed to clone template repository: %w", err)
		}
		fmt.Println("Using cloned template repository.")
		// Return the OS FS rooted at the temp dir, root is ".", cleanup path is tempDir
		return os.DirFS(tempDir), ".", tempDir, nil
	}

	// Default to the embedded templates
	fmt.Println("Using embedded templates.")
	// Return the embedded FS and specify "template" as the root directory within it
	return pkg.MixinTemplateFS, "template", "", nil
}

// Define custom template functions
var funcMap = template.FuncMap{
	"lower":          strings.ToLower,
	"now":            time.Now, // Add now function
	"hasFeature":     hasFeature,
	"splitFeatures":  splitFeatures,
	"joinFeatures":   joinFeatures,
	"featureEnabled": featureEnabled,
	"default":        defaultValue,
}

// Template functions for feature toggle support

// hasFeature checks if a specific feature is present in a comma-separated list
func hasFeature(featureList, feature string) bool {
	if featureList == "" {
		return false
	}
	features := strings.Split(featureList, ",")
	for _, f := range features {
		if strings.TrimSpace(f) == feature {
			return true
		}
	}
	return false
}

// splitFeatures splits a comma-separated feature list into a slice
func splitFeatures(featureList string) []string {
	if featureList == "" {
		return []string{}
	}
	features := strings.Split(featureList, ",")
	var result []string
	for _, f := range features {
		if trimmed := strings.TrimSpace(f); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// joinFeatures joins a slice of features into a comma-separated string
func joinFeatures(features []string) string {
	return strings.Join(features, ",")
}

// featureEnabled checks if a feature is enabled based on template data
func featureEnabled(data map[string]interface{}, category, feature string) bool {
	switch category {
	case "security":
		if enabled, ok := data["EnableSecurity"].(bool); !ok || !enabled {
			return false
		}
		if featureList, ok := data["SecurityFeatures"].(string); ok {
			return hasFeature(featureList, feature)
		}
	case "compliance":
		if enabled, ok := data["EnableCompliance"].(bool); !ok || !enabled {
			return false
		}
		if featureList, ok := data["ComplianceFrameworks"].(string); ok {
			return hasFeature(featureList, feature)
		}
	case "auth":
		if enabled, ok := data["EnableAuth"].(bool); !ok || !enabled {
			return false
		}
		if featureList, ok := data["AuthFeatures"].(string); ok {
			return hasFeature(featureList, feature)
		}
	case "observability":
		if enabled, ok := data["EnableObservability"].(bool); !ok || !enabled {
			return false
		}
		if featureList, ok := data["ObservabilityFeatures"].(string); ok {
			return hasFeature(featureList, feature)
		}
	}
	return false
}

// defaultValue provides a default value if the input is empty or nil
func defaultValue(defaultVal, value interface{}) interface{} {
	if value == nil {
		return defaultVal
	}

	// Handle string values
	if str, ok := value.(string); ok {
		if str == "" {
			return defaultVal
		}
		return str
	}

	// For other types, return the value if it's not nil
	return value
}

// --- Refactored createMixin and Helper Functions ---

// createMixin generates the mixin files from the template source or simulates if dryRun is true
func createMixin(data map[string]interface{}, tmplFS fs.FS, templateRoot, outputDir string, config *TemplateConfig, dryRun bool) error {
	if dryRun {
		fmt.Println("[Dry Run] Simulating file generation...")
	} else {
		// Use 0750 permission as recommended by gosec G301
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
		fmt.Println("Generating mixin files...")
	}

	err := fs.WalkDir(tmplFS, templateRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking template source at %s: %w", path, walkErr)
		}

		// Calculate destination path and check if the file/dir should be skipped
		destRelPath, skip := calculateDestPath(path, templateRoot, config.Ignore)
		if skip {
			// Special case: if this is the template root directory itself, don't return fs.SkipDir
			// because that would skip the entire tree. Just continue to the next iteration.
			if path == templateRoot && d.IsDir() {
				return nil // Continue walking the directory contents
			}
			if d.IsDir() {
				return fs.SkipDir // Skip ignored directories entirely
			}
			return nil // Skip ignored files
		}

		// Determine the actual source path and file info, handling conditional logic
		sourcePath, info, skip, err := determineSourcePath(tmplFS, path, destRelPath, templateRoot, config.ConditionalPaths, data)
		if err != nil {
			return err // Propagate errors from conditional path processing
		}
		if skip {
			if info != nil && info.IsDir() { // Need info to check if it was a dir
				return fs.SkipDir
			}
			return nil
		}

		// Process the final destination path using template data
		finalDestPath, err := processDestPath(destRelPath, outputDir, data)
		if err != nil {
			return err
		}
		if finalDestPath == "" { // Skip if templating resulted in empty path
			if info.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Handle directory or file processing
		if dryRun {
			if info.IsDir() {
				// Print for all directories except the root being walked (which is skipped earlier)
				fmt.Printf("[Dry Run] Would create directory: %s\n", finalDestPath)
			} else {
				fmt.Printf("[Dry Run] Would write file: %s (from source %s)\n", finalDestPath, sourcePath)
			}
			return nil // Skip actual processing in dry run mode
		}

		// Actual processing if not dry run
		if info.IsDir() {
			return processDirectory(finalDestPath, info)
		} else {
			return processTemplateFile(tmplFS, sourcePath, finalDestPath, info, data)
		}
	})

	if err != nil {
		return err
	}

	// Run post-generation validation if not a dry run
	if !dryRun {
		return runPostGenerationValidation(outputDir)
	}

	fmt.Println("\n[Dry Run] Simulation complete.")
	return nil
}

// calculateDestPath determines the relative destination path and if it should be skipped.
func calculateDestPath(originalPath, templateRoot string, ignorePatterns []string) (destRelPath string, skip bool) {
	// Calculate path relative to the template structure's root
	destRelPath = originalPath
	if templateRoot != "." && strings.HasPrefix(originalPath, templateRoot+"/") {
		destRelPath = strings.TrimPrefix(originalPath, templateRoot+"/")
	} else if originalPath == templateRoot || originalPath == "." {
		// Skip the root directory itself (".") or the specified templateRoot
		return "", true
	} else {
		// Path doesn't match expected structure, treat as relative path directly
		destRelPath = originalPath
	}

	// Skip template.json (check against original path within FS)
	if originalPath == filepath.Join(templateRoot, "template.json") {
		return "", true
	}

	// Skip .git directory
	if strings.Contains(originalPath, ".git") {
		return "", true
	}

	// Skip ignored files/dirs based on patterns (using original path relative to FS root)
	for _, pattern := range ignorePatterns {
		matched, _ := filepath.Match(pattern, originalPath) // Ignore match error for simplicity here
		if matched {
			return "", true
		}
	}

	return destRelPath, false
}

// determineSourcePath finds the correct source path and info, handling conditional logic.
func determineSourcePath(tmplFS fs.FS, originalPath, destRelPath, templateRoot string, conditionalPaths map[string]string, data map[string]interface{}) (sourcePath string, fileInfo fs.FileInfo, skip bool, err error) {
	sourcePath = originalPath // Default source is the original path walked

	// Get initial FileInfo using the original path
	var initialFileInfo fs.FileInfo
	initialFileInfo, err = fs.Stat(tmplFS, sourcePath)
	if err != nil {
		err = fmt.Errorf("could not get FileInfo for %s: %w", sourcePath, err)
		return
	}
	fileInfo = initialFileInfo // Use this unless overridden by conditional logic

	// Check conditional paths (key is relative to template structure root, which matches destRelPath)
	if sourceTemplatePathTmplStr, exists := conditionalPaths[destRelPath]; exists {
		sourceTemplatePathTmpl, parseErr := template.New("sourcePathCondition").Parse(sourceTemplatePathTmplStr)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse conditional source path template for destination %s: %w", destRelPath, parseErr)
			return
		}

		var sourcePathBuf bytes.Buffer
		if execErr := sourceTemplatePathTmpl.Execute(&sourcePathBuf, data); execErr != nil {
			err = fmt.Errorf("failed to execute conditional source path template for destination %s: %w", destRelPath, execErr)
			return
		}
		evaluatedSourceRelPath := sourcePathBuf.String()

		if evaluatedSourceRelPath == "" {
			fmt.Printf("  Skipping destination %s (conditional source path evaluated to empty)\n", destRelPath)
			skip = true
			return // Return original fileInfo in case caller needs to check IsDir for fs.SkipDir
		}

		// Construct the actual source path within the FS
		if templateRoot != "." {
			sourcePath = filepath.Join(templateRoot, evaluatedSourceRelPath)
		} else {
			// When templateRoot is ".", evaluatedSourceRelPath is relative to the embedded root,
			// but needs to be prefixed with "template/" to match the actual embedded path.
			// This assumes conditional paths always resolve to something inside "template/".
			sourcePath = filepath.Join("template", evaluatedSourceRelPath)
		}

		// Stat the *actual* source path
		newInfo, statErr := fs.Stat(tmplFS, sourcePath)
		if statErr != nil {
			fmt.Printf("  Warning: Conditional source path %s (evaluated from %s) for destination %s does not exist in FS. Skipping.\n", sourcePath, evaluatedSourceRelPath, destRelPath)
			skip = true
			err = nil // Treat as skip, not error
			return    // Return initialFileInfo
		}
		fileInfo = newInfo // Update fileInfo ONLY if conditional path is valid
	}
	// If no conditional path matched or was processed, fileInfo remains initialFileInfo
	return
}

// processDestPath processes the relative destination path with template data.
func processDestPath(destRelPath, outputDir string, data map[string]interface{}) (string, error) {
	destPathTemplate, err := template.New("destPath").Parse(destRelPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse destination path template for %s: %w", destRelPath, err)
	}
	var destPathBuf bytes.Buffer
	if err := destPathTemplate.Execute(&destPathBuf, data); err != nil {
		return "", fmt.Errorf("failed to execute destination path template for %s: %w", destRelPath, err)
	}
	processedDestRelPath := strings.TrimSuffix(destPathBuf.String(), ".tmpl")

	if processedDestRelPath == "" {
		fmt.Printf("  Skipping empty destination path derived from %s\n", destRelPath)
		return "", nil // Return empty string to indicate skip
	}

	return filepath.Join(outputDir, processedDestRelPath), nil
}

// processDirectory handles directory creation during mixin generation.
// Removed dryRun parameter
func processDirectory(destPath string, info fs.FileInfo) error {
	// Use 0750 for directories as recommended by gosec G301
	if err := os.MkdirAll(destPath, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destPath, err)
	}
	return nil
}

// processTemplateFile handles reading, templating, and writing a single file.
// Removed dryRun parameter
func processTemplateFile(tmplFS fs.FS, sourcePath, destPath string, info fs.FileInfo, data map[string]interface{}) error {
	// Read source content
	content, err := fs.ReadFile(tmplFS, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file %s from FS: %w", sourcePath, err)
	}

	processedContent := string(content)
	destRelPathForTemplateName := filepath.Base(destPath) // Use filename part for template name

	// Process as template only if it had .tmpl extension
	if strings.HasSuffix(sourcePath, ".tmpl") { // Check original source path for .tmpl
		tmpl, parseErr := template.New(destRelPathForTemplateName).Funcs(funcMap).Parse(string(content))
		if parseErr != nil {
			return fmt.Errorf("failed to parse content template for %s (source %s): %w", destRelPathForTemplateName, sourcePath, parseErr)
		}
		var templatedContentBuf bytes.Buffer
		if execErr := tmpl.Execute(&templatedContentBuf, data); execErr != nil {
			return fmt.Errorf("failed to execute content template for %s (source %s): %w", destRelPathForTemplateName, sourcePath, execErr)
		}
		processedContent = templatedContentBuf.String()
	} // End of if strings.HasSuffix(sourcePath, ".tmpl")

	// Apply Go-specific replacements (use destRelPathForTemplateName which is just the filename)
	processedContent = applyGoSpecificReplacements(processedContent, destRelPathForTemplateName, data)

	// Write the final content
	// Use 0600 for files as recommended by gosec G306 (owner rw only)
	if err := os.WriteFile(destPath, []byte(processedContent), 0600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}

// applyGoSpecificReplacements performs string replacements specific to Go files.
func applyGoSpecificReplacements(content, destRelPath string, data map[string]interface{}) string {
	if !strings.HasSuffix(destRelPath, ".go") {
		return content // Only process .go files
	}

	mixinNameRaw, _ := data["MixinName"].(string)
	authorName, _ := data["AuthorName"].(string)
	// Check against the *destination* relative path (filename only)
	// This logic might be too simple if templates generate files outside expected dirs
	// placeholderPkgDir := "pkg/mixin" // Relative path check - Simplified logic below
	// placeholderCmdDir := "cmd/mixin" // Relative path check - Simplified logic below

	// Determine if the file *should* be in the 'mixin' package based on its eventual path
	// This requires knowing the full intended relative path, not just the filename.
	// Let's assume for now that if it's a .go file, it should be package mixin.
	// A better approach might involve passing the full destRelPath from createMixin.
	// For now, we'll apply the package replacement more broadly.

	// Ensure package is always 'package mixin' for any generated .go file
	// This might be too broad if templates generate Go files not intended to be in package mixin.
	if !strings.Contains(content, "package mixin") {
		packageLineRegex := regexp.MustCompile(`^package\s+\w+`)
		content = packageLineRegex.ReplaceAllString(content, "package mixin")
	}

	// Replace placeholders
	content = strings.ReplaceAll(content, `"YOURNAME"`, `"`+authorName+`"`)
	content = strings.ReplaceAll(content, `Use:  "mixin"`, `Use:  "`+mixinNameRaw+`"`)                             // Use raw name for user-facing strings
	content = strings.ReplaceAll(content, `StartRootSpan(ctx, "mixin")`, `StartRootSpan(ctx, "`+mixinNameRaw+`")`) // Use raw name for tracing

	return content
}

// runPostGenerationValidation executes validation commands in the output directory.
func runPostGenerationValidation(outputDir string) error {
	fmt.Println("\nRunning post-generation validation...")
	commands := [][]string{
		{"go", "mod", "tidy"},
		{"go", "build", "./..."},
		{"go", "test", "./..."},
	}

	for _, cmdArgs := range commands {
		if err := runCommandInDir(outputDir, cmdArgs[0], cmdArgs[1:]...); err != nil {
			// Log warning but continue validation
			fmt.Printf("Warning: '%s' failed: %v\n", strings.Join(cmdArgs, " "), err)
		} else {
			fmt.Printf("  - %s: OK\n", strings.Join(cmdArgs, " "))
		}
	}
	fmt.Println("\nValidation complete.")
	return nil // Don't return error from validation failures, just warn
}

// --- Other helper functions (promptString, promptStringWithDefault, capitalize, runCommandInDir, buildTemplateData) remain the same ---

func promptString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptStringWithDefault(prompt, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%s] ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func runCommandInDir(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("  Running '%s %s' in %s...\n", command, strings.Join(args, " "), dir)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command '%s %s' failed: %w", command, strings.Join(args, " "), err)
	}
	return nil
}

// Update buildTemplateData signature and logic
func buildTemplateData(config *TemplateConfig, name, author, modulePath, outputDir, complianceLevel string, nonInteractive bool, extraVars []string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Add compliance level first so it can be used in default value templates
	data["ComplianceLevel"] = complianceLevel

	for _, varDef := range extraVars {
		parts := strings.SplitN(varDef, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid variable format: %s (expected KEY=VALUE)", varDef)
		}
		data[parts[0]] = parts[1]
	}
	if name != "" {
		data["MixinName"] = name
	}
	if author != "" {
		data["AuthorName"] = author
	}
	if modulePath != "" {
		data["ModulePath"] = modulePath
	}

	for varName, varConfig := range config.Variables {
		if _, exists := data[varName]; exists {
			continue
		}
		if !nonInteractive {
			var defaultValue string
			if varConfig.Default != nil {
				defaultValue = fmt.Sprintf("%v", varConfig.Default)
				if strings.Contains(defaultValue, "{{") {
					tmpl, err := template.New("default").Parse(defaultValue)
					if err == nil {
						var buf bytes.Buffer
						if err := tmpl.Execute(&buf, data); err == nil {
							defaultValue = buf.String()
						}
					}
				}
			}
			for {
				prompt := varConfig.Description
				if varConfig.Choices != nil && len(varConfig.Choices) > 0 {
					prompt = fmt.Sprintf("%s %v", prompt, varConfig.Choices)
				}
				var rawValue string
				if defaultValue != "" {
					rawValue = promptStringWithDefault(prompt+": ", defaultValue)
				} else {
					rawValue = promptString(prompt + ": ")
				}
				if varConfig.Choices != nil && len(varConfig.Choices) > 0 {
					isValidChoice := false
					for _, choice := range varConfig.Choices {
						if rawValue == choice {
							isValidChoice = true
							break
						}
					}
					if !isValidChoice {
						fmt.Printf("  Error: Invalid choice. Please select one of %v\n", varConfig.Choices)
						continue
					}
				}
				var validatedValue interface{}
				var validationErr error
				switch strings.ToLower(varConfig.Type) {
				case "bool", "boolean":
					validatedValue, validationErr = strconv.ParseBool(rawValue)
					if validationErr != nil {
						validationErr = fmt.Errorf("invalid boolean value (try true/false, 1/0)")
					}
				case "int", "integer":
					validatedValue, validationErr = strconv.Atoi(rawValue)
					if validationErr != nil {
						validationErr = fmt.Errorf("invalid integer value")
					}
				case "string", "":
					validatedValue = rawValue
				default:
					fmt.Printf("  Warning: Unknown variable type '%s' for '%s', treating as string.\n", varConfig.Type, varName)
					validatedValue = rawValue
				}
				if validationErr != nil {
					fmt.Printf("  Error: %v\n", validationErr)
					continue
				}
				data[varName] = validatedValue
				break
			}
		} else if varConfig.Default != nil {
			data[varName] = varConfig.Default
		} else if varConfig.Required {
			return nil, fmt.Errorf("required variable %s is not provided", varName)
		}
	}
	if mixinName, ok := data["MixinName"].(string); ok {
		sanitizedName := strings.ReplaceAll(mixinName, "-", "")
		data["MixinName"] = mixinName              // Keep original for paths and module name
		data["SanitizedMixinName"] = sanitizedName // For Go package names
		data["MixinNameCap"] = capitalize(mixinName)
		if outputDir == "" {
			data["OutputDir"] = "./" + mixinName
		} else {
			data["OutputDir"] = outputDir
		}
		if _, exists := data["ModulePath"]; !exists {
			data["ModulePath"] = fmt.Sprintf("github.com/getporter/%s", mixinName)
		}
	}
	for varName, varConfig := range config.Variables {
		if varConfig.Required {
			if _, exists := data[varName]; !exists || data[varName] == "" {
				return nil, fmt.Errorf("required variable %s is not provided", varName)
			}
		}
	}

	// Always add enterprise feature flags (default to false for backward compatibility)
	if _, exists := data["EnableSecurity"]; !exists {
		data["EnableSecurity"] = false
	}
	if _, exists := data["EnableCompliance"]; !exists {
		data["EnableCompliance"] = false
	}
	if _, exists := data["EnableAuth"]; !exists {
		data["EnableAuth"] = false
	}
	if _, exists := data["EnableObservability"]; !exists {
		data["EnableObservability"] = false
	}

	// Always add enterprise feature lists (default to empty)
	if _, exists := data["SecurityFeatures"]; !exists {
		data["SecurityFeatures"] = ""
	}
	if _, exists := data["ComplianceFrameworks"]; !exists {
		data["ComplianceFrameworks"] = ""
	}
	if _, exists := data["AuthFeatures"]; !exists {
		data["AuthFeatures"] = ""
	}
	if _, exists := data["ObservabilityFeatures"]; !exists {
		data["ObservabilityFeatures"] = ""
	}

	return data, nil
}

// buildTemplateDataWithFeatures creates template data with enterprise feature toggles
func buildTemplateDataWithFeatures(config *TemplateConfig, name, author, modulePath, outputDir, complianceLevel string, nonInteractive bool, extraVars []string,
	enableSecurity, enableCompliance, enableAuth, enableObservability bool,
	securityFeatures, complianceFrameworks, authFeatures, observabilityFeatures string) (map[string]interface{}, error) {

	// First build the base template data
	data, err := buildTemplateData(config, name, author, modulePath, outputDir, complianceLevel, nonInteractive, extraVars)
	if err != nil {
		return nil, err
	}

	// Add enterprise feature toggle data
	data["EnableSecurity"] = enableSecurity
	data["EnableCompliance"] = enableCompliance
	data["EnableAuth"] = enableAuth
	data["EnableObservability"] = enableObservability

	data["SecurityFeatures"] = securityFeatures
	data["ComplianceFrameworks"] = complianceFrameworks
	data["AuthFeatures"] = authFeatures
	data["ObservabilityFeatures"] = observabilityFeatures

	return data, nil
}

// End of file
