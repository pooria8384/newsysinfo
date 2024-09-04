package utils

import (
	"fmt"
)

func ToHuman(n float32, counter int) string {
	if n < 1024 {
		return fmt.Sprintf("%.2f %s", float32(n), getUnit(counter))
	}
	return ToHuman(float32(n)/1024, counter+1)
}

func getUnit(n int) string {
	switch n {
	case 0:
		return "B"
	case 1:
		return "KB"
	case 2:
		return "MB"
	case 3:
		return "GB"
	case 4:
		return "TB"
	default:
		return "B"
	}
}
