// Package iterset implements set operations using maps and iterators.
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

// Contains returns whether the key is present.
// For multiple keys, is equivalent to [MapSet.IsSuperset].
func (m MapSet[K, V]) Contains(keys ...K) bool {
	return !slices.ContainsFunc(keys, m.missing)
}

// Missing returns whether the key is not present.
// For multiple keys, is equivalent to [MapSet.IsDisjoint].
func (m MapSet[K, V]) Missing(keys ...K) bool {
	return !slices.ContainsFunc(keys, m.contains)
}

// Equal returns whether the key sets are equivalent. See also [maps.Equal].
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

// IsSubset returns whether no keys are missing.
// Note [MapSet.IsSuperset] is more efficient if the keys are from a map.
// [IsSubset] is more efficient if the receiver was not originally a map.
//   - time: Θ(k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) IsSubset(keys iter.Seq[K]) bool {
	s := Collect(filterFunc(keys, m.contains), struct{}{})
	return len(m) == len(s)
}

// IsSubset returns whether all keys are present in the sequence.
// Note [MapSet.IsSuperset] is more efficient if the sequence is from a map.
//   - time: Θ(k)
//   - space: Θ(k)
func IsSubset[K comparable](keys, seq iter.Seq[K]) bool {
	return Collect(seq, struct{}{}).IsSuperset(keys)
}

// IsSuperset returns whether all keys are present.
//   - time: O(k)
func (m MapSet[K, V]) IsSuperset(keys iter.Seq[K]) bool {
	return allFunc(keys, m.contains)
}

// IsDisjoint returns whether no keys are present.
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint(keys iter.Seq[K]) bool {
	return allFunc(keys, m.missing)
}

// Add key(s) with zero value.
func (m MapSet[K, V]) Add(keys ...K) {
	var value V
	for _, key := range keys {
		m[key] = value
	}
}

// Insert keys with default value. See also [maps.Insert].
func (m MapSet[K, V]) Insert(keys iter.Seq[K], value V) {
	for key := range keys {
		m[key] = value
	}
}

// Delete key(s).
func (m MapSet[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// Union merges all keys with successive inserts.
//   - time: Θ(m+k)
//   - space: Θ(m+k)
func (m MapSet[K, V]) Union(seqs ...iter.Seq2[K, V]) MapSet[K, V] {
	m = maps.Clone(m)
	for _, seq := range seqs {
		maps.Insert(m, seq)
	}
	return m
}

// Intersect returns the ordered key-value pairs which are present in both.
//   - time: O(k)
func (m MapSet[K, V]) Intersect(keys iter.Seq[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key := range keys {
			value, ok := m[key]
			if ok && !yield(key, value) {
				return
			}
		}
	}
}

// Intersect returns the ordered keys which are present in the sequence(s).
// Note [MapSet.Intersect] is more efficient if the sequence is from a map.
//   - time: Θ(k)
//   - space: Θ(k)
func Intersect[K comparable](keys iter.Seq[K], seqs ...iter.Seq[K]) iter.Seq[K] {
	for _, seq := range seqs {
		keys = filterFunc(keys, Collect(seq, struct{}{}).contains)
	}
	return keys
}

// Difference returns the key-value pairs which are not present in the keys.
// Note [MapSet.ReverseDifference] is more efficient if the keys are from a map.
// [Difference] is more efficient if the receiver was not originally a map.
//   - time: Ω(k)..O(m+k)
//   - space: O(min(m, k))
func (m MapSet[K, V]) Difference(keys iter.Seq[K]) iter.Seq2[K, V] {
	s := Collect(filterFunc(keys, m.contains), struct{}{})
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
// Note [MapSet.ReverseDifference] is more efficient if the sequence is from a map.
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
// For values that compare equal, the first key-value pair is returned. See [IndexBy].
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
// More efficient than [Unique] or [Count] if the keys are already grouped, e.g., sorted.
//   - time: O(k)
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
// More efficient than [UniqueBy] or [GroupBy] if the values are already grouped, e.g., sorted.
//   - time: O(k)
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

// Collect returns unique keys with a default value. See also [maps.Collect].
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
func Count[K comparable](keys iter.Seq[K]) MapSet[K, int] {
	m := map[K]int{}
	for key := range keys {
		m[key] += 1
	}
	return m
}

// IndexBy returns values indexed by key function.
// If there are collisions, the last value remains. See [GroupBy].
func IndexBy[K comparable, V any](values iter.Seq[V], key func(V) K) MapSet[K, V] {
	m := map[K]V{}
	for value := range values {
		m[key(value)] = value
	}
	return m
}

// GroupBy returns values grouped by key function. See [IndexBy].
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
