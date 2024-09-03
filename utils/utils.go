package utils

import (
	"fmt"
)

func KbToHumanReadable(kb uint) string {
	const (
		kilobyte = 1
		megabyte = 1024
		gigabyte = 1024 * megabyte
		terabyte = 1024 * gigabyte
		petabyte = 1024 * terabyte
		exabyte  = 1024 * petabyte
	)

	if kb == 0 {
		return "0 KB"
	}

	if kb > exabyte {
		return "Error: Value exceeds the exabyte limit"
	}

	switch {
	case kb >= exabyte:
		return fmt.Sprintf("%.2f EB", float64(kb)/float64(exabyte))
	case kb >= petabyte:
		return fmt.Sprintf("%.2f PB", float64(kb)/float64(petabyte))
	case kb >= terabyte:
		return fmt.Sprintf("%.2f TB", float64(kb)/float64(terabyte))
	case kb >= gigabyte:
		return fmt.Sprintf("%.2f GB", float64(kb)/float64(gigabyte))
	case kb >= megabyte:
		return fmt.Sprintf("%.2f MB", float64(kb)/float64(megabyte))
	default:
		return fmt.Sprintf("%d KB", kb)
	}

}
