package config

import (
	"errors"
	"os"
	"time"

	"image-metadata-editor/pkg/logger"
)

type Config struct {
	FolderPath     string
	DateString     string
	TimeString     string
	AutoMode       bool
	SequentialMode bool
	StartTime      string
	DateFormat     string
	DryRun         bool
	Logger         *logger.Logger
}

func (c *Config) Validate() error {
	// Check if folder exists
	if _, err := os.Stat(c.FolderPath); os.IsNotExist(err) {
		return errors.New("folder does not exist")
	}

	// Count active modes
	modeCount := 0
	if c.AutoMode {
		modeCount++
	}
	if c.SequentialMode {
		modeCount++
	}
	if c.DateString != "" && !c.SequentialMode && !c.AutoMode {
		modeCount++
	}

	if modeCount != 1 {
		return errors.New("exactly one mode must be specified: -auto, -sequential, or manual date/time")
	}

	// Validate sequential mode requirements
	if c.SequentialMode {
		if c.DateString == "" {
			return errors.New("sequential mode requires -date parameter")
		}
		if c.StartTime == "" {
			return errors.New("sequential mode requires -start-time parameter")
		}

		// Validate date format
		if _, err := time.Parse("2006-01-02", c.DateString); err != nil {
			return errors.New("invalid date format, use YYYY-MM-DD")
		}

		// Validate start time format
		if _, err := time.Parse("15:04:05", c.StartTime); err != nil {
			return errors.New("invalid start-time format, use HH:MM:SS")
		}

		return nil
	}

	// Validate other modes
	if !c.AutoMode && c.DateString != "" {
		if _, err := time.Parse("2006-01-02", c.DateString); err != nil {
			return errors.New("invalid date format, use YYYY-MM-DD")
		}
	}

	if !c.AutoMode && c.TimeString != "" {
		if _, err := time.Parse("15:04:05", c.TimeString); err != nil {
			return errors.New("invalid time format, use HH:MM:SS")
		}
	}

	if !c.AutoMode && !c.SequentialMode && c.DateString == "" {
		return errors.New("date is required when not using auto or sequential mode")
	}

	return nil
}

func (c *Config) GetTargetDateTime() (time.Time, error) {
	if c.AutoMode || c.SequentialMode {
		return time.Time{}, nil // Will be determined per file or folder
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

func (c *Config) GetSequentialStartTime() (time.Time, error) {
	if !c.SequentialMode {
		return time.Time{}, errors.New("not in sequential mode")
	}

	dateTimeStr := c.DateString + " " + c.StartTime
	return time.Parse("2006-01-02 15:04:05", dateTimeStr)
}
