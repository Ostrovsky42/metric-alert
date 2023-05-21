package validator

import (
	"metric-alert/internal/entities"
)

func ValidateUpdate(metric entities.Metrics) error {
	if metric.MType == "" {
		return errEmptyMetricType
	}

	if metric.ID == "" {
		return errEmptyMetricName
	}

	if metric.MType == entities.Counter && metric.Delta == nil {
		return errNotProvidedValue
	}

	if metric.MType == entities.Gauge && metric.Value == nil {
		return errNotProvidedValue
	}

	return nil
}

func ValidateGet(metric entities.Metrics) error {
	if metric.MType == "" {
		return errEmptyMetricType
	}

	if metric.MType != entities.Gauge && metric.MType != entities.Counter {
		return errUnknownMetricType
	}

	if metric.ID == "" {
		return errEmptyMetricName
	}

	return nil
}
