package version

import (
	"bamort/config"
	"fmt"
	"strconv"
	"strings"
)

// RequiredDBVersion defines the exact database version this backend requires
// This must be updated whenever database migrations are added
const RequiredDBVersion = "0.4.0"

// VersionCompatibility contains version comparison results
type VersionCompatibility struct {
	BackendVersion    string
	RequiredDBVersion string
	ActualDBVersion   string
	Compatible        bool
	MigrationNeeded   bool
	Reason            string
}

// CheckCompatibility checks if the database version matches the required version
func CheckCompatibility(actualDBVersion string) *VersionCompatibility {
	compatible := actualDBVersion == RequiredDBVersion
	migrationNeeded := actualDBVersion != RequiredDBVersion

	var reason string
	if compatible {
		reason = "Database version matches required version"
	} else if isOlderVersion(actualDBVersion, RequiredDBVersion) {
		reason = fmt.Sprintf("Database migration required: %s → %s",
			actualDBVersion, RequiredDBVersion)
	} else {
		reason = fmt.Sprintf("Backend too old for database version. Backend requires %s, database is %s",
			RequiredDBVersion, actualDBVersion)
	}

	return &VersionCompatibility{
		BackendVersion:    config.GetVersion(),
		RequiredDBVersion: RequiredDBVersion,
		ActualDBVersion:   actualDBVersion,
		Compatible:        compatible,
		MigrationNeeded:   migrationNeeded,
		Reason:            reason,
	}
}

// GetRequiredDBVersion returns the database version this backend requires
func GetRequiredDBVersion() string {
	return RequiredDBVersion
}

// GetBackendVersion returns the current backend version
func GetBackendVersion() string {
	return config.GetVersion()
}

// parseVersion parses a semantic version string into major, minor, patch
func parseVersion(version string) (major, minor, patch int, err error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return major, minor, patch, nil
}

// isOlderVersion checks if version1 is older than version2
func isOlderVersion(version1, version2 string) bool {
	v1Major, v1Minor, v1Patch, err1 := parseVersion(version1)
	v2Major, v2Minor, v2Patch, err2 := parseVersion(version2)

	if err1 != nil || err2 != nil {
		return false
	}

	if v1Major != v2Major {
		return v1Major < v2Major
	}
	if v1Minor != v2Minor {
		return v1Minor < v2Minor
	}
	return v1Patch < v2Patch
}

// CompareVersions returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) (int, error) {
	v1Major, v1Minor, v1Patch, err1 := parseVersion(v1)
	if err1 != nil {
		return 0, err1
	}

	v2Major, v2Minor, v2Patch, err2 := parseVersion(v2)
	if err2 != nil {
		return 0, err2
	}

	if v1Major != v2Major {
		if v1Major < v2Major {
			return -1, nil
		}
		return 1, nil
	}

	if v1Minor != v2Minor {
		if v1Minor < v2Minor {
			return -1, nil
		}
		return 1, nil
	}

	if v1Patch != v2Patch {
		if v1Patch < v2Patch {
			return -1, nil
		}
		return 1, nil
	}

	return 0, nil
}
