package utils

import "fmt"

func BytesToHumanReadable(bytes int) string {
	const (
		megabyte = 1024 * 1024
		gigabyte = 1024 * megabyte
		terabyte = 1024 * gigabyte
	)

	if bytes >= terabyte {
		tera := float64(bytes) / terabyte
		return fmt.Sprintf("%.2f tera", tera)
	} else if bytes >= gigabyte {
		gig := float64(bytes) / gigabyte
		return fmt.Sprintf("%.2f gig", gig)
	} else if bytes >= megabyte {
		mb := float64(bytes) / megabyte
		return fmt.Sprintf("%.2f mb", mb)
	} else {
		return fmt.Sprintf("%d bytes", bytes)
	}

}
