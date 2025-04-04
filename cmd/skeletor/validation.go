package main

import (
	"fmt"
	"strings"
)

// ValidateComplianceLevel checks if the provided compliance level is valid
func ValidateComplianceLevel(level string) error {
	validLevels := []string{"basic", "standard", "advanced"}
	for _, valid := range validLevels {
		if strings.ToLower(level) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid compliance level: %s (must be one of: %s)",
		level, strings.Join(validLevels, ", "))
}

// ValidateMixinName checks if the provided mixin name is valid
func ValidateMixinName(name string) error {
	if name == "" {
		return fmt.Errorf("mixin name cannot be empty")
	}

	// Check for valid characters (lowercase letters, numbers, and hyphens)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("mixin name must contain only lowercase letters, numbers, and hyphens")
		}
	}

	// Check that it doesn't start or end with a hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("mixin name cannot start or end with a hyphen")
	}

	return nil
}

// TODO: Add validation functions for other inputs if needed (e.g., AuthorName, ModulePath)
