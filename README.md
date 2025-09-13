# Metadata Editor

Have you ever clicked pictures from a digital camera, just to realize later that the date \& time were not set correctly?

Or have you moved your photos countless times across different devices and storage systems, which has corrupted the EXIF data, and now all your photos show today's date instead of when they were actually taken?

With hundreds or thousands of photos, manually changing the date and time of each photo is a nightmare that can take days or weeks of tedious work.

**This tool will save you from that hell.**

The Image Metadata Editor is a powerful command-line tool that intelligently processes your photo collections, preserving the relative timing between images while correcting their metadata timestamps. Whether your photos have corrupted EXIF data, incorrect camera dates, or missing timestamps entirely, this tool can fix them all in minutes, not days.

## 🎯 What This Tool Does

- **Fixes corrupted EXIF data** from older digital cameras and file transfers
- **Preserves time relationships** between photos while updating dates
- **Bulk processes** entire folders of images efficiently
- **Extracts dates from filenames** automatically (e.g., `IMG_20241031_143052.jpg`)
- **Handles time corrections** for camera clock drift (±seconds to hours)
- **Prepares photos for Google Photos** and other cloud services


## 📋 Prerequisites

Before using this tool, you need to install the following dependencies:

### 1. ExifTool

ExifTool is required for reading and writing image metadata.

**macOS:**

```bash
brew install exiftool
```

**Ubuntu/Debian:**

```bash
sudo apt-get update
sudo apt-get install libimage-exiftool-perl
```

