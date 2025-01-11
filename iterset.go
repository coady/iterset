// Package iterset is a set library based on maps and iterators.
package iterset

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

func filterFunc[E any](values iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for value := range values {
			if f(value) && !yield(value) {
				return
			}
		}
	}
}

func allFunc[E any](values iter.Seq[E], f func(E) bool) bool {
	for value := range values {
		if !f(value) {
			return false
		}
	}
	return true
}

// MapSet is a `map` extended with set methods.
type MapSet[K comparable, V any] map[K]V

// Get returns the key's value. A convenience method that can be passed as an argument.
func (m MapSet[K, V]) Get(key K) V {
	return m[key]
}

func (m MapSet[K, V]) contains(key K) bool {
	_, ok := m[key]
	return ok
}

func (m MapSet[K, V]) missing(key K) bool {
	_, ok := m[key]
	return !ok
}

func (m MapSet[K, V]) intersect(keys iter.Seq[K]) MapSet[K, struct{}] {
	s := Set[K]()
	for key := range keys {
		if m.contains(key) {
			s.Add(key)
		}
		if len(m) == len(s) {
			break
		}
	}
	return s
}

// Contains returns whether the key(s) is present.
//
// Related:
//   - [MapSet.IsSuperset] for many keys
func (m MapSet[K, V]) Contains(keys ...K) bool {
	return !slices.ContainsFunc(keys, m.missing)
}

// Missing returns whether the key(s) is not present.
//
// Related:
//   - [MapSet.IsDisjoint] for many keys
func (m MapSet[K, V]) Missing(keys ...K) bool {
	return len(m) == 0 || !slices.ContainsFunc(keys, m.contains)
}

// Equal returns whether the key sets are equivalent.
//
// Related:
//   - [maps.Equal] to compare values
//
// Performance:
//   - time: O(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) Equal(keys iter.Seq[K]) bool {
	s := Set[K]()
	superset := allFunc(keys, func(key K) bool {
		s.Add(key)
		return m.contains(key)
	})
	return superset && len(m) == len(s)
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

// IsSubset returns whether all keys are present in the sequence.
//
// Related:
//   - [MapSet.IsSuperset] if the sequence was a map
//
// Performance:
//   - time: Θ(k)
//   - space: Θ(k)
func IsSubset[K comparable](keys, seq iter.Seq[K]) bool {
	return Collect(seq, struct{}{}).IsSuperset(keys)
}

// IsSuperset returns whether all keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsSuperset(keys iter.Seq[K]) bool {
	return allFunc(keys, m.contains)
}

// IsDisjoint returns whether no keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint(keys iter.Seq[K]) bool {
	return len(m) == 0 || allFunc(keys, m.missing)
}

// Add key(s) with zero value.
//
// Related:
//   - [MapSet.Insert] for many keys
func (m MapSet[K, V]) Add(keys ...K) {
	var value V
	for _, key := range keys {
		m[key] = value
	}
}

// Insert keys with default value.
//
// Related:
//   - [maps.Insert] for an iter.Seq2
//   - [maps.Copy] for a map
func (m MapSet[K, V]) Insert(keys iter.Seq[K], value V) {
	for key := range keys {
		m[key] = value
	}
}

// Delete key(s).
//
// Related:
//   - [MapSet.Remove] for many keys
func (m MapSet[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// Remove keys.
//
// Related:
//   - [MapSet.Difference] to not modify in-place
func (m MapSet[K, V]) Remove(keys iter.Seq[K]) {
	for key := range keys {
		delete(m, key)
	}
}

// Toggle removes present keys, and inserts missing keys.
//
// Related:
//   - [MapSet.SymmetricDifference] to not modify in-place
func (m MapSet[K, V]) Toggle(seq iter.Seq2[K, V]) {
	for key, value := range seq {
		if m.contains(key) {
			delete(m, key)
		} else {
			m[key] = value
		}
	}
}

// Union merges all keys with successive inserts.
//
// Related:
//   - [maps.Insert] to modify in-place
//
// Performance:
//   - time: Θ(m+k)
//   - space: Ω(max(m, k))..O(m+k)
func (m MapSet[K, V]) Union(seqs ...iter.Seq2[K, V]) MapSet[K, V] {
	m = maps.Clone(m)
	if m == nil {
		m = map[K]V{}
	}
	for _, seq := range seqs {
		maps.Insert(m, seq)
	}
	return m
}

// Intersect returns the ordered key-value pairs which are present in both.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) Intersect(keys iter.Seq[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if len(m) == 0 {
			return
		}
		for key := range keys {
			value, ok := m[key]
			if ok && !yield(key, value) {
				return
			}
		}
	}
}

// Intersect returns the ordered keys which are present in the sequence(s).
//
// Related:
//   - [MapSet.Intersect] if the sequence was a map
//
// Performance:
//   - time: Θ(k)
//   - space: Θ(k)
func Intersect[K comparable](keys iter.Seq[K], seqs ...iter.Seq[K]) iter.Seq[K] {
	for _, seq := range seqs {
		s := Collect(seq, struct{}{})
		if len(s) == 0 {
			return maps.Keys(s)
		}
		keys = filterFunc(keys, s.contains)
	}
	return keys
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
	return func(yield func(K, V) bool) {
		if len(m) == len(s) {
			return
		}
		for key, value := range m {
			if s.missing(key) && !yield(key, value) {
				return
			}
		}
	}
}

// Difference returns the ordered keys which are not present in the sequence(s).
//
// Related:
//   - [MapSet.ReverseDifference] if the sequence was a map
//
// Performance:
//   - time: Θ(k)
//   - space: Θ(k)
func Difference[K comparable](keys iter.Seq[K], seqs ...iter.Seq[K]) iter.Seq[K] {
	for _, seq := range seqs {
		keys = filterFunc(keys, Collect(seq, struct{}{}).missing)
	}
	return keys
}

