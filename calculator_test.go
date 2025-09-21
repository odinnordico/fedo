package main

import (
	"testing"
)

func TestCalculatePortionSize(t *testing.T) {
	// Test with valid coefficients
	coeffs := []float64{1.0, 2.0, 3.0, 4.0}
	weight := 5.0
	expected := (1.0 + 2.0*5.0 + 3.0*25.0 + 4.0*125.0) / 2
	result := calculatePortionSize(weight, coeffs)
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}

	// Test with nil coeffs
	result = calculatePortionSize(weight, nil)
	if result != 0 {
		t.Errorf("Expected 0 for nil coeffs, got %f", result)
	}

	// Test with invalid coeffs length
	coeffs = []float64{1.0, 2.0}
	result = calculatePortionSize(weight, coeffs)
	if result != 0 {
		t.Errorf("Expected 0 for invalid coeffs length, got %f", result)
	}
}

func BenchmarkCalculatePortionSize(b *testing.B) {
	coeffs := []float64{1.0, 2.0, 3.0, 4.0}
	weight := 5.0
	for i := 0; i < b.N; i++ {
		calculatePortionSize(weight, coeffs)
	}
}
