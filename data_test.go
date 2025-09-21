package main

import (
	"testing"
)

func TestTrainModel(t *testing.T) {
	// Test with sample data
	feedingData := []Portion{
		{WeightKG: 1, DailyGr: 30},
		{WeightKG: 2, DailyGr: 50},
		{WeightKG: 3, DailyGr: 70},
	}
	coeffs := trainModel(feedingData)
	if coeffs == nil {
		t.Error("Expected coefficients, got nil")
	}
	if len(coeffs) != 4 {
		t.Errorf("Expected 4 coefficients, got %d", len(coeffs))
	}

	// Test with empty data
	coeffs = trainModel([]Portion{})
	if coeffs != nil {
		t.Error("Expected nil for empty data")
	}
}

func BenchmarkTrainModel(b *testing.B) {
	feedingData := []Portion{
		{WeightKG: 1, DailyGr: 30},
		{WeightKG: 2, DailyGr: 50},
		{WeightKG: 3, DailyGr: 70},
		{WeightKG: 4, DailyGr: 90},
		{WeightKG: 5, DailyGr: 110},
	}
	for i := 0; i < b.N; i++ {
		trainModel(feedingData)
	}
}
