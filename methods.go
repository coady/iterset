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
//   - [maps.EqualFunc] for two maps with different value types
//
// Performance:
//   - time: O(k) if seq or slice
//   - time: O(min(m, k)) if map
//   - space: O(min(m, k)) if seq or slice
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

// SubsetOf returns whether every map key is present in keys.
//
// Related:
//   - [MapSet.ContainsAll] for the other polarity (better when the argument is a sequence)
//   - [SubsetOf] if the receiver was not a map
//
// Performance:
//   - time: O(k) if seq or slice
//   - time: O(min(m, k)) if map
//   - space: O(min(m, k)) if seq or slice
func (m MapSet[K, V]) SubsetOf[S iter.Seq[K] | []K | MapSet[K, V]](keys S) bool {
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
		return len(m) <= len(keys) && keys.ContainsAll(maps.Keys(m))
	}
	return len(m) == len(m.intersection(it))
}

// Intersects reports whether any key is present.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) Intersects[S iter.Seq[K] | MapSet[K, V]](keys S) bool {
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
	return len(m) > 0 && !allFunc(it, m.Missing)
}

// IntersectionWith keeps only the keys present in both.
// Receiver values are retained.
//
// Related:
//   - [MapSet.Intersection] to not modify in-place
//
// Performance:
//   - time: O(m+k) if seq
//   - time: O(m) if map
//   - space: O(min(m, k)) if seq
func (m MapSet[K, V]) IntersectionWith[S iter.Seq[K] | MapSet[K, V]](keys S) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		s := m.intersection(keys)
		if len(s) < len(m) {
			intersectionWith(m, s)
		}
	case MapSet[K, V]:
		intersectionWith(m, keys)
	}
}

// Intersection returns the ordered key-value pairs which are present in both.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Related:
//   - [MapSet.IntersectionWith] to modify in-place
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) Intersection[S iter.Seq[K] | MapSet[K, V]](keys S) iter.Seq2[K, V] {
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

// IntersectionCount returns the number of keys present in both.
// Use [maps.Keys] on the smaller of maps with different value types.
// Sequences count duplicate keys.
//
// Related:
//   - [MapSet.Overlap] for distinct count from a sequence
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) IntersectionCount[S iter.Seq[K] | MapSet[K, V]](keys S) int {
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
	return m.intersectionCount(it)
}

// Difference returns the key-value pairs which are not present in the keys.
//
// Related:
//   - [MapSet.DifferenceWith] to modify in-place
//   - [MapSet.ReverseDifference] with [maps.Keys] for maps with different value types
//   - [Difference] if the receiver was not a map
//
// Performance:
//   - time: O(m+k) if seq
//   - time: O(m) if map
//   - space: O(min(m,k)) if seq
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

// DifferenceWith removes keys present in keys.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Related:
//   - [MapSet.Difference] to not modify in-place
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) DifferenceWith[S iter.Seq[K] | MapSet[K, V]](keys S) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		m.remove(keys)
	case MapSet[K, V]:
		if len(m) < len(keys) {
			maps.DeleteFunc(m, func(key K, _ V) bool { return keys.Contains(key) })
		} else {
			m.remove(maps.Keys(keys))
		}
	}
}

// SymmetricDifference returns keys which are not in both.
// Use [maps.Keys] for maps with different value types.
//
// Related:
//   - [MapSet.SymmetricDifferenceWith] to modify in-place
//
// Performance:
//   - time: O(m+k)
//   - space: O(min(m, k)) if seq
func (m MapSet[K, V]) SymmetricDifference[S iter.Seq[K] | MapSet[K, V]](keys S) (it iter.Seq[K]) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		it = m.symmetricDifference(keys)
	case MapSet[K, V]:
		it = func(yield func(K) bool) {
			both := 0
			for key := range m {
				if keys.Contains(key) {
					both += 1
				} else if !yield(key) {
					return
				}
			}
			if both < len(keys) {
				for key := range keys {
					if !m.Contains(key) && !yield(key) {
						return
					}
				}
			}
		}
	}
	return it
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
func (m MapSet[K, V]) Union[S iter.Seq2[K, V] | MapSet[K, V]](keys S) MapSet[K, V] {
	m = maps.Clone(m)
	if m == nil {
		m = map[K]V{}
	}
	m.UnionWith(keys)
	return m
}

// UnionWith merges keys with inserts.
// Duplicate keys overwrite values.
//
// Related:
//   - [MapSet.Union] to not modify in-place
//   - [maps.Insert] for an iter.Seq2
//   - [maps.Copy] for a map
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) UnionWith[S iter.Seq2[K, V] | MapSet[K, V]](keys S) {
	switch keys := any(keys).(type) {
	case iter.Seq2[K, V]:
		maps.Insert(m, keys)
	case MapSet[K, V]:
		maps.Copy(m, keys)
	}
}

// Overlap returns the sizes of the intersection and differences:
// left only, both, right only.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Similarity measures:
//   - overlap coefficient: both / (min(left, right) + both)
//   - Jaccard index: both / (left + both + right)
//
// Related:
//   - [MapSet.IntersectionCount] for just the intersection size
//
// Performance:
//   - time: Θ(k) if seq
//   - time: O(min(m, k)) if map
//   - space: Θ(k) if seq
func (m MapSet[K, V]) Overlap[S iter.Seq[K] | MapSet[K, V]](keys S) (left, both, right int) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		left, both, right = m.overlap(keys)
	case MapSet[K, V]:
		both = m.IntersectionCount(keys)
		left, right = len(m)-both, len(keys)-both
	}
	return left, both, right
}
