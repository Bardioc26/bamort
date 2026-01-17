package config

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Version should not be empty")
	}
	if version != Version {
		t.Errorf("Expected version %s, got %s", Version, version)
	}
}
