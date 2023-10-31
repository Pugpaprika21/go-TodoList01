package helper

import (
	"strings"
	"time"
)

func NumRows(index int) int {
	return index + 1
}

func GetDMY(createAt time.Time) string {
	parts := strings.Split(createAt.String(), " ")
	if len(parts) >= 1 {
		datePart := parts[0]
		return datePart
	}
	return ""
}
