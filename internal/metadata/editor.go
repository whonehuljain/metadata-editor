package metadata

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

type Editor struct {
	supportedFormats map[string]bool
}

func New() *Editor {
	return &Editor{
		supportedFormats: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".tiff": true,
			".tif":  true,
		},
	}
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

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".jpg", ".jpeg":
		return e.updateJPEGMetadata(filePath, dateTime)
	case ".png", ".tiff", ".tif":
		return e.updateFileSystemTimestamp(filePath, dateTime)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

func (e *Editor) updateJPEGMetadata(filePath string, dateTime time.Time) error {
	fmt.Printf("DEBUG: Processing %s -> Setting to: %s\n",
		filepath.Base(filePath),
		dateTime.Format("2006-01-02 15:04:05"))

	// Read the JPEG file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse JPEG: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Get or construct EXIF builder
	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		fmt.Printf("DEBUG: Creating new EXIF for %s\n", filepath.Base(filePath))

		// Create new EXIF data
		im, err := exifcommon.NewIfdMappingWithStandard()
		if err != nil {
			return fmt.Errorf("failed to create IFD mapping: %w", err)
		}

		ti := exif.NewTagIndex()
		if err := exif.LoadStandardTags(ti); err != nil {
			return fmt.Errorf("failed to load standard tags: %w", err)
		}

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.EncodeDefaultByteOrder)
	} else {
		fmt.Printf("DEBUG: Using existing EXIF for %s\n", filepath.Base(filePath))
	}

	// Format the datetime string in EXIF format manually (YYYY:MM:DD HH:MM:SS)
	updatedTimestampPhrase := dateTime.Format("2006:01:02 15:04:05")
	fmt.Printf("DEBUG: Setting datetime to: %s\n", updatedTimestampPhrase)

	// Get IFD0 for DateTime
	ifdPath := "IFD0"
	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	if err != nil {
		return fmt.Errorf("failed to get/create IFD0: %w", err)
	}

	// Set DateTime in IFD0
	if err := ifdIb.SetStandardWithName("DateTime", updatedTimestampPhrase); err != nil {
		fmt.Printf("Warning: Could not set DateTime: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Successfully set DateTime\n")
	}

	// Get EXIF IFD for DateTimeOriginal and DateTimeDigitized
	exifIfdPath := "IFD0/Exif"
	exifIfdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, exifIfdPath)
	if err != nil {
		fmt.Printf("Warning: Could not get/create EXIF IFD: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Got EXIF IFD successfully\n")

		// Set DateTimeOriginal (this is what Google Photos uses primarily)
		if err := exifIfdIb.SetStandardWithName("DateTimeOriginal", updatedTimestampPhrase); err != nil {
			fmt.Printf("Warning: Could not set DateTimeOriginal: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Successfully set DateTimeOriginal\n")
		}

		// Set DateTimeDigitized
		if err := exifIfdIb.SetStandardWithName("DateTimeDigitized", updatedTimestampPhrase); err != nil {
			fmt.Printf("Warning: Could not set DateTimeDigitized: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Successfully set DateTimeDigitized\n")
		}
	}

	// Update the EXIF segment
	err = sl.SetExif(rootIb)
	if err != nil {
		return fmt.Errorf("failed to set EXIF in JPEG: %w", err)
	}

	// Write the updated JPEG
	b := new(bytes.Buffer)
	if err := sl.Write(b); err != nil {
		return fmt.Errorf("failed to write JPEG: %w", err)
	}

	// Save to file
	if err := ioutil.WriteFile(filePath, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully updated EXIF metadata for: %s\n", filepath.Base(filePath))
	return nil
}

func (e *Editor) updateFileSystemTimestamp(filePath string, dateTime time.Time) error {
	return os.Chtimes(filePath, dateTime, dateTime)
}
