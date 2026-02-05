package appsystem

import (
	"bamort/database"
)

// Version is the application version
const Version = "0.2.4"

var (
	// GitCommit will be set by build flags or detected at runtime
	GitCommit = "unknown"
)

// init detects git commit if not set during build
func init() {
	/*
		if GitCommit == "" {
			// Try environment variable first
			if envCommit := os.Getenv("GIT_COMMIT"); envCommit != "" {
				GitCommit = envCommit
			} else {
				// Try to detect from git command
				GitCommit = detectGitCommit()
			}
		}
	*/
}

/*
// detectGitCommit tries to get the current git commit hash
func detectGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}
*/
// GetVersion returns the current application version
func GetVersion() string {
	return Version
}

/*
// GetGitCommit returns the git commit hash
func GetGitCommit() string {
	return GitCommit
}
*/
// Info contains version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
}
type Info2 struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	UserCount int64  `json:"userCount"`
	CharCount int64  `json:"charCount"`
	DbVersion string `json:"dbVersion"`
}

func getUserCount() int {
	// Placeholder implementation
	return 42
}
func getCharCount() int {
	// Placeholder implementation
	return 123456
}

// GetInfo returns version information as a struct
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
	}
}

func GetInfo2() Info2 {
	info := Info2{
		Version:   Version,
		GitCommit: GitCommit,
	}
	db := database.DB
	if db == nil {
		return info
	}

	db.Table("char_chars").Count(&info.CharCount)
	db.Table("users").Count(&info.UserCount)
	db.Raw("SELECT version FROM schema_version ORDER BY id DESC LIMIT 1").Scan(&info.DbVersion)

	return info
}
