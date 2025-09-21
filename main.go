// Package main implements a dog food calculator application using polynomial regression.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// main is the entry point of the application.

func main() {
	logger.Info("Starting Dog Food Calculator")

	// Create a new fyne application with a unique ID
	myApp := app.NewWithID("com.github.odinnordico.fedo")
	icon, err := fyne.LoadResourceFromPath("fedo.png")
	if err == nil {
		myApp.SetIcon(icon)
	}

	// Load the feeding data from JSON file
	feedingData, err := loadFeedingData("feeding_data.json")
	if err != nil {
		logger.Error("Failed to load feeding data, using empty data", "error", err)
		feedingData = []Portion{}
	}

	// Create app state
	appState := &App{
		FeedingData: feedingData,
	}

	// Train the polynomial regression model
	appState.CoeffsMutex.Lock()
	appState.Coeffs = trainModel(feedingData)
	appState.CoeffsMutex.Unlock()

	myWindow := createMainWindow(myApp, appState)
	myWindow.ShowAndRun()
}
