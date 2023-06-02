package validator

import "errors"

var errUnknownMetricType = errors.New("unknown metric type")
var errNotProvidedValue = errors.New("not provided metric value")
var errEmptyMetricType = errors.New("empty metric type")
var errEmptyMetricName = errors.New("empty metric name")
