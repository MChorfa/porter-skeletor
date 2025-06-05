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
	FeatureToggles   *FeatureToggles     `json:"feature_toggles,omitempty"`   // Enterprise feature toggle configuration
}

// FeatureToggles represents enterprise feature toggle configuration
type FeatureToggles struct {
	Security      *SecurityFeatures      `json:"security,omitempty"`
	Compliance    *ComplianceFeatures    `json:"compliance,omitempty"`
	Auth          *AuthFeatures          `json:"auth,omitempty"`
	Observability *ObservabilityFeatures `json:"observability,omitempty"`
}

// SecurityFeatures represents security-related feature toggles
type SecurityFeatures struct {
	Enabled               bool `json:"enabled"`
	InputValidation       bool `json:"input_validation"`
	RateLimiting          bool `json:"rate_limiting"`
	SecureHeaders         bool `json:"secure_headers"`
	VulnerabilityScanning bool `json:"vulnerability_scanning"`
	PolicyEnforcement     bool `json:"policy_enforcement"`
}

// ComplianceFeatures represents compliance framework feature toggles
type ComplianceFeatures struct {
	Enabled  bool                    `json:"enabled"`
	SOC2     bool                    `json:"soc2"`
	GDPR     bool                    `json:"gdpr"`
	HIPAA    bool                    `json:"hipaa"`
	PCIDSS   bool                    `json:"pci_dss"`
	Custom   map[string]bool         `json:"custom,omitempty"`
	Policies map[string]PolicyConfig `json:"policies,omitempty"`
}

// AuthFeatures represents authentication and authorization feature toggles
type AuthFeatures struct {
	Enabled      bool            `json:"enabled"`
	RBAC         bool            `json:"rbac"`
	LDAP         bool            `json:"ldap"`
	SSO          bool            `json:"sso"`
	MFA          bool            `json:"mfa"`
	Vault        bool            `json:"vault"`
	SessionMgmt  bool            `json:"session_management"`
	Integrations map[string]bool `json:"integrations,omitempty"`
}

// ObservabilityFeatures represents observability and monitoring feature toggles
type ObservabilityFeatures struct {
	Enabled        bool            `json:"enabled"`
	APM            bool            `json:"apm"`
	Infrastructure bool            `json:"infrastructure"`
	CustomMetrics  bool            `json:"custom_metrics"`
	HealthChecks   bool            `json:"health_checks"`
	OpenTelemetry  bool            `json:"opentelemetry"`
	AuditLogging   bool            `json:"audit_logging"`
	Tracing        bool            `json:"tracing"`
	Backends       map[string]bool `json:"backends,omitempty"`
}

