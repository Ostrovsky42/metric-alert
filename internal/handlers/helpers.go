package handlers

import (
	"strconv"
)

func prepareGauge(metric string) (float64, error) {
	value, err := strconv.ParseFloat(metric, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func prepareCounter(metric string) (int64, error) {
	value, err := strconv.Atoi(metric)
	if err != nil {
		return 0, err
	}
	return int64(value), nil
}
