package main

import (
	"log/slog"
	"os"
	"sync"
)

// Portion represents a data point from the feeding guide table.
type Portion struct {
	WeightKG float64 `json:"weight_kg"`
	DailyGr  float64 `json:"daily_gr"`
}

// A mutex to ensure thread-safe access to the feedingData slice.
var dataMutex sync.RWMutex

// Mutex for global coefficients
var coeffsMutex sync.RWMutex

// Global coefficients for the polynomial model
var globalCoeffs []float64

// Logger
var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
