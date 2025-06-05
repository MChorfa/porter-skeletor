package main

import (
	"io/fs"
	"testing"

	"github.com/getporter/skeletor/pkg"
	"github.com/stretchr/testify/require"
)

func TestCreateMixinDebug(t *testing.T) {
	t.Run("Debug createMixin dry run", func(t *testing.T) {
		// Load the template config
		config, err := LoadTemplateConfig(pkg.MixinTemplateFS, "template")
		require.NoError(t, err)

		// Create test data
		data := map[string]interface{}{
			"MixinName":       "debug-test",
			"AuthorName":      "Debug Author",
			"ModulePath":      "github.com/test/debug-test",
			"OutputDir":       "./debug-test",
			"ComplianceLevel": "basic",
		}

		// Test the createMixin function with dry run
		t.Logf("Testing createMixin with dry run...")

		// Actually call createMixin with dry run
		err = createMixin(data, pkg.MixinTemplateFS, "template", "./debug-test", config, true)
		require.NoError(t, err)

		// Also test the walking behavior separately for debugging
		t.Logf("Testing WalkDir separately...")
		err = fs.WalkDir(pkg.MixinTemplateFS, "template", func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				t.Logf("Walk error at %s: %v", path, walkErr)
				return walkErr
			}

			// Calculate destination path and check if the file/dir should be skipped
			destRelPath, skip := calculateDestPath(path, "template", config.Ignore)
			t.Logf("Path: %s -> destRelPath: %s, skip: %v", path, destRelPath, skip)

			if skip {
				if d.IsDir() {
					t.Logf("  Skipping directory: %s", path)
					return fs.SkipDir
				}
				t.Logf("  Skipping file: %s", path)
				return nil
			}

			// Determine the actual source path and file info, handling conditional logic
			sourcePath, info, skip, err := determineSourcePath(pkg.MixinTemplateFS, path, destRelPath, "template", config.ConditionalPaths, data)
			if err != nil {
				t.Logf("  Error in determineSourcePath for %s: %v", path, err)
				return err
			}
			if skip {
				t.Logf("  Skipped by conditional logic: %s", path)
				if info != nil && info.IsDir() {
					return fs.SkipDir
				}
				return nil
			}

			// Process the final destination path using template data
			finalDestPath, err := processDestPath(destRelPath, "./debug-test", data)
			if err != nil {
				t.Logf("  Error in processDestPath for %s: %v", destRelPath, err)
				return err
			}
			if finalDestPath == "" {
				t.Logf("  Empty final dest path for: %s", destRelPath)
				if info.IsDir() {
					return fs.SkipDir
				}
				return nil
			}

			t.Logf("  Would process: %s -> %s (source: %s)", path, finalDestPath, sourcePath)
			return nil
		})

		require.NoError(t, err)
	})
}
