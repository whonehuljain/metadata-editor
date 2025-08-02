package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"image-metadata-editor/internal/config"
	"image-metadata-editor/internal/processor"
	"image-metadata-editor/pkg/logger"
)

func main() {
	var (
		folderPath = flag.String("folder", "", "Path to the folder containing images")
		dateStr    = flag.String("date", "", "Date in YYYY-MM-DD format (optional if using filename parsing)")
		timeStr    = flag.String("time", "", "Time in HH:MM:SS format (optional)")
		autoMode   = flag.Bool("auto", false, "Automatically extract date/time from filename")
		dateFormat = flag.String("format", "20060102_150405", "Date format for filename parsing (default: YYYYMMDD_HHMMSS)")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
		dryRun     = flag.Bool("dry-run", false, "Show what would be changed without making actual changes")
	)
	flag.Parse()

	if *folderPath == "" {
		fmt.Println("Usage: image-metadata-editor -folder /path/to/images [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize logger
	logger := logger.New(*verbose)

	// Create configuration
	cfg := &config.Config{
		FolderPath: *folderPath,
		DateString: *dateStr,
		TimeString: *timeStr,
		AutoMode:   *autoMode,
		DateFormat: *dateFormat,
		DryRun:     *dryRun,
		Logger:     logger,
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Process images
	processor := processor.New(cfg)
	if err := processor.ProcessFolder(); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	fmt.Println("Image metadata processing completed successfully!")
}
