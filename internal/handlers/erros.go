package handlers

import "errors"

var ErrUnknownMetric = errors.New("unknown metric type")
var ErrMetricType = errors.New("invalid metric type")
var ErrEmptyMetric = errors.New("empty metric name")
