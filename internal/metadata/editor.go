package metadata

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Editor struct {
	supportedFormats map[string]bool
	hasExifTool      bool
}

func New() *Editor {
	editor := &Editor{
		supportedFormats: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".tiff": true,
			".tif":  true,
		},
	}

	// Check if exiftool is available
	_, err := exec.LookPath("exiftool")
	editor.hasExifTool = err == nil

	if !editor.hasExifTool {
		fmt.Println("Warning: exiftool not found. Install it for better EXIF support:")
		fmt.Println("  macOS: brew install exiftool")
		fmt.Println("  Ubuntu/Debian: sudo apt-get install libimage-exiftool-perl")
	}

	return editor
}

func (e *Editor) Close() {
	// Nothing to close for direct exec approach
}

func (e *Editor) IsSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return e.supportedFormats[ext]
}

func (e *Editor) UpdateMetadata(filePath string, dateTime time.Time, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would update metadata for: %s to %s\n", filePath, dateTime.Format("2006-01-02 15:04:05"))
		return nil
	}

	if !e.hasExifTool {
		fmt.Printf("ExifTool not available, skipping: %s\n", filepath.Base(filePath))
		return nil
	}

	return e.updateWithExifTool(filePath, dateTime)
}

func (e *Editor) updateWithExifTool(filePath string, dateTime time.Time) error {
	fmt.Printf("DEBUG: Processing %s -> Setting to: %s\n",
		filepath.Base(filePath),
		dateTime.Format("2006-01-02 15:04:05"))

	// Format date for exiftool (YYYY:MM:DD HH:MM:SS)
	dateTimeStr := dateTime.Format("2006:01:02 15:04:05")

	// Step 1: First, rebuild the metadata to fix corruption
	// This removes all metadata and rebuilds it cleanly
	cleanCmd := exec.Command("exiftool",
		"-all=",              // Remove all metadata
		"-tagsfromfile", "@", // Copy tags back from the same file
		"-all:all",            // Copy all supported tags
		"-unsafe",             // Allow potentially unsafe operations
		"-icc_profile",        // Preserve color profile
		"-F",                  // Fix any structure issues
		"-overwrite_original", // Don't create backup
		filePath)

	output, err := cleanCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Could not clean metadata for %s: %v - %s\n", filepath.Base(filePath), err, string(output))
		// Continue anyway, sometimes the original dates still work
	} else {
		fmt.Printf("DEBUG: Successfully cleaned metadata for %s\n", filepath.Base(filePath))
	}

	// Step 2: Now set the new dates on the cleaned file
	cmd := exec.Command("exiftool",
		"-F",                           // Fix any remaining structure issues
		"-overwrite_original",          // Don't create backup files
		"-AllDates="+dateTimeStr,       // Set all main date fields
		"-FileModifyDate="+dateTimeStr, // Also set file modification date
		filePath)

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exiftool failed to set dates: %w - Output: %s", err, string(output))
	}

	fmt.Printf("Successfully updated EXIF metadata for: %s -> %s\n",
		filepath.Base(filePath),
		dateTime.Format("2006-01-02 15:04:05"))
	return nil
}

func (e *Editor) updateFileSystemTimestamp(filePath string, dateTime time.Time) error {
	return os.Chtimes(filePath, dateTime, dateTime)
}

func (e *Editor) GetOriginalTimestamp(filePath string) (time.Time, error) {
	// Use exiftool to read the original timestamp
	cmd := exec.Command("exiftool", "-DateTimeOriginal", "-s", "-s", "-s", filePath)
	output, err := cmd.Output()
	if err != nil {
		// Try alternative date fields
		cmd = exec.Command("exiftool", "-DateTime", "-s", "-s", "-s", filePath)
		output, err = cmd.Output()
		if err != nil {
			return time.Time{}, fmt.Errorf("could not read timestamp: %w", err)
		}
	}

	// Parse the timestamp (format: "YYYY:MM:DD HH:MM:SS")
	timestampStr := strings.TrimSpace(string(output))
	if timestampStr == "" {
		return time.Time{}, fmt.Errorf("no timestamp found")
	}

	timestamp, err := time.Parse("2006:01:02 15:04:05", timestampStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse timestamp '%s': %w", timestampStr, err)
	}

	return timestamp, nil
}
