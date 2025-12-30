package config

// Version is the application version
const Version = "0.1.30"

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

// GetInfo returns version information as a struct
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
	}
}
