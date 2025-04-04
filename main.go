package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	mixinName := flag.String("name", "", "Name of the new mixin")
	authorName := flag.String("author", "", "Author name for the mixin")
	outputDir := flag.String("output", "", "Output directory (defaults to ./{mixinName})")
	flag.Parse()

	if *mixinName == "" || *authorName == "" {
		fmt.Println("Error: mixin name and author name are required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputDir == "" {
		*outputDir = "./" + *mixinName
	}

	// Clean and validate the output directory path
	cleanedOutputDir := filepath.Clean(*outputDir)
	if strings.HasPrefix(cleanedOutputDir, "../") || filepath.IsAbs(cleanedOutputDir) {
		fmt.Printf("Error: output directory must be a relative path within the current directory: %s\n", *outputDir)
		os.Exit(1)
	}
	// Use the cleaned path for the clone operation
	fmt.Println("Cloning skeletor template...")
	// #nosec G204 -- URL is hardcoded, output dir is cleaned, command is git clone
	cmd := exec.Command("git", "clone", "https://github.com/getporter/skeletor.git", cleanedOutputDir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	// Use cleaned path for subsequent operations
	// Remove .git directory
	if err := os.RemoveAll(filepath.Join(cleanedOutputDir, ".git")); err != nil {
		// Log warning but continue, as failing to remove .git might not be critical
		fmt.Fprintf(os.Stderr, "Warning: failed to remove .git directory from %s: %v\n", cleanedOutputDir, err)
	}

	// Replace skeletor with mixinName and YOURNAME with authorName
	replaceInFiles(cleanedOutputDir, "skeletor", *mixinName)
	replaceInFiles(cleanedOutputDir, "YOURNAME", *authorName)

	// Rename directories, handle potential errors
	cmdDirOld := filepath.Join(cleanedOutputDir, "cmd", "skeletor")
	cmdDirNew := filepath.Join(cleanedOutputDir, "cmd", *mixinName)
	if err := os.Rename(cmdDirOld, cmdDirNew); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not rename %s to %s: %v\n", cmdDirOld, cmdDirNew, err)
	}

	pkgDirOld := filepath.Join(cleanedOutputDir, "pkg", "skeletor")
	pkgDirNew := filepath.Join(cleanedOutputDir, "pkg", *mixinName)
	if err := os.Rename(pkgDirOld, pkgDirNew); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not rename %s to %s: %v\n", pkgDirOld, pkgDirNew, err)
	}

	fmt.Printf("Mixin '%s' successfully created in %s\n", *mixinName, cleanedOutputDir)
	fmt.Println("Next steps:")
	fmt.Println("1. cd", cleanedOutputDir)
	fmt.Println("2. Update go.mod with your module path")
	fmt.Println("3. Run 'mage build test' to verify everything works")
}

func replaceInFiles(dir, oldStr, newStr string) {
	// Handle error returned by filepath.Walk
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log the error encountered during walk, but continue if possible?
			// Or return the error to stop the walk? Let's return to stop.
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || strings.Contains(path, ".git") {
			return nil
		}

		// Read file
		// #nosec G304 -- path is derived from filepath.Walk on a controlled directory
		content, err := os.ReadFile(path)
		if err != nil {
			// Log read errors but continue the walk if possible
			fmt.Fprintf(os.Stderr, "Error reading file %q: %v\n", path, err)
			return nil // Continue walk even if one file fails to read
		}

		// Replace content
		newContent := strings.ReplaceAll(string(content), oldStr, newStr)

		// Write back
		writeErr := os.WriteFile(path, []byte(newContent), info.Mode())
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %q: %v\n", path, writeErr)
			// Return the error to stop the walk
			return writeErr
		}
		return nil
	})

	// Log if the walk itself encountered an error
	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "Error during file replacement walk: %v\n", walkErr)
		// Decide if this should cause program exit? For now, just log.
	}
}
