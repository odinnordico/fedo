package main

import (
	"testing"
)

func TestApp(t *testing.T) {
	app := &App{
		FeedingData: []Portion{
			{WeightKG: 1, DailyGr: 30},
		},
		Coeffs: []float64{1, 2, 3, 4},
	}
	if len(app.FeedingData) != 1 {
		t.Errorf("Expected 1 feeding data, got %d", len(app.FeedingData))
	}
	if len(app.Coeffs) != 4 {
		t.Errorf("Expected 4 coeffs, got %d", len(app.Coeffs))
	}
	if app.EditWindowOpen {
		t.Error("Expected EditWindowOpen false")
	}
}

func TestPortion(t *testing.T) {
	p := Portion{WeightKG: 5.5, DailyGr: 100.0}
	if p.WeightKG != 5.5 {
		t.Errorf("Expected WeightKG 5.5, got %f", p.WeightKG)
	}
	if p.DailyGr != 100.0 {
		t.Errorf("Expected DailyGr 100.0, got %f", p.DailyGr)
	}
}
