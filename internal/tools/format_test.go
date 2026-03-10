package tools

import "testing"

func TestFormatBytes_Units(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"zero", 0, "0 B"},
		{"bytes", 512, "512 B"},
		{"kilobytes", 1536, "1.50 KB"},
		{"megabytes", 1572864, "1.50 MB"},
		{"gigabytes", 1610612736, "1.50 GB"},
		{"terabytes", 1649267441664, "1.50 TB"},
		{"negative", -1, "unknown"},
		{"large negative", -1099511627776, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("formatBytes(%d) = %s, want %s", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatSpeed(t *testing.T) {
	got := formatSpeed(1024)
	if got != "1.00 KB/s" {
		t.Errorf("formatSpeed(1024) = %s, want 1.00 KB/s", got)
	}
}

func TestFormatPercent(t *testing.T) {
	got := formatPercent(0.5)
	if got != "50.0%" {
		t.Errorf("formatPercent(0.5) = %s, want 50.0%%", got)
	}
}
