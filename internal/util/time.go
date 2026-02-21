// Package util provides utility functions.
package util

import "time"

// FormatRelativeTime formats a time as a relative time string.
func FormatRelativeTime(t time.Time) string {
	now := time.Now()
	if t.Location() != nil {
		now = now.In(t.Location())
	}

	diff := now.Sub(t)
	seconds := diff.Seconds()

	switch {
	case seconds < 60:
		return "just now"
	case seconds < 3600:
		minutes := int(seconds / 60)
		return formatDuration(minutes, "m")
	case seconds < 86400:
		hours := int(seconds / 3600)
		return formatDuration(hours, "h")
	case seconds < 604800:
		days := int(seconds / 86400)
		return formatDuration(days, "d")
	default:
		weeks := int(seconds / 604800)
		return formatDuration(weeks, "w")
	}
}

// FormatRelativeTimeShort formats a time as a short relative time string.
func FormatRelativeTimeShort(t time.Time) string {
	now := time.Now()
	if t.Location() != nil {
		now = now.In(t.Location())
	}

	diff := now.Sub(t)
	seconds := diff.Seconds()

	switch {
	case seconds < 60:
		return "now"
	case seconds < 3600:
		minutes := int(seconds / 60)
		return formatDurationShort(minutes, "m")
	case seconds < 86400:
		hours := int(seconds / 3600)
		return formatDurationShort(hours, "h")
	case seconds < 604800:
		days := int(seconds / 86400)
		return formatDurationShort(days, "d")
	default:
		weeks := int(seconds / 604800)
		return formatDurationShort(weeks, "w")
	}
}

func formatDuration(value int, unit string) string {
	return formatDurationShort(value, unit) + " ago"
}

func formatDurationShort(value int, unit string) string {
	return itoa(value) + unit
}

// Simple int to string without importing strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	negative := i < 0
	if negative {
		i = -i
	}

	var buf [20]byte
	pos := len(buf)

	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}

	if negative {
		pos--
		buf[pos] = '-'
	}

	return string(buf[pos:])
}
