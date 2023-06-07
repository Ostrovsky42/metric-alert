package storage

import (
	"metric-alert/internal/entities"
	"sort"
)

func sortMetric(metrics []entities.Metrics) {
	sortFunc := func(i, j int) bool {
		if metrics[i].MType == entities.Counter {
			return false
		}

		if metrics[j].MType == entities.Counter {
			return true
		}

		return metrics[i].ID < metrics[j].ID
	}

	sort.Slice(metrics, sortFunc)
}
