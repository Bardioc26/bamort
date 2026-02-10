package importer

// Swagger model definitions for API documentation

// DetectResponse represents the format detection response
// swagger:model
type DetectResponse struct {
	// Unique identifier of the detected adapter
	// example: moam-vtt-v1
	AdapterID string `json:"adapter_id"`

	// Confidence score between 0.0 and 1.0
	// example: 0.95
	Confidence float64 `json:"confidence"`

	// Human-readable name of the suggested adapter
	// example: Moam VTT Character
	SuggestedAdapterName string `json:"suggested_adapter_name,omitempty"`

	// Detected version of the external format
	// example: 10.x
	Version string `json:"version,omitempty"`
}

// ImportResultResponse represents a successful import response
// swagger:model
type ImportResultResponse struct {
	// ID of the created character
	// example: 123
	CharacterID uint `json:"character_id"`

	// ID of the import history record
	// example: 456
	ImportID uint `json:"import_id"`

	// Adapter used for import
	// example: moam-vtt-v1
	AdapterID string `json:"adapter_id"`

	// Import status
	// example: success
	// enum: success,partial,failed
	Status string `json:"status"`

	// Validation warnings (non-blocking)
	Warnings []ValidationWarning `json:"warnings,omitempty"`

	// Count of created master data items by type
	// example: {"skills": 3, "spells": 1, "equipment": 2}
	CreatedItems map[string]int `json:"created_items,omitempty"`
}

// AdapterListResponse represents the list of available adapters
// swagger:model
type AdapterListResponse struct {
	// List of registered adapters
	Adapters []AdapterMetadata `json:"adapters"`
}

// ImportHistoryResponse represents paginated import history
// swagger:model
type ImportHistoryResponse struct {
	// List of import history records
	Histories []ImportHistory `json:"histories"`

	// Total number of imports
	// example: 42
	Total int64 `json:"total"`

	// Current page number
	// example: 1
	Page int `json:"page"`

	// Items per page
	// example: 20
	PerPage int `json:"per_page"`

	// Total number of pages
	// example: 3
	Pages int64 `json:"pages"`
}

// ImportDetailsResponse represents detailed import information
// swagger:model
type ImportDetailsResponse struct {
	// Import history record
	History ImportHistory `json:"history"`

	// Master data items created/matched during import
	MasterDataImports []MasterDataImport `json:"master_data_imports"`
}

// ErrorResponse represents an API error
// swagger:model
type ErrorResponse struct {
	// Error message
	// example: Character not found
	Error string `json:"error"`

	// Optional hint for resolving the error
	// example: Specify adapter_id query parameter
	Hint string `json:"hint,omitempty"`

	// Available options (for 409 Conflict responses)
	AvailableAdapters []string `json:"available_adapters,omitempty"`
}

// ValidationWarningResponse represents a validation warning
// swagger:model
type ValidationWarningResponse struct {
	// Field that triggered the warning
	// example: Stats.St
	Field string `json:"field"`

	// Warning message
	// example: Stat value 101 exceeds typical range (0-100)
	Message string `json:"message"`

	// Source of the validation
	// example: gamesystem
	// enum: adapter,bmrt,gamesystem
	Source string `json:"source"`
}

// AdapterMetadataResponse represents adapter information
// swagger:model
type AdapterMetadataResponse struct {
	// Unique adapter identifier
	// example: moam-vtt-v1
	ID string `json:"id"`

	// Human-readable adapter name
	// example: Moam VTT Character
	Name string `json:"name"`

	// Adapter version (semantic versioning)
	// example: 1.0
	Version string `json:"version"`

	// Supported BMRT format versions
	// example: ["1.0"]
	BmrtVersions []string `json:"bmrt_versions"`

	// Supported file extensions
	// example: [".json"]
	SupportedExtensions []string `json:"supported_extensions"`

	// Adapter capabilities
	// example: ["import", "export", "detect"]
	Capabilities []string `json:"capabilities"`

	// Current health status
	// example: true
	Healthy bool `json:"healthy"`

	// Last health check timestamp
	// example: 2026-02-10T10:30:00Z
	LastCheckedAt string `json:"last_checked_at"`

	// Last error message (if unhealthy)
	// example: Connection timeout
	LastError string `json:"last_error,omitempty"`

	// Adapter base URL
	// example: http://adapter-moam:8181
	BaseURL string `json:"base_url"`
}

// ImportHistoryRecord represents an import history entry
// swagger:model
type ImportHistoryRecord struct {
	// Import history ID
	// example: 456
	ID uint `json:"id"`

	// User who performed the import
	// example: 123
	UserID uint `json:"user_id"`

	// Created character ID
	// example: 789
	CharacterID uint `json:"character_id"`

	// Adapter used for import
	// example: moam-vtt-v1
	AdapterID string `json:"adapter_id"`

	// Source format name
	// example: moam-vtt
	SourceFormat string `json:"source_format,omitempty"`

	// Original filename
	// example: character.json
	SourceFilename string `json:"source_filename"`

	// BMRT format version used
	// example: 1.0
	BmrtVersion string `json:"bmrt_version"`

	// Import timestamp
	// example: 2026-02-10T10:00:00Z
	ImportedAt string `json:"imported_at"`

	// Import status
	// example: success
	// enum: in_progress,success,partial,failed
	Status string `json:"status"`

	// Error log (if failed)
	ErrorLog string `json:"error_log,omitempty"`
}

// MasterDataImportRecord represents a master data import entry
// swagger:model
type MasterDataImportRecord struct {
	// Master data import ID
	// example: 1001
	ID uint `json:"id"`

	// Related import history ID
	// example: 456
	ImportHistoryID uint `json:"import_history_id"`

	// Type of master data
	// example: skill
	// enum: skill,spell,weapon,equipment
	ItemType string `json:"item_type"`

	// Master data item ID in BaMoRT database
	// example: 789
	ItemID uint `json:"item_id"`

	// Original name from external format
	// example: Custom Sword Fighting
	ExternalName string `json:"external_name"`

	// How the item was matched/created
	// example: created_personal
	// enum: exact,created_personal
	MatchType string `json:"match_type"`

	// Creation timestamp
	// example: 2026-02-10T10:00:05Z
	CreatedAt string `json:"created_at"`
}
