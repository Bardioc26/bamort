package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantMajor   int
		wantMinor   int
		wantPatch   int
		expectError bool
	}{
		{
			name:        "Valid version 0.4.0",
			version:     "0.4.0",
			wantMajor:   0,
			wantMinor:   4,
			wantPatch:   0,
			expectError: false,
		},
		{
			name:        "Valid version 1.2.3",
			version:     "1.2.3",
			wantMajor:   1,
			wantMinor:   2,
			wantPatch:   3,
			expectError: false,
		},
		{
			name:        "Invalid version - too few parts",
			version:     "1.2",
			expectError: true,
		},
		{
			name:        "Invalid version - too many parts",
			version:     "1.2.3.4",
			expectError: true,
		},
		{
			name:        "Invalid version - non-numeric",
			version:     "1.a.3",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, patch, err := parseVersion(tt.version)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMajor, major)
				assert.Equal(t, tt.wantMinor, minor)
				assert.Equal(t, tt.wantPatch, patch)
			}
		})
	}
}

func TestIsOlderVersion(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		want     bool
	}{
		{
			name:     "0.3.0 is older than 0.4.0",
			version1: "0.3.0",
			version2: "0.4.0",
			want:     true,
		},
		{
			name:     "0.4.0 is not older than 0.4.0",
			version1: "0.4.0",
			version2: "0.4.0",
			want:     false,
		},
		{
			name:     "0.5.0 is not older than 0.4.0",
			version1: "0.5.0",
			version2: "0.4.0",
			want:     false,
		},
		{
			name:     "0.4.1 is not older than 0.4.0",
			version1: "0.4.1",
			version2: "0.4.0",
			want:     false,
		},
		{
			name:     "0.4.0 is older than 0.4.1",
			version1: "0.4.0",
			version2: "0.4.1",
			want:     true,
		},
		{
			name:     "0.4.0 is older than 1.0.0",
			version1: "0.4.0",
			version2: "1.0.0",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOlderVersion(tt.version1, tt.version2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		v1      string
		v2      string
		want    int
		wantErr bool
	}{
		{
			name: "Equal versions",
			v1:   "0.4.0",
			v2:   "0.4.0",
			want: 0,
		},
		{
			name: "v1 < v2 (minor)",
			v1:   "0.3.0",
			v2:   "0.4.0",
			want: -1,
		},
		{
			name: "v1 > v2 (minor)",
			v1:   "0.5.0",
			v2:   "0.4.0",
			want: 1,
		},
		{
			name: "v1 < v2 (patch)",
			v1:   "0.4.0",
			v2:   "0.4.1",
			want: -1,
		},
		{
			name: "v1 > v2 (patch)",
			v1:   "0.4.1",
			v2:   "0.4.0",
			want: 1,
		},
		{
			name: "v1 < v2 (major)",
			v1:   "0.9.0",
			v2:   "1.0.0",
			want: -1,
		},
		{
			name:    "Invalid v1",
			v1:      "invalid",
			v2:      "0.4.0",
			wantErr: true,
		},
		{
			name:    "Invalid v2",
			v1:      "0.4.0",
			v2:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareVersions(tt.v1, tt.v2)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCheckCompatibility(t *testing.T) {
	tests := []struct {
		name                string
		actualDBVersion     string
		wantCompatible      bool
		wantMigrationNeeded bool
		reasonContains      string
	}{
		{
			name:                "Exact match - compatible",
			actualDBVersion:     RequiredDBVersion,
			wantCompatible:      true,
			wantMigrationNeeded: false,
			reasonContains:      "matches required version",
		},
		{
			name:                "DB too old - migration needed",
			actualDBVersion:     "0.3.0",
			wantCompatible:      false,
			wantMigrationNeeded: true,
			reasonContains:      "migration required",
		},
		{
			name:                "DB too new - backend too old",
			actualDBVersion:     "0.5.0",
			wantCompatible:      false,
			wantMigrationNeeded: true,
			reasonContains:      "Backend too old",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckCompatibility(tt.actualDBVersion)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantCompatible, result.Compatible)
			assert.Equal(t, tt.wantMigrationNeeded, result.MigrationNeeded)
			assert.Equal(t, RequiredDBVersion, result.RequiredDBVersion)
			assert.Equal(t, tt.actualDBVersion, result.ActualDBVersion)
			assert.Contains(t, result.Reason, tt.reasonContains)
		})
	}
}

func TestGetRequiredDBVersion(t *testing.T) {
	version := GetRequiredDBVersion()
	assert.Equal(t, RequiredDBVersion, version)
	assert.NotEmpty(t, version)
}

func TestGetBackendVersion(t *testing.T) {
	version := GetBackendVersion()
	assert.NotEmpty(t, version)
	// Should be able to parse it
	_, _, _, err := parseVersion(version)
	assert.NoError(t, err)
}
