package main

// calculatePortionSize calculates the portion size for the given weight using the trained polynomial model.
// The result is divided by 2.
func calculatePortionSize(weight float64, coeffs []float64) float64 {
	if coeffs == nil || len(coeffs) != 4 {
		logger.Warn("Model not trained or invalid coefficients")
		return 0
	}

	// Predict for the weight
	prediction := coeffs[0] + coeffs[1]*weight + coeffs[2]*weight*weight + coeffs[3]*weight*weight*weight

	logger.Debug("Calculated portion size", "weight", weight, "prediction", prediction)

	return prediction / 2
}
