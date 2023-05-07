package handlers

import "strconv"

func prepareGauge(metric []string) (string, float64, error) {
	name := metric[metricName]
	value, err := strconv.ParseFloat(metric[metricValue], 64)
	if err != nil {
		return "", 0, err
	}
	return name, value, nil
}

func prepareCounter(metric []string) (string, int64, error) {
	name := metric[metricName]
	value, err := strconv.Atoi(metric[metricValue])
	if err != nil {
		return "", 0, err
	}
	return name, int64(value), nil
}
