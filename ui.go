package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Window state variables are now in App struct

// showEditWindow displays a window for editing the feeding data points.
func showEditWindow(myApp fyne.App, app *App) {
	if app.EditWindowOpen {
		return
	}
	app.EditWindowOpen = true
	editWindow := myApp.NewWindow("Edit Feeding Data")
	editWindow.Resize(fyne.NewSize(800, 600))

	var weightEntries []*widget.Entry
	var dailyEntries []*widget.Entry
	var content *fyne.Container
	var buildGrid func() *container.Scroll

	buildGrid = func() *container.Scroll {
		gridContainer := container.NewVBox()
		weightEntries = make([]*widget.Entry, len(app.FeedingData))
		dailyEntries = make([]*widget.Entry, len(app.FeedingData))

		for i, p := range app.FeedingData {
			localIndex := i

			weightEntry := widget.NewEntry()
			weightEntry.SetText(fmt.Sprintf("%.2f", p.WeightKG))
			weightEntries[localIndex] = weightEntry

			dailyEntry := widget.NewEntry()
			dailyEntry.SetText(fmt.Sprintf("%.2f", p.DailyGr))
			dailyEntries[localIndex] = dailyEntry

			deleteButton := widget.NewButton("Delete", func() {
				if len(app.FeedingData) > 3 {
					// Update feedingData with current entries
					for i := range app.FeedingData {
						if w, err := strconv.ParseFloat(weightEntries[i].Text, 64); err == nil {
							app.FeedingData[i].WeightKG = w
						}
						if d, err := strconv.ParseFloat(dailyEntries[i].Text, 64); err == nil {
							app.FeedingData[i].DailyGr = d
						}
					}
					// Remove the item
					app.FeedingData = append(app.FeedingData[:localIndex], app.FeedingData[localIndex+1:]...)
					// Rebuild the grid
					content.Objects[1] = buildGrid()
					content.Refresh()
				} else {
					label := widget.NewLabel("Cannot delete: minimum 3 records required")
					closeButton := widget.NewButton("OK", func() {})
					popup := widget.NewModalPopUp(container.NewVBox(label, closeButton), editWindow.Canvas())
					closeButton.OnTapped = func() { popup.Hide() }
					popup.Show()
				}
			})

			row := container.NewGridWithColumns(5,
				widget.NewLabel("Weight (kg):"),
				weightEntry,
				widget.NewLabel("Daily (g):"),
				dailyEntry,
				deleteButton,
			)
			gridContainer.Add(row)
		}

		editGridScroll := container.NewScroll(gridContainer)
		editGridScroll.SetMinSize(fyne.NewSize(780, 500))
		return editGridScroll
	}

	addButton := widget.NewButton("Add Row", func() {
		// Update feedingData with current entries
		for i := range app.FeedingData {
			if w, err := strconv.ParseFloat(weightEntries[i].Text, 64); err == nil {
				app.FeedingData[i].WeightKG = w
			}
			if d, err := strconv.ParseFloat(dailyEntries[i].Text, 64); err == nil {
				app.FeedingData[i].DailyGr = d
			}
		}
		app.FeedingData = append(app.FeedingData, Portion{WeightKG: 0, DailyGr: 0})
		content.Objects[1] = buildGrid()
		content.Refresh()
	})

	loadJSONButton := widget.NewButton("Load JSON", func() {
		showLoadJSONWindow(myApp, app, func() {
			content.Objects[1] = buildGrid()
			content.Refresh()
		})
	})

	saveButton := widget.NewButton("Save", func() {
		var newFeedingData []Portion
		hasError := false

		for i := range app.FeedingData {
			weight, err := strconv.ParseFloat(weightEntries[i].Text, 64)
			if err != nil {
				widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Invalid weight value: %s", weightEntries[i].Text)), editWindow.Canvas()).Show()
				hasError = true
				break
			}
			daily, err := strconv.ParseFloat(dailyEntries[i].Text, 64)
			if err != nil {
				widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Invalid daily grams value: %s", dailyEntries[i].Text)), editWindow.Canvas()).Show()
				hasError = true
				break
			}
			newFeedingData = append(newFeedingData, Portion{WeightKG: weight, DailyGr: daily})
		}

		if !hasError {
			app.DataMutex.Lock()
			app.FeedingData = newFeedingData
			app.DataMutex.Unlock()

			err := saveFeedingData("feeding_data.json", app.FeedingData)
			if err != nil {
				widget.NewModalPopUp(widget.NewLabel("Error saving data"), editWindow.Canvas()).Show()
			}

			app.CoeffsMutex.Lock()
			app.Coeffs = trainModel(app.FeedingData)
			app.CoeffsMutex.Unlock()

			editWindow.Close()
			app.EditWindowOpen = false
		}
	})

	cancelButton := widget.NewButton("Cancel", func() {
		if app.LoadJSONWindowOpen {
			app.CurrentLoadWindow.Close()
			app.LoadJSONWindowOpen = false
		}
		editWindow.Close()
		app.EditWindowOpen = false
	})

	content = container.NewVBox(
		widget.NewLabel("Edit Feeding Data Points"),
		buildGrid(),
		container.NewGridWithColumns(4, addButton, loadJSONButton, saveButton, cancelButton),
	)
	editWindow.SetContent(content)
	editWindow.Show()
}

