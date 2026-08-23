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

// IsSubset returns whether every map key is present in keys.
//
// Related:
//   - [MapSet.IsSuperset] with [maps.Keys] for maps with different value types
//   - [IsSubset] if the receiver was not a map
//
// Performance:
//   - time: O(k) if seq or slice
//   - time: O(min(m, k)) if map
//   - space: O(min(m, k)) if seq or slice
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
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
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

// Keep only the keys present in both.
// Also known as intersection update.
//
// Related:
//   - [MapSet.Intersect] to not modify in-place
//
// Performance:
//   - time: O(m+k) if seq
//   - time: O(m) if map
//   - space: O(min(m, k)) if seq
func (m MapSet[K, V]) Keep[S iter.Seq[K] | MapSet[K, V]](keys S) {
	switch keys := any(keys).(type) {
	case iter.Seq[K]:
		s := m.intersect(keys)
		if len(s) < len(m) {
			keep(m, s)
		}
	case MapSet[K, V]:
		keep(m, keys)
	}
}

// Intersect returns the ordered key-value pairs which are present in both.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Related:
//   - [MapSet.Keep] to modify in-place
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
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

// IntersectCount returns the number of keys present in both.
// Use [maps.Keys] on the smaller of maps with different value types.
// Sequences count duplicate keys.
//
// Related:
//   - [MapSet.Overlap] for distinct count from a sequence
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) IntersectCount[S iter.Seq[K] | MapSet[K, V]](keys S) int {
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
	return m.intersectCount(it)
}

// Difference returns the key-value pairs which are not present in the keys.
//
// Related:
//   - [MapSet.Remove] to modify in-place
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

// Remove keys.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Related:
//   - [MapSet.Difference] to not modify in-place
//
// Performance:
//   - time: O(k) if seq
//   - time: O(min(m, k)) if map
func (m MapSet[K, V]) Remove[S iter.Seq[K] | MapSet[K, V]](keys S) {
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
//   - [MapSet.Toggle] to modify in-place
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

// Overlap returns the sizes of the intersection and differences:
// left only, both, right only.
// Use [maps.Keys] on the smaller of maps with different value types.
//
// Similarity measures:
//   - overlap coefficient: both / (min(left, right) + both)
//   - Jaccard index: both / (left + both + right)
//
// Related:
//   - [MapSet.IntersectCount] for just the intersection size
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
		both = m.IntersectCount(keys)
		left, right = len(m)-both, len(keys)-both
	}
	return left, both, right
}
