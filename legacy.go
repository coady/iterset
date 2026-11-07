//go:build !go1.27

package iterset

import (
	"iter"
)

// Equal returns whether the key sets are equivalent.
//
// Related:
//   - [maps.Equal] to compare values
//   - [maps.EqualFunc] for two maps
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
	return len(m) == len(m.intersection(keys))
}

// IsDisjoint returns whether no keys are present.
// Use [maps.Keys] on the smaller of two maps.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint(keys iter.Seq[K]) bool {
	return len(m) == 0 || allFunc(keys, m.Missing)
}

// IntersectionWith keeps only the keys present in both.
// Receiver values are retained.
//
// Related:
//   - [MapSet.Intersection] to not modify in-place
//
// Performance:
//   - time: O(m+k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) IntersectionWith(keys iter.Seq[K]) {
	s := m.intersection(keys)
	if len(s) < len(m) {
		keep(m, s)
	}
}

// Intersection returns the ordered key-value pairs which are present in both.
// Use [maps.Keys] on the smaller of two maps.
//
// Related:
//   - [MapSet.IntersectionWith] to modify in-place
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) Intersection(keys iter.Seq[K]) iter.Seq2[K, V] {
	return m.mask(keys)
}

// IntersectionCount returns the number of keys present in both.
//
// Related:
//   - [MapSet.Overlap] for distinct count
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IntersectionCount(keys iter.Seq[K]) int {
	return m.intersectionCount(keys)
}

// Difference returns the key-value pairs which are not present in the keys.
//
// Related:
//   - [MapSet.DifferenceWith] to modify in-place
//   - [MapSet.ReverseDifference] if the keys were a map
//   - [Difference] if the receiver was not a map
//
// Performance:
//   - time:  O(m+k)
//   - space: O(min(m,k))
func (m MapSet[K, V]) Difference(keys iter.Seq[K]) iter.Seq2[K, V] {
	return m.difference(keys)
}

// DifferenceWith removes keys present in keys.
//
// Related:
//   - [MapSet.Difference] to not modify in-place
func (m MapSet[K, V]) DifferenceWith(keys iter.Seq[K]) {
	m.remove(keys)
}

// SymmetricDifference returns keys which are not in both.
//
// Related:
//   - [MapSet.Toggle] to modify in-place
//
// Performance:
//   - time: O(m+k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) SymmetricDifference(keys iter.Seq[K]) iter.Seq[K] {
	return m.symmetricDifference(keys)
}

// Overlap returns the sizes of the intersection and differences:
// left only, both, right only.
//
// Similarity measures:
//   - overlap coefficient: both / (min(left, right) + both)
//   - Jaccard index: both / (left + both + right)
//
// Related:
//   - [MapSet.IntersectionCount] for just the intersection size
//
// Performance:
//   - time: Θ(k)
//   - space: Θ(k)
func (m MapSet[K, V]) Overlap(keys iter.Seq[K]) (int, int, int) {
	return m.overlap(keys)
}
