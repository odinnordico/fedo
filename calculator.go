package main

// calculatePortionSize finds the correct portion size for a given weight.
// It uses the trained polynomial regression model to predict.
func calculatePortionSize(weight float64, feedingData []Portion) float64 {
	coeffsMutex.RLock()
	defer coeffsMutex.RUnlock()
	if globalCoeffs == nil || len(globalCoeffs) != 4 {
		logger.Warn("Model not trained or invalid coefficients")
		return 0
	}

	// Predict for the weight
	prediction := globalCoeffs[0] + globalCoeffs[1]*weight + globalCoeffs[2]*weight*weight + globalCoeffs[3]*weight*weight*weight

	logger.Debug("Calculated portion size", "weight", weight, "prediction", prediction)

	return prediction / 2
}
