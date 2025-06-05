package main

import (
	"io/fs"
	"testing"

	"github.com/getporter/skeletor/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedFilesystem(t *testing.T) {
	t.Run("Check embedded FS structure", func(t *testing.T) {
		// Test that we can read the embedded filesystem
		entries, err := fs.ReadDir(pkg.MixinTemplateFS, ".")
		require.NoError(t, err)
		
		t.Logf("Root entries: %d", len(entries))
		for _, entry := range entries {
			t.Logf("  - %s (dir: %v)", entry.Name(), entry.IsDir())
		}
		
		// Check if template directory exists
		templateEntries, err := fs.ReadDir(pkg.MixinTemplateFS, "template")
		require.NoError(t, err)
		
		t.Logf("Template entries: %d", len(templateEntries))
		for _, entry := range templateEntries {
			t.Logf("  - template/%s (dir: %v)", entry.Name(), entry.IsDir())
		}
		
		// Verify we have some expected files
		assert.True(t, len(templateEntries) > 0, "Template directory should contain files")
		
		// Check for template.json
		_, err = fs.Stat(pkg.MixinTemplateFS, "template/template.json")
		assert.NoError(t, err, "template.json should exist")
		
		// Check for some template files
		_, err = fs.Stat(pkg.MixinTemplateFS, "template/cmd/mixin/main.go.tmpl")
		assert.NoError(t, err, "main.go.tmpl should exist")
	})
	
	t.Run("Test WalkDir on embedded FS", func(t *testing.T) {
		var fileCount int
		var dirCount int
		
		err := fs.WalkDir(pkg.MixinTemplateFS, "template", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			
			t.Logf("Walking: %s (dir: %v)", path, d.IsDir())
			
			if d.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
			
			return nil
		})
		
		require.NoError(t, err)
		t.Logf("Found %d files and %d directories", fileCount, dirCount)
		
		assert.True(t, fileCount > 0, "Should find some files")
		assert.True(t, dirCount > 0, "Should find some directories")
	})
	
	t.Run("Test calculateDestPath function", func(t *testing.T) {
		tests := []struct {
			name         string
			originalPath string
			templateRoot string
			expected     string
			shouldSkip   bool
		}{
			{
				name:         "template root itself",
				originalPath: "template",
				templateRoot: "template",
				expected:     "",
				shouldSkip:   true,
			},
			{
				name:         "template.json file",
				originalPath: "template/template.json",
				templateRoot: "template",
				expected:     "",
				shouldSkip:   true,
			},
			{
				name:         "regular template file",
				originalPath: "template/cmd/mixin/main.go.tmpl",
				templateRoot: "template",
				expected:     "cmd/mixin/main.go.tmpl",
				shouldSkip:   false,
			},
			{
				name:         "directory",
				originalPath: "template/cmd",
				templateRoot: "template",
				expected:     "cmd",
				shouldSkip:   false,
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				destPath, skip := calculateDestPath(tt.originalPath, tt.templateRoot, []string{})
				assert.Equal(t, tt.expected, destPath)
				assert.Equal(t, tt.shouldSkip, skip)
			})
		}
	})
}
