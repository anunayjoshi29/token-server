package routecalc

func CalculateAmountOut(amountIn, reserveIn, reserveOut float64) (float64, error) {
	// Basic validation
	if reserveIn <= 0 || reserveOut <= 0 || amountIn <= 0 {
		return 0, nil
	}

	newReserveIn := reserveIn + amountIn
	// newReserveIn * newReserveOut = reserveIn * reserveOut  => newReserveOut = (reserveIn*reserveOut)/newReserveIn
	newReserveOut := (reserveIn * reserveOut) / newReserveIn
	amountOut := reserveOut - newReserveOut
	return amountOut, nil
}