// showLoadJSONWindow displays a window for loading feeding data from JSON.
func showLoadJSONWindow(myApp fyne.App, app *App, refresh func()) {
	if app.LoadJSONWindowOpen {
		return
	}
	app.LoadJSONWindowOpen = true
	app.CurrentLoadWindow = myApp.NewWindow("Load JSON Data")
	app.CurrentLoadWindow.Resize(fyne.NewSize(600, 400))
	loadWindow := app.CurrentLoadWindow

	jsonEntry := widget.NewMultiLineEntry()
	jsonEntry.SetPlaceHolder("Paste JSON data here...")

	loadButton := widget.NewButton("Load", func() {
		var newData []Portion
		err := json.Unmarshal([]byte(jsonEntry.Text), &newData)
		if err != nil {
			label := widget.NewLabel(fmt.Sprintf("JSON parsing error: %s", err.Error()))
			closeButton := widget.NewButton("OK", func() {})
			popup := widget.NewModalPopUp(container.NewVBox(label, closeButton), loadWindow.Canvas())
			closeButton.OnTapped = func() { popup.Hide() }
			popup.Show()
			return
		}
		app.FeedingData = newData
		refresh()
		loadWindow.Close()
		app.LoadJSONWindowOpen = false
	})

	cancelButton := widget.NewButton("Cancel", func() {
		loadWindow.Close()
		app.LoadJSONWindowOpen = false
	})

	loadWindow.SetContent(container.NewVBox(
		widget.NewLabel("Load Feeding Data from JSON"),
		jsonEntry,
		container.NewGridWithColumns(2, loadButton, cancelButton),
	))
	loadWindow.Show()
}

// createMainWindow creates and returns the main application window.
func createMainWindow(myApp fyne.App, app *App) fyne.Window {
	myWindow := myApp.NewWindow("Dog Food Calculator")
	myWindow.Resize(fyne.NewSize(400, 200))

	// Create UI widgets
	weightInput := widget.NewEntry()
	weightInput.SetPlaceHolder("Enter dog's weight in kg...")
	resultLabel := widget.NewLabel("Enter a weight to calculate.")
	resultLabel.Wrapping = fyne.TextWrapWord

	// Create the calculation button
	calculateButton := widget.NewButton("Calculate", func() {
		weightStr := weightInput.Text
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Invalid weight: %s", weightStr))
			return
		}
		if weight <= 0 || weight >= 100 {
			resultLabel.SetText("Weight must be greater than 0 and less than 100 kg.")
			return
		}

		// Use a mutex to read the feeding data safely
		app.DataMutex.RLock()
		portionSize := calculatePortionSize(weight, app.Coeffs)
		app.DataMutex.RUnlock()

		roundedPortionSize := math.Round(portionSize*100) / 100

		// Update the result label
		resultLabel.SetText(fmt.Sprintf("For a dog weighing %.2f kg, the portion size is %.2f grams.", weight, roundedPortionSize))
	})

	// Create the edit button that opens a new window
	editButton := widget.NewButton("Edit Data", func() {
		showEditWindow(myApp, app)
	})

	// Set up the UI layout
	content := container.NewVBox(
		widget.NewLabel("Dog Food Portion Calculator"),
		weightInput,
		container.NewGridWithColumns(2, calculateButton, editButton),
		resultLabel,
	)

	myWindow.SetContent(content)
	return myWindow
}
