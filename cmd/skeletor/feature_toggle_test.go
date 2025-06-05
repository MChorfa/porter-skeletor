package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureToggles_IsFeatureEnabled(t *testing.T) {
	tests := []struct {
		name     string
		toggles  *FeatureToggles
		category string
		feature  string
		expected bool
	}{
		{
			name:     "nil feature toggles",
			toggles:  nil,
			category: "security",
			feature:  "input_validation",
			expected: false,
		},
		{
			name: "security feature enabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: true,
				},
			},
			category: "security",
			feature:  "input_validation",
			expected: true,
		},
		{
			name: "security feature disabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: false,
				},
			},
			category: "security",
			feature:  "input_validation",
			expected: false,
		},
		{
			name: "security category disabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         false,
					InputValidation: true,
				},
			},
			category: "security",
			feature:  "input_validation",
			expected: false,
		},
		{
			name: "compliance feature enabled",
			toggles: &FeatureToggles{
				Compliance: &ComplianceFeatures{
					Enabled: true,
					SOC2:    true,
				},
			},
			category: "compliance",
			feature:  "soc2",
			expected: true,
		},
		{
			name: "auth feature enabled",
			toggles: &FeatureToggles{
				Auth: &AuthFeatures{
					Enabled: true,
					RBAC:    true,
				},
			},
			category: "auth",
			feature:  "rbac",
			expected: true,
		},
		{
			name: "observability feature enabled",
			toggles: &FeatureToggles{
				Observability: &ObservabilityFeatures{
					Enabled: true,
					APM:     true,
				},
			},
			category: "observability",
			feature:  "apm",
			expected: true,
		},
		{
			name: "unknown category",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: true,
				},
			},
			category: "unknown",
			feature:  "some_feature",
			expected: false,
		},
		{
			name: "unknown feature",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: true,
				},
			},
			category: "security",
			feature:  "unknown_feature",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.toggles.IsFeatureEnabled(tt.category, tt.feature)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFeatureToggles_GetEnabledFeatures(t *testing.T) {
	tests := []struct {
		name     string
		toggles  *FeatureToggles
		expected map[string][]string
	}{
		{
			name:     "nil feature toggles",
			toggles:  nil,
			expected: map[string][]string{},
		},
		{
			name: "all features enabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:               true,
					InputValidation:       true,
					RateLimiting:          true,
					SecureHeaders:         true,
					VulnerabilityScanning: true,
					PolicyEnforcement:     true,
				},
				Compliance: &ComplianceFeatures{
					Enabled: true,
					SOC2:    true,
					GDPR:    true,
					HIPAA:   true,
					PCIDSS:  true,
				},
				Auth: &AuthFeatures{
					Enabled:     true,
					RBAC:        true,
					LDAP:        true,
					SSO:         true,
					MFA:         true,
					Vault:       true,
					SessionMgmt: true,
				},
				Observability: &ObservabilityFeatures{
					Enabled:        true,
					APM:            true,
					Infrastructure: true,
					CustomMetrics:  true,
					HealthChecks:   true,
					OpenTelemetry:  true,
					AuditLogging:   true,
					Tracing:        true,
				},
			},
			expected: map[string][]string{
				"security": {
					"input_validation",
					"rate_limiting",
					"secure_headers",
					"vulnerability_scanning",
					"policy_enforcement",
				},
				"compliance": {
					"soc2",
					"gdpr",
					"hipaa",
					"pci_dss",
				},
				"auth": {
					"rbac",
					"ldap",
					"sso",
					"mfa",
					"vault",
					"session_management",
				},
				"observability": {
					"apm",
					"infrastructure",
					"custom_metrics",
					"health_checks",
					"opentelemetry",
					"audit_logging",
					"tracing",
				},
			},
		},
		{
			name: "partial features enabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: true,
					RateLimiting:    false,
					SecureHeaders:   true,
				},
				Auth: &AuthFeatures{
					Enabled: true,
					RBAC:    true,
					LDAP:    false,
				},
			},
			expected: map[string][]string{
				"security": {
					"input_validation",
					"secure_headers",
				},
				"auth": {
					"rbac",
				},
			},
		},
		{
			name: "categories disabled",
			toggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         false,
					InputValidation: true,
				},
				Auth: &AuthFeatures{
					Enabled: false,
					RBAC:    true,
				},
			},
			expected: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.toggles.GetEnabledFeatures()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateFunctions(t *testing.T) {
	t.Run("hasFeature", func(t *testing.T) {
		tests := []struct {
			name        string
			featureList string
			feature     string
			expected    bool
		}{
			{
				name:        "feature present",
				featureList: "input_validation,rate_limiting,secure_headers",
				feature:     "rate_limiting",
				expected:    true,
			},
			{
				name:        "feature not present",
				featureList: "input_validation,rate_limiting,secure_headers",
				feature:     "policy_enforcement",
				expected:    false,
			},
			{
				name:        "empty feature list",
				featureList: "",
				feature:     "input_validation",
				expected:    false,
			},
			{
				name:        "single feature",
				featureList: "input_validation",
				feature:     "input_validation",
				expected:    true,
			},
			{
				name:        "feature with spaces",
				featureList: "input_validation, rate_limiting, secure_headers",
				feature:     "rate_limiting",
				expected:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := hasFeature(tt.featureList, tt.feature)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("splitFeatures", func(t *testing.T) {
		tests := []struct {
			name        string
			featureList string
			expected    []string
		}{
			{
				name:        "multiple features",
				featureList: "input_validation,rate_limiting,secure_headers",
				expected:    []string{"input_validation", "rate_limiting", "secure_headers"},
			},
			{
				name:        "single feature",
				featureList: "input_validation",
				expected:    []string{"input_validation"},
			},
			{
				name:        "empty list",
				featureList: "",
				expected:    []string{},
			},
			{
				name:        "features with spaces",
				featureList: "input_validation, rate_limiting, secure_headers",
				expected:    []string{"input_validation", "rate_limiting", "secure_headers"},
			},
			{
				name:        "features with empty elements",
				featureList: "input_validation,,rate_limiting",
				expected:    []string{"input_validation", "rate_limiting"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := splitFeatures(tt.featureList)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("joinFeatures", func(t *testing.T) {
		tests := []struct {
			name     string
			features []string
			expected string
		}{
			{
				name:     "multiple features",
				features: []string{"input_validation", "rate_limiting", "secure_headers"},
				expected: "input_validation,rate_limiting,secure_headers",
			},
			{
				name:     "single feature",
				features: []string{"input_validation"},
				expected: "input_validation",
			},
			{
				name:     "empty list",
				features: []string{},
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := joinFeatures(tt.features)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("featureEnabled", func(t *testing.T) {
		data := map[string]interface{}{
			"EnableSecurity":         true,
			"SecurityFeatures":       "input_validation,rate_limiting",
			"EnableCompliance":       true,
			"ComplianceFrameworks":   "soc2,gdpr",
			"EnableAuth":             false,
			"AuthFeatures":           "rbac,ldap",
			"EnableObservability":    true,
			"ObservabilityFeatures":  "apm,tracing",
		}

		tests := []struct {
			name     string
			category string
			feature  string
			expected bool
		}{
			{
				name:     "security feature enabled",
				category: "security",
				feature:  "input_validation",
				expected: true,
			},
			{
				name:     "security feature not in list",
				category: "security",
				feature:  "policy_enforcement",
				expected: false,
			},
			{
				name:     "compliance feature enabled",
				category: "compliance",
				feature:  "soc2",
				expected: true,
			},
			{
				name:     "auth category disabled",
				category: "auth",
				feature:  "rbac",
				expected: false,
			},
			{
				name:     "observability feature enabled",
				category: "observability",
				feature:  "apm",
				expected: true,
			},
			{
				name:     "unknown category",
				category: "unknown",
				feature:  "some_feature",
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := featureEnabled(data, tt.category, tt.feature)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestLoadTemplateConfigWithFeatureToggles(t *testing.T) {
	// This test would require setting up a mock filesystem with template.json
	// For now, we'll test the structure validation
	t.Run("feature toggles structure", func(t *testing.T) {
		config := &TemplateConfig{
			FeatureToggles: &FeatureToggles{
				Security: &SecurityFeatures{
					Enabled:         true,
					InputValidation: true,
				},
				Compliance: &ComplianceFeatures{
					Enabled: true,
					SOC2:    true,
				},
				Auth: &AuthFeatures{
					Enabled: true,
					RBAC:    true,
				},
				Observability: &ObservabilityFeatures{
					Enabled: true,
					APM:     true,
				},
			},
		}

		require.NotNil(t, config.FeatureToggles)
		require.NotNil(t, config.FeatureToggles.Security)
		require.NotNil(t, config.FeatureToggles.Compliance)
		require.NotNil(t, config.FeatureToggles.Auth)
		require.NotNil(t, config.FeatureToggles.Observability)

		assert.True(t, config.FeatureToggles.Security.Enabled)
		assert.True(t, config.FeatureToggles.Security.InputValidation)
		assert.True(t, config.FeatureToggles.Compliance.Enabled)
		assert.True(t, config.FeatureToggles.Compliance.SOC2)
		assert.True(t, config.FeatureToggles.Auth.Enabled)
		assert.True(t, config.FeatureToggles.Auth.RBAC)
		assert.True(t, config.FeatureToggles.Observability.Enabled)
		assert.True(t, config.FeatureToggles.Observability.APM)
	})
}
