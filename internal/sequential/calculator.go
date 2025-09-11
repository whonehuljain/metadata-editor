package sequential

import (
	"fmt"
	"sort"
	"time"
)

type PhotoTime struct {
	FilePath     string
	OriginalTime time.Time
	NewTime      time.Time
}

type Calculator struct {
	baseDateTime time.Time
	photos       []PhotoTime
}

func NewCalculator(baseDateTime time.Time) *Calculator {
	return &Calculator{
		baseDateTime: baseDateTime,
		photos:       make([]PhotoTime, 0),
	}
}

func (c *Calculator) AddPhoto(filePath string, originalTime time.Time) {
	c.photos = append(c.photos, PhotoTime{
		FilePath:     filePath,
		OriginalTime: originalTime,
	})
}

func (c *Calculator) CalculateNewTimes() {
	if len(c.photos) == 0 {
		return
	}

	// Sort photos by their original timestamps
	sort.Slice(c.photos, func(i, j int) bool {
		return c.photos[i].OriginalTime.Before(c.photos[j].OriginalTime)
	})

	// The first photo gets the base time
	firstOriginalTime := c.photos[0].OriginalTime
	c.photos[0].NewTime = c.baseDateTime

	fmt.Printf("DEBUG: First photo time: %s -> %s\n",
		firstOriginalTime.Format("2006-01-02 15:04:05"),
		c.baseDateTime.Format("2006-01-02 15:04:05"))

	// Calculate time differences for subsequent photos
	for i := 1; i < len(c.photos); i++ {
		// Calculate time difference from the first photo
		timeDiff := c.photos[i].OriginalTime.Sub(firstOriginalTime)

		// Add this difference to the base time
		c.photos[i].NewTime = c.baseDateTime.Add(timeDiff)

		fmt.Printf("DEBUG: Photo %d: %s -> %s (diff: %v)\n",
			i+1,
			c.photos[i].OriginalTime.Format("2006-01-02 15:04:05"),
			c.photos[i].NewTime.Format("2006-01-02 15:04:05"),
			timeDiff)
	}
}

func (c *Calculator) GetPhotoNewTime(filePath string) (time.Time, bool) {
	for _, photo := range c.photos {
		if photo.FilePath == filePath {
			return photo.NewTime, true
		}
	}
	return time.Time{}, false
}

func (c *Calculator) GetPhotoCount() int {
	return len(c.photos)
}
