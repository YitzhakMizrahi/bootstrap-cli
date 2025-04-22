package system

import (
	"testing"
)

func TestDetect(t *testing.T) {
	info, err := Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	if info.OS == "" {
		t.Error("Detect() OS is empty")
	}

	if info.Arch == "" {
		t.Error("Detect() Arch is empty")
	}
} 