package handlers

import (
	"metric-alert/internal/types"
)

func ValidateUpdate(metric *types.Metric, mValue string) error {
	if metric.MetricType == "" {
		return ErrEmptyMetricType
	}

	if metric.MetricName == "" {
		return ErrEmptyMetricName
	}

	var err error
	switch metric.MetricType {
	case types.Gauge:
		metric.GaugeValue, err = prepareGauge(mValue)
	case types.Counter:
		metric.CounterValue, err = prepareCounter(mValue)
	default:
		err = ErrUnknownMetricType
	}
	if err != nil {
		return err
	}

	return nil
}

func ValidateGet(metric types.Metric) error {
	if metric.MetricType == "" {
		return ErrEmptyMetricType
	}

	if metric.MetricType != types.Gauge && metric.MetricType != types.Counter {
		return ErrUnknownMetricType
	}

	if metric.MetricName == "" {
		return ErrEmptyMetricName
	}

	return nil
}
