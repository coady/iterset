//go:build go1.27

package iterset

import (
	"iter"
	"maps"
	"slices"
)

// KeyInput is the set of key sources accepted by generic [MapSet] methods.
//
// May be an iterator, slice, or map. Slices and maps use better algorithms where applicable.
type KeyInput[K comparable, V any] interface {
	iter.Seq[K] | []K | map[K]V | MapSet[K, V]
}

// IsSubset returns whether every map key is present in keys.
//
// Related:
//   - [MapSet.IsSuperset] if the keys were a map
//   - [IsSubset] if the receiver was not a map
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) IsSubset[S KeyInput[K, V]](keys S) bool {
	var it iter.Seq[K]
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = keys
	case []K:
		return len(m) <= len(keys) && len(m) == len(m.intersect(slices.Values(keys)))
	case map[K]V:
		return len(m) <= len(keys) && Cast(keys).IsSuperset(maps.Keys(m))
	case MapSet[K, V]:
		return len(m) <= len(keys) && keys.IsSuperset(maps.Keys(m))
	}
	return len(m) == len(m.intersect(it))
}
