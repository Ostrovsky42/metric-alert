package validator

import (
	"strconv"

	"metric-alert/internal/server/entities"
)

func ValidateUpdateWithBody(metric entities.Metrics) error {
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

func ValidateMetrics(metrics []entities.Metrics) error {
	for _, metric := range metrics {
		if err := ValidateUpdateWithBody(metric); err != nil {
			return err
		}
	}

	return nil
}

func ValidateGetWithBody(metric entities.Metrics) error {
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

func ValidateUpdate(metric *entities.Metrics, mValue string) error {
	if metric.MType == "" {
		return errEmptyMetricType
	}

	if metric.MType == "" {
		return errEmptyMetricName
	}

	var err error
	switch metric.MType {
	case entities.Gauge:
		metric.Value, err = prepareGauge(mValue)
	case entities.Counter:
		metric.Delta, err = prepareCounter(mValue)
	default:
		err = errUnknownMetricType
	}
	if err != nil {
		return err
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

func prepareGauge(metric string) (*float64, error) {
	value, err := strconv.ParseFloat(metric, 64)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func prepareCounter(metric string) (*int64, error) {
	value, err := strconv.Atoi(metric)
	if err != nil {
		return nil, err
	}
	intVal := int64(value)
	return &intVal, nil
}
