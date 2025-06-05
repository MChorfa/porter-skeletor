package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnterpriseFeatureToggleCombinations(t *testing.T) {
	tests := []struct {
		name                      string
		enableSecurity            bool
		enableCompliance          bool
		enableAuth                bool
		enableObservability       bool
		securityFeatures          string
		complianceFrameworks      string
		authFeatures              string
		observabilityFeatures     string
		expectedSecurityEnabled   bool
		expectedComplianceEnabled bool
		expectedAuthEnabled       bool
		expectedObsEnabled        bool
	}{
		{
			name:                      "No enterprise features enabled",
			enableSecurity:            false,
			enableCompliance:          false,
			enableAuth:                false,
			enableObservability:       false,
			securityFeatures:          "",
			complianceFrameworks:      "",
			authFeatures:              "",
			observabilityFeatures:     "",
			expectedSecurityEnabled:   false,
			expectedComplianceEnabled: false,
			expectedAuthEnabled:       false,
			expectedObsEnabled:        false,
		},
		{
			name:                      "Only security features enabled",
			enableSecurity:            true,
			enableCompliance:          false,
			enableAuth:                false,
			enableObservability:       false,
			securityFeatures:          "input_validation,rate_limiting",
			complianceFrameworks:      "",
			authFeatures:              "",
			observabilityFeatures:     "",
			expectedSecurityEnabled:   true,
			expectedComplianceEnabled: false,
			expectedAuthEnabled:       false,
			expectedObsEnabled:        false,
		},
		{
			name:                      "All enterprise features enabled",
			enableSecurity:            true,
			enableCompliance:          true,
			enableAuth:                true,
			enableObservability:       true,
			securityFeatures:          "input_validation,rate_limiting,secure_headers",
			complianceFrameworks:      "soc2,gdpr",
			authFeatures:              "rbac,ldap,sso",
			observabilityFeatures:     "apm,opentelemetry,audit_logging",
			expectedSecurityEnabled:   true,
			expectedComplianceEnabled: true,
			expectedAuthEnabled:       true,
			expectedObsEnabled:        true,
		},
		{
			name:                      "Mixed feature combination",
			enableSecurity:            true,
			enableCompliance:          false,
			enableAuth:                true,
			enableObservability:       false,
			securityFeatures:          "input_validation,policy_enforcement",
			complianceFrameworks:      "",
			authFeatures:              "rbac,vault",
			observabilityFeatures:     "",
			expectedSecurityEnabled:   true,
			expectedComplianceEnabled: false,
			expectedAuthEnabled:       true,
			expectedObsEnabled:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock config
			config := &TemplateConfig{
				Name: "test-mixin",
				Variables: map[string]Variable{
					"MixinName": {Default: "test-mixin"},
					"Author":    {Default: "Test Author"},
				},
			}

			// Build template data with enterprise features
			data, err := buildTemplateDataWithFeatures(
				config,
				"test-mixin",
				"Test Author",
				"github.com/test/test-mixin",
				"./test-output",
				"basic",
				true, // non-interactive
				[]string{},
				tt.enableSecurity,
				tt.enableCompliance,
				tt.enableAuth,
				tt.enableObservability,
				tt.securityFeatures,
				tt.complianceFrameworks,
				tt.authFeatures,
				tt.observabilityFeatures,
			)

			require.NoError(t, err)

			// Verify enterprise feature flags are set correctly
			assert.Equal(t, tt.expectedSecurityEnabled, data["EnableSecurity"])
			assert.Equal(t, tt.expectedComplianceEnabled, data["EnableCompliance"])
			assert.Equal(t, tt.expectedAuthEnabled, data["EnableAuth"])
			assert.Equal(t, tt.expectedObsEnabled, data["EnableObservability"])

			// Verify feature lists are set correctly
			assert.Equal(t, tt.securityFeatures, data["SecurityFeatures"])
			assert.Equal(t, tt.complianceFrameworks, data["ComplianceFrameworks"])
			assert.Equal(t, tt.authFeatures, data["AuthFeatures"])
			assert.Equal(t, tt.observabilityFeatures, data["ObservabilityFeatures"])

			// Test template functions work with this data
			if tt.enableSecurity && tt.securityFeatures != "" {
				assert.True(t, featureEnabled(data, "security", "input_validation"))
			}
			if tt.enableCompliance && tt.complianceFrameworks != "" {
				assert.True(t, featureEnabled(data, "compliance", "soc2"))
			}
			if tt.enableAuth && tt.authFeatures != "" {
				assert.True(t, featureEnabled(data, "auth", "rbac"))
			}
			if tt.enableObservability && tt.observabilityFeatures != "" {
				assert.True(t, featureEnabled(data, "observability", "apm"))
			}
		})
	}
}

