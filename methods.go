//go:build go1.27

package iterset

import (
	"iter"
	"maps"
	"slices"
)

// Equal returns whether the key sets are equivalent.
//
// Related:
//   - [maps.Equal] to compare values
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) Equal[S iter.Seq[K] | []K | MapSet[K, V]](keys S) bool {
	var it iter.Seq[K]
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = keys
	case []K:
		if len(m) > len(keys) {
			return false
		}
		it = slices.Values(keys)
	case MapSet[K, V]:
		return maps.EqualFunc(m, keys, func(V, V) bool { return true })
	}
	return m.equal(it)
}

// IsSubset returns whether every map key is present in keys.
//
// Related:
//   - [IsSubset] if the receiver was not a map
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) IsSubset[S iter.Seq[K] | []K | MapSet[K, V]](keys S) bool {
	var it iter.Seq[K]
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = keys
	case []K:
		if len(m) > len(keys) {
			return false
		}
		it = slices.Values(keys)
	case MapSet[K, V]:
		return len(m) <= len(keys) && keys.IsSuperset(maps.Keys(m))
	}
	return len(m) == len(m.intersect(it))
}

// IsDisjoint returns whether no keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint[S iter.Seq[K] | MapSet[K, V]](keys S) bool {
	var it iter.Seq[K]
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = keys
	case MapSet[K, V]:
		it = maps.Keys(keys)
		if len(m) < len(keys) {
			m, it = keys, maps.Keys(m)
		}
	}
	return len(m) == 0 || allFunc(it, m.Missing)
}

// Intersect returns the ordered key-value pairs which are present in both.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) Intersect[S iter.Seq[K] | MapSet[K, V]](keys S) iter.Seq2[K, V] {
	var it iter.Seq[K]
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = keys
	case MapSet[K, V]:
		if len(m) < len(keys) {
			return m.filter(keys.Contains)
		}
		it = maps.Keys(keys)
	}
	return m.mask(it)
}

// Difference returns the key-value pairs which are not present in the keys.
//
// Related:
//   - [MapSet.Remove] to modify in-place
//   - [Difference] if the receiver was not a map
//
// Performance:
//   - time:  O(m+k)
//   - space: O(min(m,k))
func (m MapSet[K, V]) Difference[S iter.Seq[K] | MapSet[K, V]](keys S) (it iter.Seq2[K, V]) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = m.difference(keys)
	case MapSet[K, V]:
		if len(keys) == 0 {
			return maps.All(m)
		}
		it = m.filter(keys.Missing)
	}
	return it
}
