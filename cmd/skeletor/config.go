package main

import (
	"bytes"
	"encoding/json"
	"errors" // Import errors
	"fmt"
	"io/fs" // Import io/fs
	"os"
	"os/exec" // Import os/exec
	"path/filepath"
	"runtime" // Import runtime
	"strings"
	"text/template"
)

// TemplateConfig represents the configuration for a template
type TemplateConfig struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Variables        map[string]Variable `json:"variables"`
	Hooks            map[string][]string `json:"hooks"`
	Ignore           []string            `json:"ignore"`
	ConditionalPaths map[string]string   `json:"conditional_paths,omitempty"` // Map of relative path -> Go template condition string
}

// Variable represents a template variable with its properties
type Variable struct {
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Type        string      `json:"type,omitempty"` // string, bool, int, etc.
	Required    bool        `json:"required,omitempty"`
	Choices     []string    `json:"choices,omitempty"` // For enum-like variables
}

// LoadTemplateConfig loads the template configuration from the given filesystem and root directory
func LoadTemplateConfig(tmplFS fs.FS, templateRoot string) (*TemplateConfig, error) {
	configPath := filepath.Join(templateRoot, "template.json") // Path within the FS

	// Check if config file exists within the FS using fs.Stat
	if _, err := fs.Stat(tmplFS, configPath); errors.Is(err, fs.ErrNotExist) {
		// Return default config if no config file exists in the source FS
		fmt.Println("Warning: template.json not found in template source, using default configuration.")
		return &TemplateConfig{
			Name:        "Porter Mixin Template (Default)", // Indicate default
			Description: "Default Porter mixin template",
			Variables: map[string]Variable{
				"MixinName": {
					Description: "Name of the mixin (lowercase)",
					Type:        "string",
					Required:    true,
				},
				"AuthorName": {
					Description: "Author name",
					Type:        "string",
					Required:    true,
				},
				"ModulePath": {
					Description: "Go module path",
					Type:        "string",
					Default:     "github.com/getporter/{{ .MixinName }}",
				},
				// ComplianceLevel is now defined in template.json
				"MixinFeedRepoURL": { // Add MixinFeedRepoURL variable
					Description: "Git URL for the mixin feed repository (e.g., git@github.com:YOUR/packages.git)",
					Type:        "string",
					Required:    false, // Optional
				},
				"MixinFeedBranch": { // Add MixinFeedBranch variable
					Description: "Branch in the mixin feed repository to commit to",
					Type:        "string",
					Default:     "main",
					Required:    false,
				},
				"AuthorEmail": { // Add AuthorEmail variable
					Description: "Author's email for security contact",
					Type:        "string",
					Required:    false, // Optional, defaults in security.txt
				},
			},
		}, nil
	} else if err != nil {
		// Other error during stat
		return nil, fmt.Errorf("failed to stat template config %s: %w", configPath, err)
	}

	// Read and parse config file from the FS using fs.ReadFile
	data, err := fs.ReadFile(tmplFS, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template config %s from FS: %w", configPath, err)
	}

	var config TemplateConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return &config, nil
}

// RunHooks executes the hooks defined in the template configuration, substituting variables
func RunHooks(config *TemplateConfig, hookName string, outputDir string, data map[string]interface{}) error {
	hooks, exists := config.Hooks[hookName]
	if !exists || len(hooks) == 0 {
		return nil // No hooks for this stage
	}

	fmt.Printf("Running %s hooks...\n", hookName)
	for _, commandTmplStr := range hooks {
		// Process command string as a template
		commandTmpl, err := template.New("hook-cmd").Parse(commandTmplStr)
		if err != nil {
			// If parsing fails, treat it as a literal command for backward compatibility? Or error out?
			// Let's error out for now to encourage proper templating.
			return fmt.Errorf("failed to parse hook command template '%s': %w", commandTmplStr, err)
		}

		var commandBuf bytes.Buffer // Need to import "bytes"
		if err := commandTmpl.Execute(&commandBuf, data); err != nil {
			return fmt.Errorf("failed to execute hook command template '%s': %w", commandTmplStr, err)
		}
		processedCommand := commandBuf.String()

		// Split command into executable and args (basic split, might need refinement for complex cases)
		parts := strings.Fields(processedCommand)
		if len(parts) == 0 {
			continue // Skip empty commands
		}
		executable := parts[0]
		args := parts[1:]

		fmt.Printf("  Executing: %s\n", processedCommand)
		cmd := createCommand(executable, args...) // Use the helper function
		cmd.Dir = outputDir                       // Run in the generated directory
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook '%s' failed: %w", processedCommand, err)
		}
	}
	fmt.Printf("%s hooks completed.\n", hookName)

	return nil
}

// Allowed commands for hooks
var allowedHookCommands = map[string]bool{
	"go":  true,
	"git": true,
	// Add other safe commands here if needed, e.g., "echo", "mage"
}

// Helper function to create OS-specific commands, checking against an allow-list
func createCommand(name string, args ...string) *exec.Cmd {
	// Check if the command is allowed
	if !allowedHookCommands[name] {
		// Return a command that will error out immediately
		// This prevents execution of arbitrary commands from templates
		// We use "false" which is a standard shell utility that just exits with status 1
		fmt.Fprintf(os.Stderr, "Error: Hook command '%s' is not allowed.\n", name)
		return exec.Command("false") // "false" command always fails
	}

	// Handle potential issue with 'go' command on windows needing cmd /c
	if runtime.GOOS == "windows" && (name == "git" || name == "go") {
		// #nosec G204 -- Command is restricted by allow-list, args come from trusted template or user input
		return exec.Command("cmd", append([]string{"/c", name}, args...)...)
	}
	// #nosec G204 -- Command is restricted by allow-list, args come from trusted template or user input
	return exec.Command(name, args...)
}

// Removed duplicate buildTemplateData function. The correct one is in main.go.
