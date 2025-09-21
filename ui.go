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

var editWindowOpen bool
var loadJSONWindowOpen bool
var currentLoadWindow fyne.Window

func showEditWindow(myApp fyne.App, feedingData *[]Portion) {
	if editWindowOpen {
		return
	}
	editWindowOpen = true
	editWindow := myApp.NewWindow("Edit Feeding Data")
	editWindow.Resize(fyne.NewSize(800, 600))

	var weightEntries []*widget.Entry
	var dailyEntries []*widget.Entry
	var content *fyne.Container
	var buildGrid func() *container.Scroll

	buildGrid = func() *container.Scroll {
		gridContainer := container.NewVBox()
		weightEntries = make([]*widget.Entry, len(*feedingData))
		dailyEntries = make([]*widget.Entry, len(*feedingData))

		for i, p := range *feedingData {
			localIndex := i

			weightEntry := widget.NewEntry()
			weightEntry.SetText(fmt.Sprintf("%.2f", p.WeightKG))
			weightEntries[localIndex] = weightEntry

			dailyEntry := widget.NewEntry()
			dailyEntry.SetText(fmt.Sprintf("%.2f", p.DailyGr))
			dailyEntries[localIndex] = dailyEntry

			deleteButton := widget.NewButton("Delete", func() {
				if len(*feedingData) > 3 {
					// Update feedingData with current entries
					for i := range *feedingData {
						if w, err := strconv.ParseFloat(weightEntries[i].Text, 64); err == nil {
							(*feedingData)[i].WeightKG = w
						}
						if d, err := strconv.ParseFloat(dailyEntries[i].Text, 64); err == nil {
							(*feedingData)[i].DailyGr = d
						}
					}
					// Remove the item
					*feedingData = append((*feedingData)[:localIndex], (*feedingData)[localIndex+1:]...)
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
		for i := range *feedingData {
			if w, err := strconv.ParseFloat(weightEntries[i].Text, 64); err == nil {
				(*feedingData)[i].WeightKG = w
			}
			if d, err := strconv.ParseFloat(dailyEntries[i].Text, 64); err == nil {
				(*feedingData)[i].DailyGr = d
			}
		}
		*feedingData = append(*feedingData, Portion{WeightKG: 0, DailyGr: 0})
		content.Objects[1] = buildGrid()
		content.Refresh()
	})

	loadJSONButton := widget.NewButton("Load JSON", func() {
		showLoadJSONWindow(myApp, feedingData, func() {
			content.Objects[1] = buildGrid()
			content.Refresh()
		})
	})

	saveButton := widget.NewButton("Save", func() {
		var newFeedingData []Portion
		hasError := false

		for i := range *feedingData {
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
			dataMutex.Lock()
			*feedingData = newFeedingData
			dataMutex.Unlock()

			err := saveFeedingData("feeding_data.json", *feedingData)
			if err != nil {
				widget.NewModalPopUp(widget.NewLabel("Error saving data"), editWindow.Canvas()).Show()
			}

			coeffsMutex.Lock()
			globalCoeffs = trainModel(*feedingData)
			coeffsMutex.Unlock()

			editWindow.Close()
			editWindowOpen = false
		}
	})

	cancelButton := widget.NewButton("Cancel", func() {
		if loadJSONWindowOpen {
			currentLoadWindow.Close()
			loadJSONWindowOpen = false
		}
		editWindow.Close()
		editWindowOpen = false
	})

	content = container.NewVBox(
		widget.NewLabel("Edit Feeding Data Points"),
		buildGrid(),
		container.NewGridWithColumns(4, addButton, loadJSONButton, saveButton, cancelButton),
	)
	editWindow.SetContent(content)
	editWindow.Show()
}

func showLoadJSONWindow(myApp fyne.App, feedingData *[]Portion, refresh func()) {
	if loadJSONWindowOpen {
		return
	}
	loadJSONWindowOpen = true
	currentLoadWindow = myApp.NewWindow("Load JSON Data")
	currentLoadWindow.Resize(fyne.NewSize(600, 400))
	loadWindow := currentLoadWindow

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
		*feedingData = newData
		refresh()
		loadWindow.Close()
		loadJSONWindowOpen = false
	})

	cancelButton := widget.NewButton("Cancel", func() {
		loadWindow.Close()
		loadJSONWindowOpen = false
	})

	loadWindow.SetContent(container.NewVBox(
		widget.NewLabel("Load Feeding Data from JSON"),
		jsonEntry,
		container.NewGridWithColumns(2, loadButton, cancelButton),
	))
	loadWindow.Show()
}

func createMainWindow(myApp fyne.App, feedingData []Portion) fyne.Window {
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
		dataMutex.RLock()
		portionSize := calculatePortionSize(weight, feedingData)
		dataMutex.RUnlock()

		roundedPortionSize := math.Round(portionSize*100) / 100

		// Update the result label
		resultLabel.SetText(fmt.Sprintf("For a dog weighing %.2f kg, the portion size is %.2f grams.", weight, roundedPortionSize))
	})

	// Create the edit button that opens a new window
	editButton := widget.NewButton("Edit Data", func() {
		showEditWindow(myApp, &feedingData)
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
