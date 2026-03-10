package tools

import "fmt"

const (
	bytesPerKB = 1024
	bytesPerMB = bytesPerKB * 1024
	bytesPerGB = bytesPerMB * 1024
	bytesPerTB = bytesPerGB * 1024

	percentMultiplier = 100
)

// formatBytes formats bytes into human-readable form.
func formatBytes(bytes int64) string {
	if bytes < 0 {
		return "unknown"
	}

	switch {
	case bytes >= bytesPerTB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(bytesPerTB))
	case bytes >= bytesPerGB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(bytesPerGB))
	case bytes >= bytesPerMB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(bytesPerMB))
	case bytes >= bytesPerKB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(bytesPerKB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// formatSpeed formats bytes per second into human-readable speed.
func formatSpeed(bytesPerSec int64) string {
	return formatBytes(bytesPerSec) + "/s"
}

// formatPercent formats a 0.0-1.0 float as a percentage string.
func formatPercent(ratio float64) string {
	return fmt.Sprintf("%.1f%%", ratio*percentMultiplier)
}
