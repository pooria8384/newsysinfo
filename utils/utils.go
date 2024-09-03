package utils

import "fmt"

func KbToHumanReadable(kb uint) string {
	const (
		kilobyte = 1
		megabyte = 1024
		gigabyte = 1024 * megabyte
		terabyte = 1024 * gigabyte
		petabyte = 1024 * terabyte
		exabyte  = 1024 * petabyte
	)
	if kb >= exabyte {
		exa := float64(kb) / exabyte
		return fmt.Sprintf("%.2f EB", exa)
	} else if kb >= petabyte {
		peta := float64(kb) / petabyte
		return fmt.Sprintf("%.2f PB", peta)
	} else if kb >= terabyte {
		tera := float64(kb) / terabyte
		return fmt.Sprintf("%.2f TB", tera)
	} else if kb >= gigabyte {
		gig := float64(kb) / gigabyte
		return fmt.Sprintf("%.2f GB", gig)
	} else if kb >= megabyte {
		mb := float64(kb) / megabyte
		return fmt.Sprintf("%.2f MB", mb)
	} else {
		return fmt.Sprintf("%d KB", kb)
	}

}
