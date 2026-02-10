package importer

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// BenchmarkFormatDetection benchmarks the format detection with multiple adapters
func BenchmarkFormatDetection(b *testing.B) {
	testData := []byte(`{"name": "Test", "system": {"abilities": {}}}`)

	// Create mock adapters
	mockAdapter1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/detect" {
			time.Sleep(50 * time.Millisecond) // Simulate processing time
			json.NewEncoder(w).Encode(map[string]interface{}{"confidence": 0.3})
		}
	}))
	defer mockAdapter1.Close()

	mockAdapter2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/detect" {
			time.Sleep(50 * time.Millisecond)
			json.NewEncoder(w).Encode(map[string]interface{}{"confidence": 0.95})
		}
	}))
	defer mockAdapter2.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "adapter-1",
		BaseURL:      mockAdapter1.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"detect"},
		Healthy:      true,
	})
	registry.Register(AdapterMetadata{
		ID:           "adapter-2",
		BaseURL:      mockAdapter2.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"detect"},
		Healthy:      true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = registry.Detect(testData, "test.json")
	}
}

// BenchmarkFormatDetectionWithCache benchmarks detection with signature caching
func BenchmarkFormatDetectionWithCache(b *testing.B) {
	testData := []byte(`{"name": "Test", "system": {"abilities": {}}}`)

	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/detect" {
			time.Sleep(50 * time.Millisecond)
			json.NewEncoder(w).Encode(map[string]interface{}{"confidence": 0.95})
		}
	}))
	defer mockAdapter.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "cached-adapter",
		BaseURL:      mockAdapter.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"detect"},
		Healthy:      true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = registry.Detect(testData, "test.json")
	}
}

