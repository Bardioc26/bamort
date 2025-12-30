package character

import (
	"testing"
)

func TestGetGradeBonus_Abwehr(t *testing.T) {
	tests := []struct {
		name     string
		grad     int
		expected int
	}{
		{"Grad 1", 1, 11},
		{"Grad 2", 2, 12},
		{"Grad 3", 3, 12},
		{"Grad 4", 4, 12},
		{"Grad 5", 5, 13},
		{"Grad 9", 9, 13},
		{"Grad 10", 10, 14},
		{"Grad 14", 14, 14},
		{"Grad 15", 15, 15},
		{"Grad 19", 19, 15},
		{"Grad 20", 20, 16},
		{"Grad 24", 24, 16},
		{"Grad 25", 25, 17},
		{"Grad 29", 29, 17},
		{"Grad 30", 30, 18},
		{"Grad 35", 35, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAbwehrBaseByGrade(tt.grad)
			if result != tt.expected {
				t.Errorf("Expected %d for Grad %d, got %d", tt.expected, tt.grad, result)
			}
		})
	}
}

func TestGetGradeBonus_Resistenz(t *testing.T) {
	tests := []struct {
		name     string
		grad     int
		expected int
	}{
		{"Grad 1", 1, 11},
		{"Grad 2", 2, 12},
		{"Grad 3", 3, 12},
		{"Grad 4", 4, 12},
		{"Grad 5", 5, 13},
		{"Grad 9", 9, 13},
		{"Grad 10", 10, 14},
		{"Grad 14", 14, 14},
		{"Grad 15", 15, 15},
		{"Grad 19", 19, 15},
		{"Grad 20", 20, 16},
		{"Grad 24", 24, 16},
		{"Grad 25", 25, 17},
		{"Grad 29", 29, 17},
		{"Grad 30", 30, 18},
		{"Grad 35", 35, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getResistenzBaseByGrade(tt.grad)
			if result != tt.expected {
				t.Errorf("Expected %d for Grad %d, got %d", tt.expected, tt.grad, result)
			}
		})
	}
}

func TestGetGradeBonus_Zaubern(t *testing.T) {
	tests := []struct {
		name     string
		grad     int
		expected int
	}{
		{"Grad 1", 1, 11},
		{"Grad 2", 2, 12},
		{"Grad 3", 3, 12},
		{"Grad 4", 4, 13},
		{"Grad 5", 5, 13},
		{"Grad 6", 6, 14},
		{"Grad 7", 7, 14},
		{"Grad 8", 8, 15},
		{"Grad 9", 9, 15},
		{"Grad 10", 10, 16},
		{"Grad 14", 14, 16},
		{"Grad 15", 15, 17},
		{"Grad 19", 19, 17},
		{"Grad 20", 20, 18},
		{"Grad 25", 25, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getZaubernBaseByGrade(tt.grad)
			if result != tt.expected {
				t.Errorf("Expected %d for Grad %d, got %d", tt.expected, tt.grad, result)
			}
		})
	}
}

func TestCalculateStaticFieldsLogic_WithGrade(t *testing.T) {
	tests := []struct {
		name                  string
		grad                  int
		expectedAbwehr        int
		expectedResistenzBody int
		expectedResistenzMind int
		expectedZaubern       int
	}{
		{
			name:                  "Grad 1",
			grad:                  1,
			expectedAbwehr:        11, // base 11 + bonus 0
			expectedResistenzBody: 12, // base 11 + bonus 1 (Krieger)
			expectedResistenzMind: 11, // base 11 + bonus 0 (Krieger gets no mind bonus)
			expectedZaubern:       11, // base 11 + bonus 0
		},
		{
			name:                  "Grad 5",
			grad:                  5,
			expectedAbwehr:        13, // base 13 + bonus 0
			expectedResistenzBody: 14, // base 13 + bonus 1 (Krieger)
			expectedResistenzMind: 13, // base 13 + bonus 0
			expectedZaubern:       13, // base 13 + bonus 0
		},
		{
			name:                  "Grad 10",
			grad:                  10,
			expectedAbwehr:        14, // base 14 + bonus 0
			expectedResistenzBody: 15, // base 14 + bonus 1 (Krieger)
			expectedResistenzMind: 14, // base 14 + bonus 0
			expectedZaubern:       16, // base 16 + bonus 0
		},
		{
			name:                  "Grad 20",
			grad:                  20,
			expectedAbwehr:        16, // base 16 + bonus 0
			expectedResistenzBody: 17, // base 16 + bonus 1 (Krieger)
			expectedResistenzMind: 16, // base 16 + bonus 0
			expectedZaubern:       18, // base 18 + bonus 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CalculateStaticFieldsRequest{
				St:    50,
				Gs:    50,
				Gw:    50,
				Ko:    50,
				In:    50,
				Zt:    50,
				Au:    50,
				Rasse: "Menschen",
				Typ:   "Krieger",
				Grad:  tt.grad,
			}

			result := CalculateStaticFieldsLogic(req)

			if result.Abwehr != tt.expectedAbwehr {
				t.Errorf("Expected Abwehr %d for Grad %d, got %d", tt.expectedAbwehr, tt.grad, result.Abwehr)
			}
			if result.ResistenzKoerper != tt.expectedResistenzBody {
				t.Errorf("Expected ResistenzKoerper %d for Grad %d, got %d", tt.expectedResistenzBody, tt.grad, result.ResistenzKoerper)
			}
			if result.ResistenzGeist != tt.expectedResistenzMind {
				t.Errorf("Expected ResistenzGeist %d for Grad %d, got %d", tt.expectedResistenzMind, tt.grad, result.ResistenzGeist)
			}
			if result.Zaubern != tt.expectedZaubern {
				t.Errorf("Expected Zaubern %d for Grad %d, got %d", tt.expectedZaubern, tt.grad, result.Zaubern)
			}
		})
	}
}
