package equipment

import (
	"bamort/database"
	"bamort/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		status         int
		message        string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Bad Request Error",
			status:         http.StatusBadRequest,
			message:        "Invalid request",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Invalid request"},
		},
		{
			name:           "Internal Server Error",
			status:         http.StatusInternalServerError,
			message:        "Database connection failed",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Database connection failed"},
		},
		{
			name:           "Not Found Error",
			status:         http.StatusNotFound,
			message:        "Resource not found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "Resource not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			respondWithError(c, tt.status, tt.message)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestCreateAusruestung(t *testing.T) {
	database.SetupTestDB(true)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name: "Valid Ausruestung Creation",
			payload: models.EqAusruestung{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Test Sword",
					},
					CharacterID: 21,
					UserID:      4,
				},
				Magisch: models.Magisch{
					IstMagisch:  false,
					Abw:         0,
					Ausgebrannt: false,
				},
				Beschreibung:  "A test sword",
				Anzahl:        1,
				BeinhaltetIn:  "",
				ContainedIn:   0,
				ContainerType: "",
				Bonus:         0,
				Gewicht:       2.5,
				Wert:          100.0,
			},
			expectedStatus: http.StatusCreated,
			shouldContain:  "Test Sword",
		},
		{
			name:           "Invalid JSON",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "error",
		},
		{
			name:           "Empty JSON",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusNotFound,
			shouldContain:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// u := user.User{}
			// u.FirstId(1)
			// token := user.GenerateToken(&u)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var body bytes.Buffer
			if tt.name == "Invalid JSON" {
				body = *bytes.NewBufferString("invalid json")
			} else {
				jsonData, _ := json.Marshal(tt.payload)
				body = *bytes.NewBuffer(jsonData)
			}

			c.Set("userID", uint(4)) // Simulate logged-in user with ID 4
			c.Request = httptest.NewRequest("POST", "/ausruestung", &body)
			c.Request.Header.Set("Content-Type", "application/json")

			//c.Request.Header.Set("Authorization", "Bearer "+token)

			CreateAusruestung(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.shouldContain != "" {
				assert.Contains(t, w.Body.String(), tt.shouldContain)
			}
		})
	}
}

func TestListAusruestung(t *testing.T) {
	database.SetupTestDB(true)
	gin.SetMode(gin.TestMode)

	// Create test data
	testAusruestung := models.EqAusruestung{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Test Equipment",
			},
			CharacterID: 123,
			UserID:      1,
		},
		Magisch: models.Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:  "Test equipment description",
		Anzahl:        1,
		BeinhaltetIn:  "",
		ContainedIn:   0,
		ContainerType: "",
		Bonus:         0,
		Gewicht:       1.0,
		Wert:          50.0,
	}

	err := database.DB.Create(&testAusruestung).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		characterID    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Valid Character ID",
			characterID:    "123",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Non-existent Character ID",
			characterID:    "999",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "Invalid Character ID",
			characterID:    "invalid",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				{Key: "character_id", Value: tt.characterID},
			}

			ListAusruestung(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []models.EqAusruestung
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)

			if tt.expectedCount > 0 {
				assert.Equal(t, "Test Equipment", response[0].Name)
				assert.Equal(t, uint(123), response[0].CharacterID)
			}
		})
	}
}

func TestUpdateAusruestung(t *testing.T) {
	database.SetupTestDB(true)
	gin.SetMode(gin.TestMode)

	// Create test data
	testAusruestung := models.EqAusruestung{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Original Equipment",
			},
			CharacterID: 21,
			UserID:      4,
		},
		Magisch: models.Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:  "Original description",
		Anzahl:        1,
		BeinhaltetIn:  "",
		ContainedIn:   0,
		ContainerType: "",
		Bonus:         0,
		Gewicht:       1.0,
		Wert:          50.0,
	}

	err := database.DB.Create(&testAusruestung).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		ausruestungID  string
		payload        interface{}
		expectedStatus int
		shouldContain  string
	}{
		{
			name:          "Valid Update",
			ausruestungID: strconv.Itoa(int(testAusruestung.ID)),
			payload: map[string]interface{}{
				"name":         "Updated Equipment",
				"beschreibung": "Updated description",
				"wert":         75.0,
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "Updated Equipment",
		},
		{
			name:          "Non-existent Ausruestung",
			ausruestungID: "999",
			payload: map[string]interface{}{
				"name": "Updated Equipment",
			},
			expectedStatus: http.StatusNotFound,
			shouldContain:  "error",
		},
		{
			name:           "Invalid JSON",
			ausruestungID:  strconv.Itoa(int(testAusruestung.ID)),
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("userID", uint(4))

			c.Params = gin.Params{
				{Key: "ausruestung_id", Value: tt.ausruestungID},
			}

			var body bytes.Buffer
			if tt.name == "Invalid JSON" {
				body = *bytes.NewBufferString("invalid json")
			} else {
				jsonData, _ := json.Marshal(tt.payload)
				body = *bytes.NewBuffer(jsonData)
			}

			c.Request = httptest.NewRequest("PUT", "/ausruestung/"+tt.ausruestungID, &body)
			c.Request.Header.Set("Content-Type", "application/json")

			UpdateAusruestung(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.shouldContain != "" {
				assert.Contains(t, w.Body.String(), tt.shouldContain)
			}
		})
	}
}

func TestDeleteAusruestung(t *testing.T) {
	database.SetupTestDB(true)
	gin.SetMode(gin.TestMode)

	// Create test data
	testAusruestung := models.EqAusruestung{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Equipment to Delete",
			},
			CharacterID: 21,
			UserID:      4,
		},
		Magisch: models.Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:  "Equipment for deletion test",
		Anzahl:        1,
		BeinhaltetIn:  "",
		ContainedIn:   0,
		ContainerType: "",
		Bonus:         0,
		Gewicht:       1.0,
		Wert:          50.0,
	}

	err := database.DB.Create(&testAusruestung).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		ausruestungID  string
		expectedStatus int
		shouldContain  string
	}{
		{
			name:           "Valid Deletion",
			ausruestungID:  strconv.Itoa(int(testAusruestung.ID)),
			expectedStatus: http.StatusOK,
			shouldContain:  "deleted successfully",
		},
		{
			name:           "Non-existent Ausruestung",
			ausruestungID:  "999",
			expectedStatus: http.StatusNotFound, // GORM doesn't fail on deleting non-existent records
			shouldContain:  "Ausruestung not found",
		},
		{
			name:           "Invalid Ausruestung ID",
			ausruestungID:  "invalid",
			expectedStatus: http.StatusNotFound,
			shouldContain:  "Ausruestung not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("userID", uint(4))

			c.Params = gin.Params{
				{Key: "ausruestung_id", Value: tt.ausruestungID},
			}
			DeleteAusruestung(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.shouldContain != "" {
				assert.Contains(t, w.Body.String(), tt.shouldContain)
			}

			// For successful deletion, verify the record is actually deleted
			if tt.name == "Valid Deletion" && tt.expectedStatus == http.StatusOK {
				var count int64
				database.DB.Model(&models.EqAusruestung{}).Where("id = ?", tt.ausruestungID).Count(&count)
				assert.Equal(t, int64(0), count, "Equipment should be deleted from database")
			}
		})
	}
}
