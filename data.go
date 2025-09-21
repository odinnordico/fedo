package main

import (
	"encoding/json"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gonum.org/v1/gonum/mat"
)

func loadFeedingData(filename string) ([]Portion, error) {
	app := fyne.CurrentApp()
	rootURI := app.Storage().RootURI()
	fileURI, err := storage.ParseURI(rootURI.String() + "/" + filename)
	initialData := []Portion{
		{
			WeightKG: 1,
			DailyGr:  30,
		},
		{
			WeightKG: 2.25,
			DailyGr:  50,
		},
		{
			WeightKG: 4.5,
			DailyGr:  100,
		},
		{
			WeightKG: 9,
			DailyGr:  150,
		},
		{
			WeightKG: 13.5,
			DailyGr:  200,
		},
		{
			WeightKG: 18,
			DailyGr:  233,
		},
		{
			WeightKG: 27,
			DailyGr:  325,
		},
		{
			WeightKG: 36,
			DailyGr:  375,
		},
		{
			WeightKG: 45,
			DailyGr:  450,
		},
		{
			WeightKG: 57,
			DailyGr:  525,
		},
		{
			WeightKG: 68,
			DailyGr:  600,
		},
		{
			WeightKG: 79,
			DailyGr:  666,
		},
	}
	if err != nil {
		logger.Error("Failed to parse URI", "error", err)
		return initialData, err
	}

	// Check if file exists
	exists, err := storage.Exists(fileURI)
	if err != nil || !exists {
		// File does not exist, create with initial data
		err := saveFeedingData(filename, initialData)
		if err != nil {
			logger.Error("Failed to create initial feeding data file", "error", err)
			return nil, err
		}
		logger.Info("Created initial feeding data file", "count", len(initialData))
		return initialData, nil
	}

	reader, err := storage.Reader(fileURI)
	if err != nil {
		logger.Error("Failed to read feeding data file", "error", err)
		return nil, err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		logger.Error("Failed to read feeding data", "error", err)
		return nil, err
	}
	var feedingData []Portion
	err = json.Unmarshal(data, &feedingData)
	if err != nil {
		logger.Error("Failed to unmarshal feeding data", "error", err)
		return nil, err
	}
	logger.Info("Loaded feeding data", "count", len(feedingData))
	return feedingData, nil
}

func saveFeedingData(filename string, feedingData []Portion) error {
	data, err := json.MarshalIndent(feedingData, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal feeding data", "error", err)
		return err
	}
	app := fyne.CurrentApp()
	rootURI := app.Storage().RootURI()
	fileURI, err := storage.ParseURI(rootURI.String() + "/" + filename)
	if err != nil {
		logger.Error("Failed to parse URI", "error", err)
		return err
	}
	writer, err := storage.Writer(fileURI)
	if err != nil {
		logger.Error("Failed to write feeding data file", "error", err)
		return err
	}
	defer writer.Close()
	_, err = writer.Write(data)
	if err != nil {
		logger.Error("Failed to write feeding data", "error", err)
		return err
	}
	logger.Info("Saved feeding data", "count", len(feedingData))
	return nil
}

func trainModel(feedingData []Portion) []float64 {
	if len(feedingData) == 0 {
		logger.Warn("No feeding data to train model")
		return nil
	}

	x := make([]float64, len(feedingData))
	y := make([]float64, len(feedingData))
	for i, p := range feedingData {
		x[i] = p.WeightKG
		y[i] = p.DailyGr
	}

	// Fit polynomial of degree 3
	n := len(x)
	v := mat.NewDense(n, 4, nil)
	for i := 0; i < n; i++ {
		v.Set(i, 0, 1)
		v.Set(i, 1, x[i])
		v.Set(i, 2, x[i]*x[i])
		v.Set(i, 3, x[i]*x[i]*x[i])
	}
	yVec := mat.NewVecDense(n, y)
	var coeffsVec mat.VecDense
	err := coeffsVec.SolveVec(v, yVec)
	if err != nil {
		logger.Error("Failed to train model", "error", err)
		return nil
	}
	coeffs := coeffsVec.RawVector().Data
	logger.Info("Trained polynomial regression model", "degree", 3, "data_points", n)
	return coeffs
}
