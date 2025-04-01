package utils

import (
	"strings"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// Generate a UUID
	uuid := GenerateUUID()

	// Check that it's not empty
	if uuid == "" {
		t.Error("GenerateUUID returned an empty string")
	}

	// Check that it has the correct format (8-4-4-4-12 characters)
	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		t.Errorf("UUID has incorrect format, expected 5 parts separated by hyphens, got %d parts", len(parts))
	}

	// Check the length of each part
	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			t.Errorf("Part %d of UUID has incorrect length, expected %d, got %d", i+1, expectedLengths[i], len(part))
		}
	}
}

func TestRoundDown(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		places int
		want   float64
	}{
		{"Round down to 2 decimal places", 3.14159, 2, 3.14},
		{"Round down to 0 decimal places", 3.99, 0, 3.0},
		{"Round down negative number", -3.14159, 2, -3.15},
		{"Round down to 4 decimal places", 1.23456789, 4, 1.2345},
		{"Round down already rounded number", 3.14, 2, 3.14},
		{"Round down zero", 0.0, 2, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundDown(tt.value, tt.places)
			if got != tt.want {
				t.Errorf("RoundDown(%f, %d) = %f, want %f", tt.value, tt.places, got, tt.want)
			}
		})
	}
}

func TestValidateNAB(t *testing.T) {
	tests := []struct {
		name         string
		totalBalance float64
		totalUnits   float64
		want         float64
	}{
		{"Normal case", 1000.0, 100.0, 10.0},
		{"Zero units", 1000.0, 0.0, 1.0},
		{"Negative units", 1000.0, -10.0, 1.0},
		{"Zero balance", 0.0, 100.0, 0.0},
		{"Round down to 4 decimal places", 1000.0, 3.0, 333.3333},
		{"Both zero", 0.0, 0.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateNAB(tt.totalBalance, tt.totalUnits)
			if got != tt.want {
				t.Errorf("ValidateNAB(%f, %f) = %f, want %f", tt.totalBalance, tt.totalUnits, got, tt.want)
			}
		})
	}
}
