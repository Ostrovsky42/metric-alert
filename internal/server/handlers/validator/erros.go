package validator

import "errors"

const (
	unknownMetricType = "unknown metric type"
	notProvidedValue  = "not provided metric value"
	emptyMetricType   = "empty metric type"
	EmptyMetricName   = "empty metric name"
)

var errUnknownMetricType = errors.New(unknownMetricType)
var errNotProvidedValue = errors.New(notProvidedValue)
var errEmptyMetricType = errors.New(emptyMetricType)
var errEmptyMetricName = errors.New(EmptyMetricName)
