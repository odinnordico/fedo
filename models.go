package main

import (
	"log/slog"
	"os"
	"sync"

	"fyne.io/fyne/v2"
)

// Portion represents a data point from the feeding guide table.
type Portion struct {
	WeightKG float64 `json:"weight_kg"`
	DailyGr  float64 `json:"daily_gr"`
}

// App holds the application state.
type App struct {
	FeedingData        []Portion
	Coeffs             []float64
	DataMutex          sync.RWMutex
	CoeffsMutex        sync.RWMutex
	EditWindowOpen     bool
	LoadJSONWindowOpen bool
	CurrentLoadWindow  fyne.Window
}

// Logger
var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
