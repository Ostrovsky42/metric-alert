package handlers

import (
	"metric-alert/internal/types"
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

func parseURL(req []string) (types.Metric, error) {
	if len(req) < 5 {
		return types.Metric{}, ErrEmptyMetric

	}
	if req[metricType] != types.Gauge && req[metricType] != types.Counter {
		return types.Metric{}, ErrUnknownMetric
	}

	metric := types.Metric{
		MetricName: req[metricName],
		MetricType: req[metricType],
	}

	var err error
	switch metric.MetricType {
	case types.Gauge:
		metric.GaugeValue, err = prepareGauge(req[metricValue])
	case types.Counter:
		metric.CounterValue, err = prepareCounter(req[metricValue])
	default:
		err = ErrMetricType
	}
	if err != nil {
		return types.Metric{}, err
	}

	return metric, nil
}
