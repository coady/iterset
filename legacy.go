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
//   - [maps.EqualFunc] for two maps
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) Equal(keys iter.Seq[K]) bool {
	return m.equal(keys)
}

// SubsetOf returns whether every map key is present in keys.
//
// Related:
//   - [MapSet.ContainsAll] for the other polarity (better when the argument is a sequence)
//   - [SubsetOf] if the receiver was not a map
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) SubsetOf(keys iter.Seq[K]) bool {
	return len(m) == len(m.intersection(keys))
}

// Intersects reports whether any key is present.
// Use [maps.Keys] on the smaller of two maps.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) Intersects(keys iter.Seq[K]) bool {
	return len(m) > 0 && !allFunc(keys, m.Missing)
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
		intersectionWith(m, s)
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
//   - [MapSet.SymmetricDifferenceWith] to modify in-place
//
// Performance:
//   - time: O(m+k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) SymmetricDifference(keys iter.Seq[K]) iter.Seq[K] {
	return m.symmetricDifference(keys)
}

// Union merges keys with inserts.
// Duplicate keys overwrite values.
//
// Related:
//   - [MapSet.UnionWith] to modify in-place
//   - [SortedUnion] for sorted sequences
//
// Performance:
//   - time: Θ(m+k)
//   - space: Ω(max(m, k))..O(m+k)
func (m MapSet[K, V]) Union(seq iter.Seq2[K, V]) MapSet[K, V] {
	m = maps.Clone(m)
	if m == nil {
		m = map[K]V{}
	}
	maps.Insert(m, seq)
	return m
}

// UnionWith merges keys with inserts.
// Duplicate keys overwrite values.
//
// Related:
//   - [MapSet.Union] to not modify in-place
//   - [maps.Insert] for the same operation
func (m MapSet[K, V]) UnionWith(seq iter.Seq2[K, V]) {
	maps.Insert(m, seq)
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
