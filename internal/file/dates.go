package file

import (
	"fmt"
	"regexp"
	"time"
)

func ParseDate(fileName string) (time.Time, error) {
	dateStr, err := findDate(fileName)
	if err != nil {
		return time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC), fmt.Errorf("failed to find date: %w", err)
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC), fmt.Errorf("failed to parse date: %w", err)
	}
	return date, err
}

func findDate(str string) (string, error) {
	const datePattern = `\d{4}-\d{2}-\d{2}` // yyyy-mm-dd
	r, err := regexp.Compile(datePattern)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}
	dateStr := r.FindString(str)
	if dateStr == "" {
		return "", fmt.Errorf("%s doesn't contain date with pattern yyyy-mm-dd", str)
	}
	return dateStr, nil
}
