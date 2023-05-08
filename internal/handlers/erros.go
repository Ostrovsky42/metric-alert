package handlers

import "errors"

var ErrUnknownMetricType = errors.New("unknown metric type")
var ErrEmptyMetricType = errors.New("empty metric type")
var ErrEmptyMetricName = errors.New("empty metric name")
