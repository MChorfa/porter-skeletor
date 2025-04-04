package main

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateMixinName checks if the mixin name is valid
func ValidateMixinName(name string) error {
	if name == "" {
		return fmt.Errorf("mixin name cannot be empty")
	}

	// Check if name contains only lowercase letters, numbers, and hyphens
	validName := regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("mixin name must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens")
	}

	// Check if name is a reserved word
	reservedWords := []string{
		"porter", "mixin", "mixins", "bundle", "bundles", "installation", "installations",
		"credential", "credentials", "parameter", "parameters", "claim", "claims",
		"agent", "help", "version", "schema", "build", "install", "invoke", "upgrade", "uninstall",
		// Add any other relevant reserved words
	}
	for _, reserved := range reservedWords {
		if strings.ToLower(name) == reserved {
			return fmt.Errorf("mixin name '%s' is a reserved word", name)
		}
	}

	return nil
}

// TODO: Add validation functions for other inputs if needed (e.g., AuthorName, ModulePath)
