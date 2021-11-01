package healthcheck

import (
	"strconv"
	"strings"
)

// nolint
const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

const (
	bitSize = 64
)

// BytesToString returns a human-readable byte string of the form 10M, 12.5K, and so forth.
// The following units are available:
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func BytesToString(bytes uint64) string {
	suffix := "B"
	value := float64(bytes)

	switch {
	case bytes >= TERABYTE:
		suffix = "T" + suffix
		value /= TERABYTE
	case bytes >= GIGABYTE:
		suffix = "G" + suffix
		value /= GIGABYTE
	case bytes >= MEGABYTE:
		suffix = "M" + suffix
		value /= MEGABYTE
	case bytes >= KILOBYTE:
		suffix = "K" + suffix
		value /= KILOBYTE
	case bytes >= BYTE:
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, bitSize)
	result = strings.TrimSuffix(result, ".0")

	return result + suffix
}