// PolicyConfig represents configuration for compliance policies
type PolicyConfig struct {
	Enabled    bool     `json:"enabled"`
	Severity   string   `json:"severity"`
	Rules      []string `json:"rules,omitempty"`
	Exceptions []string `json:"exceptions,omitempty"`
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

// Feature toggle evaluation functions

// IsFeatureEnabled checks if a specific feature is enabled based on the feature toggle configuration
func (ft *FeatureToggles) IsFeatureEnabled(category, feature string) bool {
	if ft == nil {
		return false
	}

	switch category {
	case "security":
		return ft.isSecurityFeatureEnabled(feature)
	case "compliance":
		return ft.isComplianceFeatureEnabled(feature)
	case "auth":
		return ft.isAuthFeatureEnabled(feature)
	case "observability":
		return ft.isObservabilityFeatureEnabled(feature)
	default:
		return false
	}
}

// isSecurityFeatureEnabled checks if a security feature is enabled
func (ft *FeatureToggles) isSecurityFeatureEnabled(feature string) bool {
	if ft.Security == nil || !ft.Security.Enabled {
		return false
	}

	switch feature {
	case "input_validation":
		return ft.Security.InputValidation
	case "rate_limiting":
		return ft.Security.RateLimiting
	case "secure_headers":
		return ft.Security.SecureHeaders
	case "vulnerability_scanning":
		return ft.Security.VulnerabilityScanning
	case "policy_enforcement":
		return ft.Security.PolicyEnforcement
	default:
		return false
	}
}

// isComplianceFeatureEnabled checks if a compliance feature is enabled
func (ft *FeatureToggles) isComplianceFeatureEnabled(feature string) bool {
	if ft.Compliance == nil || !ft.Compliance.Enabled {
		return false
	}

	switch feature {
	case "soc2":
		return ft.Compliance.SOC2
	case "gdpr":
		return ft.Compliance.GDPR
	case "hipaa":
		return ft.Compliance.HIPAA
	case "pci_dss":
		return ft.Compliance.PCIDSS
	default:
		// Check custom compliance features
		if ft.Compliance.Custom != nil {
			return ft.Compliance.Custom[feature]
		}
		return false
	}
}

// isAuthFeatureEnabled checks if an authentication feature is enabled
func (ft *FeatureToggles) isAuthFeatureEnabled(feature string) bool {
	if ft.Auth == nil || !ft.Auth.Enabled {
		return false
	}

	switch feature {
	case "rbac":
		return ft.Auth.RBAC
	case "ldap":
		return ft.Auth.LDAP
	case "sso":
		return ft.Auth.SSO
	case "mfa":
		return ft.Auth.MFA
	case "vault":
		return ft.Auth.Vault
	case "session_management":
		return ft.Auth.SessionMgmt
	default:
		// Check integration-specific features
		if ft.Auth.Integrations != nil {
			return ft.Auth.Integrations[feature]
		}
		return false
	}
}

// isObservabilityFeatureEnabled checks if an observability feature is enabled
func (ft *FeatureToggles) isObservabilityFeatureEnabled(feature string) bool {
	if ft.Observability == nil || !ft.Observability.Enabled {
		return false
	}

	switch feature {
	case "apm":
		return ft.Observability.APM
	case "infrastructure":
		return ft.Observability.Infrastructure
	case "custom_metrics":
		return ft.Observability.CustomMetrics
	case "health_checks":
		return ft.Observability.HealthChecks
	case "opentelemetry":
		return ft.Observability.OpenTelemetry
	case "audit_logging":
		return ft.Observability.AuditLogging
	case "tracing":
		return ft.Observability.Tracing
	default:
		// Check backend-specific features
		if ft.Observability.Backends != nil {
			return ft.Observability.Backends[feature]
		}
		return false
	}
}

// GetEnabledFeatures returns a map of all enabled features organized by category
func (ft *FeatureToggles) GetEnabledFeatures() map[string][]string {
	enabled := make(map[string][]string)

	if ft == nil {
		return enabled
	}

	// Security features
	if ft.Security != nil && ft.Security.Enabled {
		var securityFeatures []string
		if ft.Security.InputValidation {
			securityFeatures = append(securityFeatures, "input_validation")
		}
		if ft.Security.RateLimiting {
			securityFeatures = append(securityFeatures, "rate_limiting")
		}
		if ft.Security.SecureHeaders {
			securityFeatures = append(securityFeatures, "secure_headers")
		}
		if ft.Security.VulnerabilityScanning {
			securityFeatures = append(securityFeatures, "vulnerability_scanning")
		}
		if ft.Security.PolicyEnforcement {
			securityFeatures = append(securityFeatures, "policy_enforcement")
		}
		if len(securityFeatures) > 0 {
			enabled["security"] = securityFeatures
		}
	}

	// Compliance features
	if ft.Compliance != nil && ft.Compliance.Enabled {
		var complianceFeatures []string
		if ft.Compliance.SOC2 {
			complianceFeatures = append(complianceFeatures, "soc2")
		}
		if ft.Compliance.GDPR {
			complianceFeatures = append(complianceFeatures, "gdpr")
		}
		if ft.Compliance.HIPAA {
			complianceFeatures = append(complianceFeatures, "hipaa")
		}
		if ft.Compliance.PCIDSS {
			complianceFeatures = append(complianceFeatures, "pci_dss")
		}
		// Add custom compliance features
		for feature, enabled := range ft.Compliance.Custom {
			if enabled {
				complianceFeatures = append(complianceFeatures, feature)
			}
		}
		if len(complianceFeatures) > 0 {
			enabled["compliance"] = complianceFeatures
		}
	}

	// Auth features
	if ft.Auth != nil && ft.Auth.Enabled {
		var authFeatures []string
		if ft.Auth.RBAC {
			authFeatures = append(authFeatures, "rbac")
		}
		if ft.Auth.LDAP {
			authFeatures = append(authFeatures, "ldap")
		}
		if ft.Auth.SSO {
			authFeatures = append(authFeatures, "sso")
		}
		if ft.Auth.MFA {
			authFeatures = append(authFeatures, "mfa")
		}
		if ft.Auth.Vault {
			authFeatures = append(authFeatures, "vault")
		}
		if ft.Auth.SessionMgmt {
			authFeatures = append(authFeatures, "session_management")
		}
		// Add integration-specific features
		for feature, enabled := range ft.Auth.Integrations {
			if enabled {
				authFeatures = append(authFeatures, feature)
			}
		}
		if len(authFeatures) > 0 {
			enabled["auth"] = authFeatures
		}
	}

	// Observability features
	if ft.Observability != nil && ft.Observability.Enabled {
		var obsFeatures []string
		if ft.Observability.APM {
			obsFeatures = append(obsFeatures, "apm")
		}
		if ft.Observability.Infrastructure {
			obsFeatures = append(obsFeatures, "infrastructure")
		}
		if ft.Observability.CustomMetrics {
			obsFeatures = append(obsFeatures, "custom_metrics")
		}
		if ft.Observability.HealthChecks {
			obsFeatures = append(obsFeatures, "health_checks")
		}
		if ft.Observability.OpenTelemetry {
			obsFeatures = append(obsFeatures, "opentelemetry")
		}
		if ft.Observability.AuditLogging {
			obsFeatures = append(obsFeatures, "audit_logging")
		}
		if ft.Observability.Tracing {
			obsFeatures = append(obsFeatures, "tracing")
		}
		// Add backend-specific features
		for feature, enabled := range ft.Observability.Backends {
			if enabled {
				obsFeatures = append(obsFeatures, feature)
			}
		}
		if len(obsFeatures) > 0 {
			enabled["observability"] = obsFeatures
		}
	}

	return enabled
}
