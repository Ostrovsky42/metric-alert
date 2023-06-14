package storage

import (
	"metric-alert/internal/entities"
	"sort"
)

const NotFound = "not found metric"

func SortMetric(metrics []entities.Metrics) {
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

func RemoveDuplicatesIDs(IDs []string) []string {
	uniqueIDs := make(map[string]struct{})
	result := make([]string, 0, len(IDs))

	for _, item := range IDs {
		if _, ok := uniqueIDs[item]; !ok {
			uniqueIDs[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
