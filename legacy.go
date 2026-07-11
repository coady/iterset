//go:build !go1.27

package iterset

import (
	"iter"
	"maps"
)

// Equal returns whether the key sets are equivalent.
//
// Related:
//   - [maps.Equal] to compare values
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) Equal(keys iter.Seq[K]) bool {
	return m.equal(keys)
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
func (m MapSet[K, V]) IsSubset(keys iter.Seq[K]) bool {
	return len(m) == len(m.intersect(keys))
}

// IsDisjoint returns whether no keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint(keys iter.Seq[K]) bool {
	return len(m) == 0 || allFunc(keys, m.Missing)
}

// Intersect returns the ordered key-value pairs which are present in both.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) Intersect(keys iter.Seq[K]) iter.Seq2[K, V] {
	if len(m) == 0 {
		return maps.All(m)
	}
	return m.mask(keys)
}

// Difference returns the key-value pairs which are not present in the keys.
//
// Related:
//   - [MapSet.Remove] to modify in-place
//   - [MapSet.ReverseDifference] if the keys were a map
//   - [Difference] if the receiver was not a map
//
// Performance:
//   - time:  O(m+k)
//   - space: O(min(m,k))
func (m MapSet[K, V]) Difference(keys iter.Seq[K]) iter.Seq2[K, V] {
	s := m.intersect(keys)
	if len(m) == len(s) {
		return func(func(K, V) bool) {}
	}
	return m.filter(s.Missing)
}
