package validator

import (
	"strconv"

	"metric-alert/internal/types"
)

func ValidateUpdate(metric *types.Metric, mValue string) error {
	if metric.MetricType == "" {
		return errEmptyMetricType
	}

	if metric.MetricName == "" {
		return errEmptyMetricName
	}

	var err error
	switch metric.MetricType {
	case types.Gauge:
		metric.GaugeValue, err = prepareGauge(mValue)
	case types.Counter:
		metric.CounterValue, err = prepareCounter(mValue)
	default:
		err = errUnknownMetricType
	}
	if err != nil {
		return err
	}

	return nil
}

func ValidateGet(metric types.Metric) error {
	if metric.MetricType == "" {
		return errEmptyMetricType
	}

	if metric.MetricType != types.Gauge && metric.MetricType != types.Counter {
		return errUnknownMetricType
	}

	if metric.MetricName == "" {
		return errEmptyMetricName
	}

	return nil
}

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
