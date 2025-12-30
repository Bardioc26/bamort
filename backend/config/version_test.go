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

/*
	func TestGetGitCommit(t *testing.T) {
		commit := GetGitCommit()
		if commit == "" {
			t.Error("GitCommit should not be empty")
		}
		// Should be either "unknown" or a valid git hash
		if commit != "unknown" && len(commit) < 7 {
			t.Errorf("Invalid git commit format: %s", commit)
		}
	}
*/
func TestGetInfo(t *testing.T) {
	info := GetInfo()

	if info.Version == "" {
		t.Error("Info.Version should not be empty")
	}

	if info.GitCommit == "" {
		t.Error("Info.GitCommit should not be empty")
	}

	if info.Version != Version {
		t.Errorf("Expected info.Version %s, got %s", Version, info.Version)
	}
}
