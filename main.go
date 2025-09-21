package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	logger.Info("Starting Dog Food Calculator")

	// Load the feeding data from JSON file
	feedingData, err := loadFeedingData("feeding_data.json")
	if err != nil {
		logger.Error("Failed to load feeding data, using empty data", "error", err)
		feedingData = []Portion{}
	}

	// Train the polynomial regression model
	coeffsMutex.Lock()
	globalCoeffs = trainModel(feedingData)
	coeffsMutex.Unlock()

	// Create a new fyne application with a unique ID
	myApp := app.NewWithID("com.github.odinnordico.fedo")
	icon, err := fyne.LoadResourceFromPath("fedo.png")
	if err == nil {
		myApp.SetIcon(icon)
	}

	myWindow := createMainWindow(myApp, feedingData)
	myWindow.ShowAndRun()
}
