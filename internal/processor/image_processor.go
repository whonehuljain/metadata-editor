package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"image-metadata-editor/internal/config"
	"image-metadata-editor/internal/metadata"
	"image-metadata-editor/internal/parser"
	"image-metadata-editor/internal/sequential"
)

type ImageProcessor struct {
	config     *config.Config
	editor     *metadata.Editor
	parser     *parser.DateParser
	calculator *sequential.Calculator
	stats      ProcessingStats
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

func (ip *ImageProcessor) Close() {
	if ip.editor != nil {
		ip.editor.Close()
	}
}

func (ip *ImageProcessor) ProcessFolder() error {
	ip.config.Logger.Info(fmt.Sprintf("Processing folder: %s", ip.config.FolderPath))

	if ip.config.DryRun {
		ip.config.Logger.Info("DRY RUN MODE - No files will be modified")
	}

	if ip.config.SequentialMode {
		return ip.processSequentialMode()
	}

	err := filepath.Walk(ip.config.FolderPath, ip.processFile)
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	ip.printStats()
	return nil
}

func (ip *ImageProcessor) processSequentialMode() error {
	ip.config.Logger.Info("Running in SEQUENTIAL mode - preserving time differences")

	// Get the base date/time for sequential mode
	baseDateTime, err := ip.config.GetSequentialStartTime()
	if err != nil {
		return fmt.Errorf("failed to get sequential start time: %w", err)
	}

	ip.calculator = sequential.NewCalculator(baseDateTime)

	// First pass: collect all photos and their original timestamps
	err = filepath.Walk(ip.config.FolderPath, ip.collectPhotoTimestamps)
	if err != nil {
		return fmt.Errorf("error collecting timestamps: %w", err)
	}

	if ip.calculator.GetPhotoCount() == 0 {
		return fmt.Errorf("no supported image files found in folder")
	}

	ip.config.Logger.Info(fmt.Sprintf("Found %d photos, calculating sequential times...", ip.calculator.GetPhotoCount()))

	// Calculate new times based on original time differences
	ip.calculator.CalculateNewTimes()

	// Second pass: update metadata with calculated times
	err = filepath.Walk(ip.config.FolderPath, ip.processSequentialFile)
	if err != nil {
		return fmt.Errorf("error processing sequential files: %w", err)
	}

	ip.printStats()
	return nil
}

func (ip *ImageProcessor) collectPhotoTimestamps(path string, info os.FileInfo, err error) error {
	if err != nil {
		return nil // Skip errors during collection
	}

	if info.IsDir() {
		return nil
	}

	ip.stats.TotalFiles++

	if !ip.editor.IsSupported(info.Name()) {
		ip.stats.SkippedFiles++
		return nil
	}

	// Get original timestamp from the file
	originalTime, err := ip.editor.GetOriginalTimestamp(path)
	if err != nil {
		ip.config.Logger.Error(fmt.Sprintf("Could not read original timestamp from %s: %v", info.Name(), err))
		// Use file modification time as fallback
		originalTime = info.ModTime()
		ip.config.Logger.Debug(fmt.Sprintf("Using file modification time for %s: %s", info.Name(), originalTime.Format("2006-01-02 15:04:05")))
	}

	ip.calculator.AddPhoto(path, originalTime)
	ip.config.Logger.Debug(fmt.Sprintf("Collected timestamp for %s: %s", info.Name(), originalTime.Format("2006-01-02 15:04:05")))

	return nil
}

func (ip *ImageProcessor) processSequentialFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	if !ip.editor.IsSupported(info.Name()) {
		return nil
	}

	// Get the calculated new time for this file
	newTime, found := ip.calculator.GetPhotoNewTime(path)
	if !found {
		ip.config.Logger.Error(fmt.Sprintf("No calculated time found for %s", info.Name()))
		ip.stats.ErrorFiles++
		return nil
	}

	// Update metadata
	if err := ip.editor.UpdateMetadata(path, newTime, ip.config.DryRun); err != nil {
		ip.config.Logger.Error(fmt.Sprintf("Failed to update metadata for %s: %v", path, err))
		ip.stats.ErrorFiles++
		return nil
	}

	ip.config.Logger.Info(fmt.Sprintf("Successfully processed: %s -> %s", info.Name(), newTime.Format("2006-01-02 15:04:05")))
	ip.stats.ProcessedFiles++

	return nil
}

func (ip *ImageProcessor) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		ip.config.Logger.Error(fmt.Sprintf("Error accessing file %s: %v", path, err))
		return nil
	}

	if info.IsDir() {
		return nil
	}

	ip.stats.TotalFiles++

	if !ip.editor.IsSupported(info.Name()) {
		ip.config.Logger.Debug(fmt.Sprintf("Skipping unsupported file: %s", path))
		ip.stats.SkippedFiles++
		return nil
	}

	var targetDateTime time.Time
	var err2 error

	if ip.config.AutoMode {
		targetDateTime, err2 = ip.parser.ParseFromFilename(info.Name())
		if err2 != nil {
			ip.config.Logger.Error(fmt.Sprintf("Could not parse date from filename %s: %v", info.Name(), err2))
			ip.stats.ErrorFiles++
			return nil
		}
		ip.config.Logger.Debug(fmt.Sprintf("Extracted date from filename %s: %s", info.Name(), targetDateTime.Format("2006-01-02 15:04:05")))
	} else {
		targetDateTime, err2 = ip.config.GetTargetDateTime()
		if err2 != nil {
			ip.config.Logger.Error(fmt.Sprintf("Invalid target date/time: %v", err2))
			ip.stats.ErrorFiles++
			return nil
		}
	}

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
