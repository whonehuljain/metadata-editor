package parser

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
)

type DateParser struct {
	dateFormat string
}

func New(dateFormat string) *DateParser {
	return &DateParser{
		dateFormat: dateFormat,
	}
}

func (dp *DateParser) ParseFromFilename(filename string) (time.Time, error) {
	// Remove file extension
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try to parse using the configured format
	parsedTime, err := time.Parse(dp.dateFormat, nameWithoutExt)
	if err != nil {
		// Try common alternative formats
		formats := []string{
			"20060102_150405",     // YYYYMMDD_HHMMSS
			"2006-01-02_15-04-05", // YYYY-MM-DD_HH-MM-SS
			"20060102",            // YYYYMMDD (date only)
			"2006-01-02",          // YYYY-MM-DD (date only)
			"IMG_20060102_150405", // IMG_YYYYMMDD_HHMMSS
			"VID_20060102_150405", // VID_YYYYMMDD_HHMMSS
		}

		for _, format := range formats {
			// Try exact match
			if parsedTime, err = time.Parse(format, nameWithoutExt); err == nil {
				return parsedTime, nil
			}

			// Try with common prefixes removed
			cleanName := dp.removeCommonPrefixes(nameWithoutExt)
			if parsedTime, err = time.Parse(format, cleanName); err == nil {
				return parsedTime, nil
			}
		}

		return time.Time{}, errors.New("could not parse date from filename: " + filename)
	}

	return parsedTime, nil
}

func (dp *DateParser) removeCommonPrefixes(filename string) string {
	prefixes := []string{"IMG_", "VID_", "PHOTO_", "PIC_", "image_", "video_"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToUpper(filename), strings.ToUpper(prefix)) {
			return filename[len(prefix):]
		}
	}
	return filename
}

func (dp *DateParser) IsValidDateFormat(filename string) bool {
	_, err := dp.ParseFromFilename(filename)
	return err == nil
}