**Windows:**
Download from [exiftool.org](https://exiftool.org/) and add to your PATH.

### 2. Go Programming Language

This tool is written in Go, so you need Go installed to build it.

**Installation:**

- Download from [golang.org](https://golang.org/dl/)
- Follow the installation instructions for your operating system
- Verify installation: `go version`


## 🚀 Setup and Installation

1. **Clone or download this repository:**
```bash
git clone https://github.com/whonehuljain/metadata-editor.git
cd metadata-editor
```

2. **Build the tool:**
```bash
go mod tidy
go build -o image-metadata-editor ./cmd
```

3. **Verify installation:**
```bash
./image-metadata-editor --help
```


## 📖 Usage Examples

### Basic Examples

**Test with dry-run (recommended first step):**

```bash
./image-metadata-editor -folder /path/to/photos -sequential -date 2024-10-31 -start-time 15:30:00 -dry-run -verbose
```

**Process photos from a vacation folder:**

```bash
./image-metadata-editor -folder ./vacation-photos -sequential -date 2024-07-15 -start-time 09:00:00 -verbose
```

**Extract dates from filenames automatically:**

```bash
./image-metadata-editor -folder ./camera-dump -auto -verbose
```

**Correct camera time that was 2 minutes fast:**

```bash
./image-metadata-editor -folder ./photos -sequential -date 2024-10-31 -start-time 14:00:00 -offset -2m -verbose
```


### Advanced Examples

**Handle different filename formats:**

```bash
./image-metadata-editor -folder ./photos -auto -format "2006-01-02_15-04-05" -verbose
```

**Set same timestamp for all photos:**

```bash
./image-metadata-editor -folder ./group-photos -date 2024-10-31 -time 18:30:00 -verbose
```


## 🎛️ Processing Modes

The tool offers three distinct modes to handle different scenarios:

### 1. **Auto Mode** (`-auto`)

**What it does:** Automatically extracts date and time from image filenames.

**Best for:**

- Photos from smartphones with timestamp filenames
- Screenshots with date stamps
- Camera files with embedded dates

**Supported filename patterns:**

- `IMG_20241031_143052.jpg` → October 31, 2024 at 2:30:52 PM
- `20241031_143052.jpg` → October 31, 2024 at 2:30:52 PM
- `2024-10-31_14-30-52.jpg` → October 31, 2024 at 2:30:52 PM

**Example:**

```bash
./image-metadata-editor -folder ./smartphone-photos -auto -verbose
```


### 2. **Sequential Mode** (`-sequential`)

**What it does:** Sets all photos to a specific date while preserving the original time differences between them.

**Best for:**

- Camera photos where the date was wrong but the timing sequence is correct
- Vacation photos that need to be moved to the correct date
- Event photos with systematic time errors

**How it works:**

1. Reads original timestamps from all photos
2. Calculates time differences between photos
3. Applies the same differences to your new base date/time

**Example scenario:**
Your camera photos were taken at:

- Photo 1: 5:30:15 AM
- Photo 2: 5:32:30 AM
- Photo 3: 5:35:45 AM

With base time 3:00:00 PM, they become:

- Photo 1: 3:00:00 PM
- Photo 2: 3:02:15 PM
- Photo 3: 3:05:30 PM

**Example:**

```bash
./image-metadata-editor -folder ./camera-photos -sequential -date 2024-10-31 -start-time 15:00:00 -verbose
```


### 3. **Manual Mode** (default)

**What it does:** Sets the exact same date and time for all photos in the folder.

**Best for:**

- Group photos taken at the same event
- Scanned photos from a specific date
- Photos where exact timing doesn't matter

**Example:**

```bash
./image-metadata-editor -folder ./group-photos -date 2024-10-31 -time 18:00:00 -verbose
```


## 🔧 Advanced Features

### Time Offset Correction (`-offset`)

Compensates for camera clock errors or timezone differences.

**Supported formats:**

- `+30s` or `-30s` (seconds)
- `+2m30s` or `-2m30s` (minutes and seconds)
- `+1h30m` or `-1h30m` (hours and minutes)
- `+1h5m30s` or `-1h5m30s` (full format)

**Examples:**

```bash
# Camera was 30 seconds slow
./image-metadata-editor -folder ./photos -sequential -date 2024-10-31 -start-time 15:00:00 -offset +30s

# Camera was 2 hours fast (wrong timezone)
./image-metadata-editor -folder ./photos -auto -offset -2h
```


### Custom Date Formats (`-format`)

Specify custom patterns for filename date extraction.

**Default format:** `20060102_150405` (YYYYMMDD_HHMMSS)

**Custom examples:**

```bash
# For format: 2024-10-31_14-30-52
./image-metadata-editor -folder ./photos -auto -format "2006-01-02_15-04-05"

# For format: 31Oct2024_143052
./image-metadata-editor -folder ./photos -auto -format "02Jan2006_150405"
```


### Dry Run Mode (`-dry-run`)

Preview changes without actually modifying files.

```bash
./image-metadata-editor -folder ./photos -sequential -date 2024-10-31 -start-time 15:00:00 -dry-run -verbose
```


### Verbose Logging (`-verbose`)

Get detailed information about the processing steps.

## 📁 Supported File Formats

- **JPEG/JPG**: Full EXIF metadata support with corruption repair
- **PNG**: Basic metadata support
- **TIFF**: Complete metadata handling


## ⚠️ Important Notes

1. **Backup your photos** before running the tool (though it doesn't delete original data)
2. **Test with `-dry-run`** first to verify the results
3. **Use `-verbose`** to understand what the tool is doing
4. The tool **fixes corrupted EXIF data** automatically during processing
5. **Google Photos compatibility**: Updated metadata works perfectly with cloud photo services

## 🆘 Troubleshooting

**ExifTool not found:**

```bash
# Verify ExifTool installation
exiftool -ver

# If not installed, install it:
brew install exiftool  # macOS
sudo apt-get install libimage-exiftool-perl  # Ubuntu
```

**No dates extracted from filenames:**

- Check your filename format with `-verbose`
- Use custom format with `-format` parameter
- Consider using sequential mode instead

**Corrupted EXIF data:**

- The tool automatically fixes most corruption issues
- Use `-verbose` to see repair progress
- Some very damaged files may need manual inspection


## 🎯 Real-World Scenarios

### Scenario 1: Digital Camera Import

```bash
# Camera date was wrong, but photo sequence timing is correct
./image-metadata-editor -folder ./DCIM/Camera -sequential -date 2024-10-31 -start-time 09:00:00 -verbose
```


### Scenario 2: Phone Backup Processing

```bash
# Smartphone photos with timestamps in filenames
./image-metadata-editor -folder ./phone-backup -auto -verbose
```


### Scenario 3: Scanned Photo Dating

```bash
# All scanned photos from same event
./image-metadata-editor -folder ./scanned-1995-birthday -date 1995-06-15 -time 16:00:00 -verbose
```


### Scenario 4: Camera Clock Correction

```bash
# Camera was 1 hour 30 minutes fast
./image-metadata-editor -folder ./photos -auto -offset -1h30m -verbose
```


***

**Made with ❤️ to save you from metadata hell**