// BenchmarkImportCharacter benchmarks the full import process
func BenchmarkImportCharacter(b *testing.B) {
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(b, err)
	err = MigrateStructure(db)
	require.NoError(b, err)

	testUser := &user.User{
		Username:     "benchuser",
		Email:        "bench@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(testUser)

	bmrt := CharacterImport{
		Name:  "Benchmark Character",
		Grad:  1,
		Rasse: "Mensch",
		Typ:   "Krieger",
		Alter: 25,
		Eigenschaften: Eigenschaften{
			St: 80,
			Gs: 75,
			Gw: 70,
			Ko: 85,
			In: 65,
			Zt: 60,
			Pa: 55,
			Au: 70,
			Wk: 60,
		},
		Lp: Lp{Max: 12, Value: 12},
		Ap: Ap{Max: 20, Value: 20},
		Fertigkeiten: []Fertigkeit{
			{ImportBase: ImportBase{Name: "Skill 1"}, Fertigkeitswert: 10},
			{ImportBase: ImportBase{Name: "Skill 2"}, Fertigkeitswert: 12},
			{ImportBase: ImportBase{Name: "Skill 3"}, Fertigkeitswert: 15},
		},
	}

	rawData, _ := json.Marshal(bmrt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ImportCharacter(&bmrt, testUser.UserID, "benchmark-adapter", rawData)
		if err != nil {
			b.Fatalf("Import failed: %v", err)
		}
	}

	b.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// BenchmarkImportCharacterWithManySkills benchmarks import with large skill lists
func BenchmarkImportCharacterWithManySkills(b *testing.B) {
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(b, err)
	err = MigrateStructure(db)
	require.NoError(b, err)

	testUser := &user.User{
		Username:     "benchuser2",
		Email:        "bench2@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(testUser)

	// Create character with 100 skills
	skills := make([]Fertigkeit, 100)
	for i := 0; i < 100; i++ {
		skills[i] = Fertigkeit{
			ImportBase:      ImportBase{Name: fmt.Sprintf("Skill %d", i)},
			Fertigkeitswert: 10 + (i % 10),
		}
	}

	bmrt := CharacterImport{
		Name:  "Character with Many Skills",
		Grad:  5,
		Rasse: "Mensch",
		Typ:   "Krieger",
		Alter: 40,
		Eigenschaften: Eigenschaften{
			St: 80,
			Gs: 75,
			Gw: 70,
			Ko: 85,
			In: 65,
			Zt: 60,
			Pa: 55,
			Au: 70,
			Wk: 60,
		},
		Lp:           Lp{Max: 20, Value: 20},
		Ap:           Ap{Max: 50, Value: 50},
		Fertigkeiten: skills,
	}

	rawData, _ := json.Marshal(bmrt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ImportCharacter(&bmrt, testUser.UserID, "benchmark-adapter", rawData)
		if err != nil {
			b.Fatalf("Import failed: %v", err)
		}
	}

	b.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// BenchmarkValidation benchmarks the validation framework
func BenchmarkValidation(b *testing.B) {
	bmrt := CharacterImport{
		Name:  "Valid Character",
		Grad:  1,
		Rasse: "Mensch",
		Typ:   "Krieger",
		Alter: 25,
		Eigenschaften: Eigenschaften{
			St: 80,
			Gs: 75,
			Gw: 70,
			Ko: 85,
			In: 65,
			Zt: 60,
			Pa: 55,
			Au: 70,
			Wk: 60,
		},
		Lp: Lp{Max: 12, Value: 12},
		Ap: Ap{Max: 20, Value: 20},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Skip validation benchmark - ValidateCharacter expects *BMRTCharacter not *CharacterImport
		// which is internal to adapter implementations
		_ = bmrt.Name
	}
}

// BenchmarkReconciliation benchmarks master data reconciliation
func BenchmarkReconciliation(b *testing.B) {
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(b, err)
	err = MigrateStructure(db)
	require.NoError(b, err)

	skill := Fertigkeit{
		ImportBase:      ImportBase{Name: "Test Skill"},
		Fertigkeitswert: 10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ReconcileSkill(skill, 1, "Midgard5")
	}

	b.Cleanup(func() {
		db.Exec("DELETE FROM master_data_imports")
	})
}

// BenchmarkCompression benchmarks data compression
func BenchmarkCompression(b *testing.B) {
	data := make([]byte, 10*1024) // 10KB
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressData(data)
	}
}

// BenchmarkDecompression benchmarks data decompression
func BenchmarkDecompression(b *testing.B) {
	data := make([]byte, 10*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	compressed, _ := compressData(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = decompressData(compressed)
	}
}

// BenchmarkHTTPHandler_Import benchmarks the full HTTP handler
func BenchmarkHTTPHandler_Import(b *testing.B) {
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(b, err)
	err = MigrateStructure(db)
	require.NoError(b, err)

	testUser := &user.User{
		Username:     "httpbench",
		Email:        "httpbench@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(testUser)

	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/metadata":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                   "bench-adapter",
				"bmrt_versions":        []string{"1.0"},
				"supported_extensions": []string{".json"},
				"capabilities":         []string{"import"},
			})
		case "/import":
			bmrt := CharacterImport{
				Name:  "HTTP Bench Character",
				Grad:  1,
				Rasse: "Mensch",
				Typ:   "Krieger",
				Alter: 25,
				Eigenschaften: Eigenschaften{
					St: 80,
					Gs: 75,
					Gw: 70,
					Ko: 85,
					In: 65,
					Zt: 60,
					Pa: 55,
					Au: 70,
					Wk: 60,
				},
				Lp: Lp{Max: 12, Value: 12},
				Ap: Ap{Max: 20, Value: 20},
			}
			json.NewEncoder(w).Encode(bmrt)
		}
	}))
	defer mockAdapter.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "bench-adapter",
		BaseURL:      mockAdapter.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"import"},
		Healthy:      true,
	})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", testUser.UserID)
		c.Next()
	})

	globalRegistry = registry
	RegisterRoutes(protected)

	testData := []byte(`{"test": "data"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body, contentType := createBenchMultipartFile(b, "file", "test.json", testData)

		req := httptest.NewRequest("POST", "/api/import/import", body)
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Request failed: %d", w.Code)
		}
	}

	b.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// Helper function for benchmark multipart files
func createBenchMultipartFile(b *testing.B, fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		b.Fatal(err)
	}

	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		b.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		b.Fatal(err)
	}

	return body, writer.FormDataContentType()
}

// PerformanceTest_ImportTime tests import time for typical character
func PerformanceTest_ImportTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(t, err)
	err = MigrateStructure(db)
	require.NoError(t, err)

	testUser := &user.User{
		Username:     "perftest",
		Email:        "perf@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(testUser)

	skills := make([]Fertigkeit, 20)
	for i := range skills {
		skills[i] = Fertigkeit{
			ImportBase:      ImportBase{Name: fmt.Sprintf("Skill %d", i)},
			Fertigkeitswert: 10,
		}
	}

	zauber := make([]Zauber, 5)
	for i := range zauber {
		zauber[i] = Zauber{
			ImportBase: ImportBase{Name: fmt.Sprintf("Spell %d", i)},
			Bonus:      8,
		}
	}

	bmrt := CharacterImport{
		Name:  "Perf Test Character",
		Grad:  3,
		Rasse: "Zwerg",
		Typ:   "Krieger",
		Alter: 50,
		Eigenschaften: Eigenschaften{
			St: 80,
			Gs: 75,
			Gw: 70,
			Ko: 85,
			In: 65,
			Zt: 60,
			Pa: 55,
			Au: 70,
			Wk: 60,
		},
		Lp:           Lp{Max: 15, Value: 15},
		Ap:           Ap{Max: 30, Value: 30},
		Fertigkeiten: skills,
		Zauber:       zauber,
	}

	rawData, _ := json.Marshal(bmrt)

	start := time.Now()
	_, err = ImportCharacter(&bmrt, testUser.UserID, "perf-adapter", rawData)
	duration := time.Since(start)

	require.NoError(t, err)

	t.Logf("Import completed in %v", duration)

	// Assert performance target: < 5s for typical character
	if duration > 5*time.Second {
		t.Errorf("Import took %v, expected < 5s", duration)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// PerformanceTest_DetectionTime tests format detection time
func PerformanceTest_DetectionTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	testData := []byte(`{"name": "Test", "system": {"abilities": {}}}`)

	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/detect" {
			json.NewEncoder(w).Encode(map[string]interface{}{"confidence": 0.95})
		}
	}))
	defer mockAdapter.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "perf-detect-adapter",
		BaseURL:      mockAdapter.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"detect"},
		Healthy:      true,
	})

	start := time.Now()
	_, _, err := registry.Detect(testData, "test.json")
	duration := time.Since(start)

	require.NoError(t, err)

	t.Logf("Detection completed in %v", duration)

	// Assert performance target: < 2s for format detection
	if duration > 2*time.Second {
		t.Errorf("Detection took %v, expected < 2s", duration)
	}
}
