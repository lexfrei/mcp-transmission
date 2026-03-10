package tools

import "github.com/lexfrei/go-transmission/api/transmission"

func deref(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}

	return fallback
}

func derefInt64(ptr *int64) int64 {
	if ptr != nil {
		return *ptr
	}

	return 0
}

func derefFloat64(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}

	return 0
}

func derefBool(ptr *bool) bool {
	if ptr != nil {
		return *ptr
	}

	return false
}

func derefInt(ptr *int) int {
	if ptr != nil {
		return *ptr
	}

	return 0
}

func derefStatus(ptr *transmission.TorrentStatus) string {
	if ptr != nil {
		return ptr.String()
	}

	return "Unknown"
}
