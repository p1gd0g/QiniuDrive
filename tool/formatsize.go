package tool

import "strconv"

// FormatSize formats the file size.
func FormatSize(n int64) string {
	if n < 1024 {
		return strconv.FormatInt(n, 10) + " B"
	}

	nf := float64(n)

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " KB"
	}

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " MB"
	}

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " GB"
	}

	nf /= 1024
	return strconv.FormatFloat(nf, 'f', 2, 64) + " TB"
}
