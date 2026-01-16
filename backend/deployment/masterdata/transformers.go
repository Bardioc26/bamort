package masterdata

import (
"bamort/logger"
"fmt"
)

// ImportTransformer transforms export data from one version to another
type ImportTransformer interface {
CanTransform(exportVersion string) bool
Transform(data *ExportData) (*ExportData, error)
TargetVersion() string
}

// transformerRegistry holds all registered transformers
var transformerRegistry = []ImportTransformer{
// Add transformers here as needed
// Example: &V1ToV2Transformer{},
}

// TransformToCurrentVersion transforms export data to the current version
func TransformToCurrentVersion(data *ExportData) (*ExportData, error) {
if data.ExportVersion == CurrentExportVersion {
logger.Debug("Export already at current version %s", CurrentExportVersion)
return data, nil
}

logger.Info("Transforming export from version %s to %s", data.ExportVersion, CurrentExportVersion)

// Apply transformers in sequence
currentVersion := data.ExportVersion
transformedData := data

for currentVersion != CurrentExportVersion {
transformed := false

for _, transformer := range transformerRegistry {
if transformer.CanTransform(currentVersion) {
logger.Debug("Applying transformer: %s → %s", currentVersion, transformer.TargetVersion())

var err error
transformedData, err = transformer.Transform(transformedData)
if err != nil {
return nil, fmt.Errorf("transformation failed (%s → %s): %w",
currentVersion, transformer.TargetVersion(), err)
}

currentVersion = transformedData.ExportVersion
transformed = true
break
}
}

if !transformed {
return nil, fmt.Errorf("no transformer found for version %s", currentVersion)
}
}

logger.Info("Transformation complete: %s → %s", data.ExportVersion, CurrentExportVersion)
return transformedData, nil
}

// RegisterTransformer adds a transformer to the registry
func RegisterTransformer(transformer ImportTransformer) {
transformerRegistry = append(transformerRegistry, transformer)
}
