package appsystem

import (
	"testing"

	"bamort/database"
	"bamort/testutils"
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

func TestGetInfo2UsesDatabaseValues(t *testing.T) {
	testutils.SetupTestEnvironment(t)
	database.ResetTestDB()
	t.Cleanup(database.ResetTestDB)
	database.SetupTestDB()

	info := GetInfo2()

	if info.CharCount <= 0 {
		t.Fatalf("expected character count from database, got %d", info.CharCount)
	}

	if info.UserCount <= 0 {
		t.Fatalf("expected user count from database, got %d", info.UserCount)
	}

	if info.DbVersion == "" {
		t.Fatalf("expected database schema version, got empty string")
	}
}
