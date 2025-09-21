package main

import (
	"encoding/json"
	"os"

	"gonum.org/v1/gonum/mat"
)

func loadFeedingData(filename string) ([]Portion, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, create with initial data
		initialData := []Portion{
			{WeightKG: 1, DailyGr: 30},
			{WeightKG: 5, DailyGr: 100},
			{WeightKG: 10, DailyGr: 150},
		}
		err := saveFeedingData(filename, initialData)
		if err != nil {
			logger.Error("Failed to create initial feeding data file", "error", err)
			return nil, err
		}
		logger.Info("Created initial feeding data file", "count", len(initialData))
		return initialData, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("Failed to read feeding data file", "error", err)
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
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		logger.Error("Failed to write feeding data file", "error", err)
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
