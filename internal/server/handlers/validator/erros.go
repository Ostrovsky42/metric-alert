package validator

import "errors"

const (
	UnknownMetricType = "unknown metric type"
	NotProvidedValue  = "not provided metric value"
	EmptyMetricType   = "empty metric type"
	EmptyMetricName   = "empty metric name"
)

var errUnknownMetricType = errors.New(UnknownMetricType)
var errNotProvidedValue = errors.New(NotProvidedValue)
var errEmptyMetricType = errors.New(EmptyMetricType)
var errEmptyMetricName = errors.New(EmptyMetricName)