func TestConditionalTemplateRendering(t *testing.T) {
	tests := []struct {
		name           string
		templateData   map[string]interface{}
		category       string
		feature        string
		expectedResult bool
	}{
		{
			name: "Security feature enabled in template data",
			templateData: map[string]interface{}{
				"EnableSecurity":   true,
				"SecurityFeatures": "input_validation,rate_limiting",
			},
			category:       "security",
			feature:        "input_validation",
			expectedResult: true,
		},
		{
			name: "Security feature disabled in template data",
			templateData: map[string]interface{}{
				"EnableSecurity":   false,
				"SecurityFeatures": "input_validation,rate_limiting",
			},
			category:       "security",
			feature:        "input_validation",
			expectedResult: false,
		},
		{
			name: "Feature not in list",
			templateData: map[string]interface{}{
				"EnableSecurity":   true,
				"SecurityFeatures": "rate_limiting,secure_headers",
			},
			category:       "security",
			feature:        "input_validation",
			expectedResult: false,
		},
		{
			name: "Compliance feature enabled",
			templateData: map[string]interface{}{
				"EnableCompliance":     true,
				"ComplianceFrameworks": "soc2,gdpr,hipaa",
			},
			category:       "compliance",
			feature:        "gdpr",
			expectedResult: true,
		},
		{
			name: "Auth feature enabled",
			templateData: map[string]interface{}{
				"EnableAuth":   true,
				"AuthFeatures": "rbac,ldap,sso,mfa",
			},
			category:       "auth",
			feature:        "sso",
			expectedResult: true,
		},
		{
			name: "Observability feature enabled",
			templateData: map[string]interface{}{
				"EnableObservability":   true,
				"ObservabilityFeatures": "apm,opentelemetry,audit_logging",
			},
			category:       "observability",
			feature:        "opentelemetry",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := featureEnabled(tt.templateData, tt.category, tt.feature)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	t.Run("No enterprise features - backward compatible", func(t *testing.T) {
		config := &TemplateConfig{
			Name: "legacy-mixin",
			Variables: map[string]Variable{
				"MixinName": {Default: "legacy-mixin"},
				"Author":    {Default: "Legacy Author"},
			},
		}

		// Build template data without enterprise features (legacy mode)
		data, err := buildTemplateData(
			config,
			"legacy-mixin",
			"Legacy Author",
			"github.com/legacy/legacy-mixin",
			"./legacy-output",
			"basic",
			true, // non-interactive
			[]string{},
		)

		require.NoError(t, err)

		// Verify legacy fields are still present
		assert.Equal(t, "legacy-mixin", data["MixinName"])
		assert.Equal(t, "Legacy Author", data["Author"])
		assert.Equal(t, "github.com/legacy/legacy-mixin", data["ModulePath"])
		assert.Equal(t, "basic", data["ComplianceLevel"])

		// Verify enterprise features are present but disabled (backward compatibility)
		// In the new version, these flags are always present but default to false
		assert.Equal(t, false, data["EnableSecurity"], "EnableSecurity should be false in legacy mode")
		assert.Equal(t, false, data["EnableCompliance"], "EnableCompliance should be false in legacy mode")
		assert.Equal(t, false, data["EnableAuth"], "EnableAuth should be false in legacy mode")
		assert.Equal(t, false, data["EnableObservability"], "EnableObservability should be false in legacy mode")

		// Verify enterprise feature lists are empty
		assert.Equal(t, "", data["SecurityFeatures"], "SecurityFeatures should be empty in legacy mode")
		assert.Equal(t, "", data["ComplianceFrameworks"], "ComplianceFrameworks should be empty in legacy mode")
		assert.Equal(t, "", data["AuthFeatures"], "AuthFeatures should be empty in legacy mode")
		assert.Equal(t, "", data["ObservabilityFeatures"], "ObservabilityFeatures should be empty in legacy mode")
	})
}
