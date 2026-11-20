package model

func ValuesToGroups[T any, K comparable](values []T, valueToKey func(T) K) [][]T {
	var groups [][]T
	indexByKey := map[K]int{}
	for _, value := range values {
		key := valueToKey(value)
		index, found := indexByKey[key]
		if !found {
			index = len(groups)
			indexByKey[key] = index
			groups = append(groups, nil)
		}
		groups[index] = append(groups[index], value)
	}
	return groups
}