// ReverseDifference returns the ordered keys which are not present in the map.
// Also known as the relative complement.
//   - time: O(k)
func (m MapSet[K, V]) ReverseDifference(keys iter.Seq[K]) iter.Seq[K] {
	return filterFunc(keys, m.missing)
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
	s := Set[K]()
	return func(yield func(K) bool) {
		for key := range keys {
			if m.contains(key) {
				s.Add(key)
			} else if !yield(key) {
				return
			}
		}
		if len(m) == len(s) {
			return
		}
		for key := range m {
			if s.missing(key) && !yield(key) {
				return
			}
		}
	}
}

// Cast returns a zero-copy [MapSet].
// Equivalent to `MapSet[K, V](m)` without having to specify concrete types.
func Cast[K comparable, V any](m map[K]V) MapSet[K, V] {
	return m
}

// Unique returns keys in order without duplicates.
//
// Related:
//   - [Index] to return a map
//   - [Compact] if the keys are already grouped
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func Unique[K comparable](keys iter.Seq[K]) iter.Seq[K] {
	s := Set[K]()
	return filterFunc(keys, func(key K) bool {
		defer s.Add(key)
		return s.missing(key)
	})
}

// UniqueBy is like [Unique] but uses a key function to compare values.
// For values that compare equal, the first key-value pair is returned.
//
// Related:
//   - [IndexBy] to return a map
//   - [CompactBy] if the values are already grouped by key
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func UniqueBy[K comparable, V any](values iter.Seq[V], key func(V) K) iter.Seq2[K, V] {
	s := Set[K]()
	return func(yield func(K, V) bool) {
		for value := range values {
			k := key(value)
			if s.missing(k) && !yield(k, value) {
				return
			}
			s.Add(k)
		}
	}
}

// Compact returns consecutive runs of deduplicated keys, with counts.
//
// Related:
//   - [Unique] to ignore adjacency
//   - [Count] to return a map
func Compact[K comparable](keys iter.Seq[K]) iter.Seq2[K, int] {
	var current K
	count := 0
	return func(yield func(K, int) bool) {
		for key := range keys {
			if count > 0 && key != current {
				if !yield(current, count) {
					return
				}
				count = 0
			}
			current = key
			count += 1
		}
		if count > 0 {
			yield(current, count)
		}
	}
}

// CompactBy is like [Compact] but uses a key function and collects all values.
//
// Related:
//   - [UniqueBy] to ignore adjacency
//   - [GroupBy] to return a map
func CompactBy[K comparable, V any](values iter.Seq[V], key func(V) K) iter.Seq2[K, []V] {
	var current K
	var group []V
	return func(yield func(K, []V) bool) {
		for value := range values {
			k := key(value)
			if group != nil && k != current {
				if !yield(current, group) {
					return
				}
				group = nil
			}
			current = k
			group = append(group, value)
		}
		if group != nil {
			yield(current, group)
		}
	}
}

// Collect returns unique keys with a default value.
//
// Related:
//   - [maps.Collect] for an iter.Seq2
func Collect[K comparable, V any](keys iter.Seq[K], value V) MapSet[K, V] {
	m := MapSet[K, V]{}
	m.Insert(keys, value)
	return m
}

// Set returns unique keys with an empty struct value.
func Set[K comparable](keys ...K) MapSet[K, struct{}] {
	return Collect(slices.Values(keys), struct{}{})
}

// Index returns unique keys with their first index position.
//
// Related:
//   - [Unique] to return an ordered sequence
//   - [Sorted] to restore original order
func Index[K comparable](keys iter.Seq[K]) MapSet[K, int] {
	m := MapSet[K, int]{}
	i := 0
	for key := range keys {
		if m.missing(key) {
			m[key] = i
		}
		i += 1
	}
	return m
}

// Count returns unique keys with their counts.
//
// Related:
//   - [Compact] if the keys are already grouped
func Count[K comparable](keys iter.Seq[K]) MapSet[K, int] {
	m := map[K]int{}
	for key := range keys {
		m[key] += 1
	}
	return m
}

// IndexBy returns values indexed by key function.
// If there are collisions, the last value remains.
//
// Related:
//   - [UniqueBy] to return an ordered sequence
//   - [GroupBy] to retain all values
func IndexBy[K comparable, V any](values iter.Seq[V], key func(V) K) MapSet[K, V] {
	m := map[K]V{}
	for value := range values {
		m[key(value)] = value
	}
	return m
}

// GroupBy returns values grouped by key function.
//
// Related:
//   - [IndexBy] to retain single value
//   - [CompactBy] if the values are already grouped by key
func GroupBy[K comparable, V any](values iter.Seq[V], key func(V) K) MapSet[K, []V] {
	m := map[K][]V{}
	for value := range values {
		k := key(value)
		m[k] = append(m[k], value)
	}
	return m
}

// Memoize caches function call.
func Memoize[K comparable, V any](keys iter.Seq[K], f func(K) V) MapSet[K, V] {
	m := MapSet[K, V]{}
	for key := range keys {
		if m.missing(key) {
			m[key] = f(key)
		}
	}
	return m
}

// Sorted returns keys sorted by values.
// When used with [Index], this will retain the original key order.
func Sorted[K comparable, V cmp.Ordered](m map[K]V) []K {
	compare := func(a, b K) int { return cmp.Compare(m[a], m[b]) }
	return slices.SortedFunc(maps.Keys(m), compare)
}
