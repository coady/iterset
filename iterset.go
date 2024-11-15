// Iterset implements set operations using maps and iterators.
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

// allFunc return whether all values pass the predicate function.
func allFunc[E any](values iter.Seq[E], f func(E) bool) bool {
	for value := range values {
		if !f(value) {
			return false
		}
	}
	return true
}

// MapSet is a `map` subtype with set methods.
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

// IsSuperset returns whether all keys are present.
func (m MapSet[K, V]) IsSuperset(keys iter.Seq[K]) bool {
	return allFunc(keys, m.contains)
}

// IsDisjoint returns whether no keys are present.
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

// Delete key(s).
func (m MapSet[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// Intersect returns the ordered key-value pairs which are present in both.
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

// Difference returns the ordered keys which are not present in the map.
// Note the keys are from the "right" parameter. See [MapSet.Delete] for reverse.
func (m MapSet[K, V]) Difference(keys iter.Seq[K]) iter.Seq[K] {
	return filterFunc(keys, m.missing)
}

// Cast returns a zero-copy [MapSet].
func Cast[K comparable, V any](m map[K]V) MapSet[K, V] {
	return m
}

// Unique returns keys in order without duplicates.
func Unique[K comparable](keys iter.Seq[K]) iter.Seq[K] {
	s := Set[K]()
	return filterFunc(keys, func(key K) bool {
		defer s.Add(key)
		return s.missing(key)
	})
}

// Set returns unique keys with an empty struct value.
func Set[K comparable](keys ...K) MapSet[K, struct{}] {
	m := MapSet[K, struct{}]{}
	m.Add(keys...)
	return m
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

// Sorted returns keys sorted by values.
// When used with [Index], this will retain the original key order.
func Sorted[K comparable, V cmp.Ordered](m map[K]V) []K {
	compare := func(a, b K) int { return cmp.Compare(m[a], m[b]) }
	return slices.SortedFunc(maps.Keys(m), compare)
}
