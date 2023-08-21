package storage

import (
	"sort"

	"metric-alert/internal/server/entities"
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

func RemoveDuplicatesIDs(ids []string) []string {
	uniqueIDs := make(map[string]struct{})
	result := make([]string, 0, len(ids))

	for _, item := range ids {
		if _, ok := uniqueIDs[item]; !ok {
			uniqueIDs[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
