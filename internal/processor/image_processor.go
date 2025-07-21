package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"image-metadata-editor/internal/config"
	"image-metadata-editor/internal/metadata"
	"image-metadata-editor/internal/parser"
)

type ImageProcessor struct {
	config *config.Config
	editor *metadata.Editor
	parser *parser.DateParser
	stats  ProcessingStats
}

type ProcessingStats struct {
	TotalFiles     int
	ProcessedFiles int
	SkippedFiles   int
	ErrorFiles     int
}

func New(cfg *config.Config) *ImageProcessor {
	return &ImageProcessor{
		config: cfg,
		editor: metadata.New(),
		parser: parser.New(cfg.DateFormat),
		stats:  ProcessingStats{},
	}
}

func (ip *ImageProcessor) ProcessFolder() error {
	ip.config.Logger.Info(fmt.Sprintf("Processing folder: %s", ip.config.FolderPath))

	if ip.config.DryRun {
		ip.config.Logger.Info("DRY RUN MODE - No files will be modified")
	}

	err := filepath.Walk(ip.config.FolderPath, ip.processFile)
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	ip.printStats()
	return nil
}

func (ip *ImageProcessor) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		ip.config.Logger.Error(fmt.Sprintf("Error accessing file %s: %v", path, err))
		return nil // Continue processing other files
	}

	// Skip directories
	if info.IsDir() {
		return nil
	}

	ip.stats.TotalFiles++

	// Check if file is supported
	if !ip.editor.IsSupported(info.Name()) {
		ip.config.Logger.Debug(fmt.Sprintf("Skipping unsupported file: %s", path))
		ip.stats.SkippedFiles++
		return nil
	}

	// Determine target date/time
	var targetDateTime time.Time
	var err2 error

	if ip.config.AutoMode {
		// Extract date/time from filename
		targetDateTime, err2 = ip.parser.ParseFromFilename(info.Name())
		if err2 != nil {
			ip.config.Logger.Error(fmt.Sprintf("Could not parse date from filename %s: %v", info.Name(), err2))
			ip.stats.ErrorFiles++
			return nil
		}
		ip.config.Logger.Debug(fmt.Sprintf("Extracted date from filename %s: %s", info.Name(), targetDateTime.Format("2006-01-02 15:04:05")))
	} else {
		// Use provided date/time
		targetDateTime, err2 = ip.config.GetTargetDateTime()
		if err2 != nil {
			ip.config.Logger.Error(fmt.Sprintf("Invalid target date/time: %v", err2))
			ip.stats.ErrorFiles++
			return nil
		}
	}

	// Update metadata
	if err := ip.editor.UpdateMetadata(path, targetDateTime, ip.config.DryRun); err != nil {
		ip.config.Logger.Error(fmt.Sprintf("Failed to update metadata for %s: %v", path, err))
		ip.stats.ErrorFiles++
		return nil
	}

	ip.config.Logger.Info(fmt.Sprintf("Successfully processed: %s -> %s", info.Name(), targetDateTime.Format("2006-01-02 15:04:05")))
	ip.stats.ProcessedFiles++

	return nil
}

func (ip *ImageProcessor) printStats() {
	fmt.Println("\n=== Processing Summary ===")
	fmt.Printf("Total files found: %d\n", ip.stats.TotalFiles)
	fmt.Printf("Successfully processed: %d\n", ip.stats.ProcessedFiles)
	fmt.Printf("Skipped (unsupported): %d\n", ip.stats.SkippedFiles)
	fmt.Printf("Errors: %d\n", ip.stats.ErrorFiles)
}
