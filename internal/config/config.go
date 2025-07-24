package config

import (
	"errors"
	"os"
	"time"

	"image-metadata-editor/pkg/logger"
)

type Config struct {
	FolderPath string
	DateString string
	TimeString string
	AutoMode   bool
	DateFormat string
	DryRun     bool
	Logger     *logger.Logger
}

func (c *Config) Validate() error {
	// Check if folder exists
	if _, err := os.Stat(c.FolderPath); os.IsNotExist(err) {
		return errors.New("folder does not exist")
	}

	// Validate date format if provided
	if !c.AutoMode && c.DateString != "" {
		if _, err := time.Parse("2006-01-02", c.DateString); err != nil {
			return errors.New("invalid date format, use YYYY-MM-DD")
		}
	}

	// Validate time format if provided
	if !c.AutoMode && c.TimeString != "" {
		if _, err := time.Parse("15:04:05", c.TimeString); err != nil {
			return errors.New("invalid time format, use HH:MM:SS")
		}
	}

	// If not auto mode, at least date should be provided
	if !c.AutoMode && c.DateString == "" {
		return errors.New("date is required when not using auto mode")
	}

	return nil
}

func (c *Config) GetTargetDateTime() (time.Time, error) {
	if c.AutoMode {
		return time.Time{}, nil // Will be determined per file
	}

	dateStr := c.DateString
	if c.TimeString != "" {
		dateStr += " " + c.TimeString
		return time.Parse("2006-01-02 15:04:05", dateStr)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
