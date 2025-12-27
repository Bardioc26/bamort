package pdfrender

import (
	"bamort/logger"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// InitializeTemplates copies default templates to target directory if they don't exist
// or if the content has changed (for development updates)
// This should be called once during application startup
func InitializeTemplates(defaultDir, targetDir string) error {
	// Check if default directory exists
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		return fmt.Errorf("default templates directory does not exist: %s", defaultDir)
	}

	// Read default templates
	entries, err := os.ReadDir(defaultDir)
	if err != nil {
		return fmt.Errorf("failed to read default templates: %w", err)
	}

	// Process each template directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		srcDir := filepath.Join(defaultDir, entry.Name())
		dstDir := filepath.Join(targetDir, entry.Name())

		// Copy or update template directory
		updated, err := syncDir(srcDir, dstDir)
		if err != nil {
			return fmt.Errorf("failed to sync template %s: %w", entry.Name(), err)
		}

		if updated {
			logger.Info("Initialized/updated template: %s", entry.Name())
		} else {
			logger.Debug("Template %s is up to date", entry.Name())
		}
	}

	return nil
}

// syncDir synchronizes source directory to destination, copying only changed files
// Returns true if any files were updated
func syncDir(src, dst string) (bool, error) {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false, err
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return false, err
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return false, err
	}

	updated := false
	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively sync subdirectory
			subUpdated, err := syncDir(srcPath, dstPath)
			if err != nil {
				return false, err
			}
			if subUpdated {
				updated = true
			}
		} else {
			// Check if file needs to be copied
			needsCopy, err := fileNeedsUpdate(srcPath, dstPath)
			if err != nil {
				return false, err
			}

			if needsCopy {
				if err := copyFile(srcPath, dstPath); err != nil {
					return false, err
				}
				updated = true
			}
		}
	}

	return updated, nil
}

// fileNeedsUpdate checks if destination file is missing or differs from source
func fileNeedsUpdate(src, dst string) (bool, error) {
	// If destination doesn't exist, needs update
	dstInfo, err := os.Stat(dst)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	// Get source info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false, err
	}

	// Quick check: if sizes differ, files differ
	if srcInfo.Size() != dstInfo.Size() {
		return true, nil
	}

	// Compare file contents
	return filesContentDiffer(src, dst)
}

// filesContentDiffer compares file contents
func filesContentDiffer(file1, file2 string) (bool, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}

	// Files differ if contents don't match
	return string(content1) != string(content2), nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
