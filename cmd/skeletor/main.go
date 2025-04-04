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
	"strconv"
	"strings"
	"text/template" // Import text/template
	"time"          // Import time package

	"github.com/spf13/cobra"

	"github.com/getporter/skeletor/pkg" // Import the local pkg
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
	MixinName    string // Name of the mixin (lowercase)
	MixinNameCap string // Capitalized mixin name
	AuthorName   string // Author's name
	ModulePath   string // Go module path
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
	return cmd
}

func buildCreateCommand() *cobra.Command {
	var (
		name           string
		author         string
		modulePath     string
		outputDir      string
		nonInteractive bool
		templateUrl    string
		templateDir    string
		extraVars      []string
		dryRun         bool // Add dryRun variable
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
			data, err := buildTemplateData(config, name, author, modulePath, outputDir, nonInteractive, extraVars)
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
	"lower": strings.ToLower,
	"now":   time.Now, // Add now function
}

// createMixin generates the mixin files from the template source or simulates if dryRun is true
func createMixin(data map[string]interface{}, tmplFS fs.FS, templateRoot, outputDir string, config *TemplateConfig, dryRun bool) error { // Add dryRun parameter
	if dryRun {
		fmt.Println("[Dry Run] Simulating file generation...")
	} else {
		// Only create the output directory if not a dry run
		if err := os.MkdirAll(outputDir, 0750); err != nil { // Changed permission to 0750
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		fmt.Println("Generating mixin files...")
	}

	fmt.Println("Generating mixin files...")
	err := fs.WalkDir(tmplFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking template source at %s: %w", path, err)
		}

		// path is relative to the tmplFS root, potentially including templateRoot (e.g., "template/README.md.tmpl")
		// We need the path relative to the *template structure's root* for destination calculation.
		destRelPath := path
		if templateRoot != "." && strings.HasPrefix(path, templateRoot+"/") {
			destRelPath = strings.TrimPrefix(path, templateRoot+"/")
		} else if path == templateRoot {
			// Skip processing the root directory itself if templateRoot isn't "."
			return nil
		}

		// Skip template.json itself by checking the original path within the FS
		if path == filepath.Join(templateRoot, "template.json") {
			return nil
		}

		// Use original path for checking ignores relative to FS root
		if path == "." || strings.Contains(path, ".git") {
			return nil // Skip FS root or .git dirs
		}

		// Skip ignored files
		// Use original path for matching ignore patterns relative to FS root
		for _, pattern := range config.Ignore {
			matched, matchErr := filepath.Match(pattern, path) // Use path here
			if matchErr != nil {
				return fmt.Errorf("error matching ignore pattern %s: %w", pattern, matchErr)
			}
			if matched {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		// Determine the source path within the FS, considering conditionals
		// Use original path for source reading and conditional logic relative to FS root
		sourcePath := path    // Use path directly for source reading initially
		info, err := d.Info() // Get FileInfo from DirEntry
		if err != nil {
			return fmt.Errorf("could not get FileInfo for %s: %w", path, err)
		}

		// Conditional paths are defined relative to the template structure's root, so use destRelPath for lookup
		if sourceTemplatePathTmplStr, exists := config.ConditionalPaths[destRelPath]; exists {
			sourceTemplatePathTmpl, err := template.New("sourcePathCondition").Parse(sourceTemplatePathTmplStr)
			if err != nil {
				// Use destRelPath in error message as that's the key used
				return fmt.Errorf("failed to parse conditional source path template for destination %s: %w", destRelPath, err)
			}

			var sourcePathBuf bytes.Buffer
			if err := sourceTemplatePathTmpl.Execute(&sourcePathBuf, data); err != nil {
				// Use destRelPath in error message
				return fmt.Errorf("failed to execute conditional source path template for destination %s: %w", destRelPath, err)
			}
			evaluatedSourceRelPath := sourcePathBuf.String() // This is relative to template structure root

			if evaluatedSourceRelPath == "" {
				// Use destRelPath in message
				fmt.Printf("  Skipping destination %s (conditional source path evaluated to empty)\n", destRelPath)
				if info.IsDir() {
					return fs.SkipDir // Skip the original directory path
				}
				return nil
			}

			// Construct the actual source path within the FS by joining with templateRoot if needed
			// sourceRelPath = evaluatedSourceRelPath // Keep this relative for potential logging - REMOVED as unused
			if templateRoot != "." {
				sourcePath = filepath.Join(templateRoot, evaluatedSourceRelPath)
			} else {
				sourcePath = evaluatedSourceRelPath
			}

			// Stat the actual source path in the FS
			newInfo, statErr := fs.Stat(tmplFS, sourcePath)
			if statErr != nil {
				// Use destRelPath and sourcePath in message
				fmt.Printf("  Warning: Conditional source path %s (evaluated from %s) for destination %s does not exist in FS. Skipping.\n", sourcePath, evaluatedSourceRelPath, destRelPath)
				if info.IsDir() {
					return fs.SkipDir // Skip the original directory path
				}
				return nil
			}
			info = newInfo // Use the FileInfo of the actual source file/dir
		}

		// Process destination path using the *stripped* destRelPath template logic
		destPathTemplate, err := template.New("destPath").Parse(destRelPath)
		if err != nil {
			return fmt.Errorf("failed to parse destination path template for %s: %w", destRelPath, err)
		}
		var destPathBuf bytes.Buffer
		if err := destPathTemplate.Execute(&destPathBuf, data); err != nil {
			return fmt.Errorf("failed to execute destination path template for %s: %w", destRelPath, err)
		}
		processedDestRelPath := strings.TrimSuffix(destPathBuf.String(), ".tmpl")
		destPath := filepath.Join(outputDir, processedDestRelPath)

		// Handle directories
		if info.IsDir() {
			// Ensure we don't try to create the output directory itself if destRelPath is "."
			if processedDestRelPath == "." {
				return nil
			}
			if dryRun {
				fmt.Printf("[Dry Run] Would create directory: %s\n", destPath)
				return nil // Don't actually create in dry run
			}
			// Use 0750 for directories as recommended by gosec G301
			return os.MkdirAll(destPath, 0750)
		}

		// Handle files: Read content from sourcePath within the FS
		content, err := fs.ReadFile(tmplFS, sourcePath)
		// Use sourcePath (which might have been updated by conditional logic)
		if err != nil {
			return fmt.Errorf("failed to read source file %s from FS: %w", sourcePath, err)
		}

		// Process file content as a template using the *destination* relative path as the template name for consistency
		// Add the custom function map here
		tmpl, err := template.New(destRelPath).Funcs(funcMap).Parse(string(content))
		processedContent := string(content) // Default to raw content if not a template

		// Only execute if it's a .tmpl file or explicitly marked for processing
		// Using TrimSuffix check on processedDestRelPath vs destRelPath to see if .tmpl was removed
		if processedDestRelPath != destRelPath || strings.HasSuffix(destRelPath, ".tmpl") { // Heuristic: if paths differ after trim, it was likely a template
			if err == nil { // Check if template parsing succeeded
				var templatedContentBuf bytes.Buffer
				if err := tmpl.Execute(&templatedContentBuf, data); err != nil {
					// Use destRelPath and sourcePath in error
					return fmt.Errorf("failed to execute content template for %s (source %s): %w", destRelPath, sourcePath, err)
				}
				processedContent = templatedContentBuf.String()
			} else {
				// Log if parsing failed but maybe shouldn't halt? Or maybe it should?
				// Let's halt for now. If a file looks like a template but fails parsing, it's an issue.
				return fmt.Errorf("failed to parse content template for %s (source %s): %w", destRelPath, sourcePath, err)
			}
		} // else: treat as plain file, use raw content

		// Perform Go-specific string replacements *after* template execution
		mixinNameRaw, _ := data["MixinName"].(string)           // Use raw mixin name from data
		sanitizedName, _ := data["SanitizedMixinName"].(string) // Use sanitized name from data
		authorName, _ := data["AuthorName"].(string)
		// Use destRelPath for checking file type/location relative to template structure root
		placeholderPkgDir := filepath.Join("pkg", "mixin") // Assuming template structure root contains pkg/mixin
		placeholderCmdDir := filepath.Join("cmd", "mixin") // Assuming template structure root contains cmd/mixin

		// Check against the *destination* relative path for Go file replacements
		if strings.HasSuffix(destRelPath, ".go") {
			// Replace package declaration ONLY if it's still 'package mixin' and in the expected template dirs
			if strings.HasPrefix(destRelPath, placeholderPkgDir) || strings.HasPrefix(destRelPath, placeholderCmdDir) {
				// Use sanitizedName for package name
				processedContent = strings.ReplaceAll(processedContent, "package mixin", "package "+sanitizedName)
			}

			// Replace ONLY the specific placeholder string for author name
			processedContent = strings.ReplaceAll(processedContent, `"YOURNAME"`, `"`+authorName+`"`)

			// Replace ONLY the specific placeholder string for mixin name in comments/strings
			// Use mixinNameRaw here
			processedContent = strings.ReplaceAll(processedContent, `Use:  "mixin"`, `Use:  "`+mixinNameRaw+`"`)
			processedContent = strings.ReplaceAll(processedContent, `StartRootSpan(ctx, "mixin")`, `StartRootSpan(ctx, "`+mixinNameRaw+`")`)
		}

		// Write file or simulate if dry run
		if dryRun {
			// Use sourcePath in log message
			fmt.Printf("[Dry Run] Would write file: %s (from source %s)\n", destPath, sourcePath)
			// Optionally print content diff or summary here if needed
			return nil // Don't actually write in dry run
		}
		// Use 0600 for files as recommended by gosec G306 (owner rw only)
		return os.WriteFile(destPath, []byte(processedContent), 0600)
	})

	if err != nil {
		return err
	}

	// --- Post Generation Validation ---
	if dryRun {
		fmt.Println("\n[Dry Run] Skipping post-generation validation.")
		fmt.Println("[Dry Run] Would run the following validation commands:")
		fmt.Println("  - go mod tidy")
		fmt.Println("  - go build ./...")
		fmt.Println("  - go test ./...")
		fmt.Println("\n[Dry Run] Simulation complete.") // Final dry run message
		return nil
	}

	fmt.Println("\nRunning post-generation validation...")
	if err := runCommandInDir(outputDir, "go", "mod", "tidy"); err != nil {
		fmt.Printf("Warning: 'go mod tidy' failed: %v\n", err)
	} else {
		fmt.Println("  - go mod tidy: OK")
	}
	if err := runCommandInDir(outputDir, "go", "build", "./..."); err != nil {
		fmt.Printf("Warning: 'go build ./...' failed: %v\n", err)
	} else {
		fmt.Println("  - go build ./...: OK")
	}
	if err := runCommandInDir(outputDir, "go", "test", "./..."); err != nil {
		fmt.Printf("Warning: 'go test ./...' failed: %v\n", err)
	} else {
		fmt.Println("  - go test ./...: OK")
	}

	fmt.Println("\nValidation complete.")
	// Next steps message moved to RunE for non-dry run case

	return nil
}

// --- Other helper functions ---

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

func buildTemplateData(config *TemplateConfig, name, author, modulePath, outputDir string, nonInteractive bool, extraVars []string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
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
	return data, nil
}
