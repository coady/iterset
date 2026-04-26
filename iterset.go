// Package iterset is a set library based on maps and iterators.
package iterset

import (
	"iter"
	"maps"
)

func allFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for value := range seq {
		if !f(value) {
			return false
		}
	}
	return true
}

// MapSet is a `map` extended with set methods.
//
// Set operations are key based; values are retained but otherwise ignored.
type MapSet[K comparable, V any] map[K]V

func (m MapSet[K, V]) add(key K) {
	var value V
	m[key] = value
}

func (m MapSet[K, V]) pop(key K) bool {
	defer delete(m, key)
	return m.Contains(key)
}

func (m MapSet[K, V]) intersect(keys iter.Seq[K]) MapSet[K, struct{}] {
	s := Set[K]()
	for key := range keys {
		if m.Contains(key) {
			s.add(key)
		}
		if len(m) == len(s) {
			break
		}
	}
	return s
}

// Contains returns whether the key is present.
//
// Related:
//   - [MapSet.IsSuperset] for multiple keys
func (m MapSet[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}

// Missing returns whether the key is not present.
// Negation of [MapSet.Contains]; useful to pass as a bound method.
//
// Related:
//   - [MapSet.IsDisjoint] for multiple keys
func (m MapSet[K, V]) Missing(key K) bool {
	_, ok := m[key]
	return !ok
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
		s.add(key)
		return m.Contains(key)
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

// IsSuperset returns whether all keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsSuperset(keys iter.Seq[K]) bool {
	return allFunc(keys, m.Contains)
}

// IsDisjoint returns whether no keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsDisjoint(keys iter.Seq[K]) bool {
	return len(m) == 0 || allFunc(keys, m.Missing)
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
		if len(m) == 0 {
			return
		}
	}
}

// Remove keys.
//
// Related:
//   - [MapSet.Difference] to not modify in-place
func (m MapSet[K, V]) Remove(keys iter.Seq[K]) {
	for key := range keys {
		delete(m, key)
		if len(m) == 0 {
			return
		}
	}
}

// Toggle removes present keys, and inserts missing keys.
//
// Related:
//   - [MapSet.SymmetricDifference] to not modify in-place
func (m MapSet[K, V]) Toggle(keys iter.Seq[K], value V) {
	for key := range keys {
		if m.Contains(key) {
			delete(m, key)
		} else {
			m[key] = value
		}
	}
}

// Union merges all keys with successive inserts.
// Duplicate keys overwrite values.
//
// Related:
//   - [maps.Insert] to modify in-place
//   - [SortedUnion] for sorted sequences
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
	if len(m) == 0 {
		return maps.All(m)
	}
	return func(yield func(K, V) bool) {
		for key := range keys {
			value, ok := m[key]
			if ok && !yield(key, value) {
				return
			}
		}
	}
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
	return func(yield func(K, V) bool) {
		for key, value := range m {
			if s.Missing(key) && !yield(key, value) {
				return
			}
		}
	}
}

// ReverseDifference returns the ordered keys which are not present in the map.
// Also known as the relative complement.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) ReverseDifference(keys iter.Seq[K]) iter.Seq[K] {
	if len(m) == 0 {
		return keys
	}
	return func(yield func(K) bool) {
		keys(func(key K) bool { return m.Contains(key) || yield(key) })
	}
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
	if len(m) == 0 {
		return keys
	}
	s := Set[K]()
	return func(yield func(K) bool) {
		for key := range keys {
			if m.Contains(key) {
				s.add(key)
			} else if !yield(key) {
				return
			}
		}
		if len(m) == len(s) {
			return
		}
		for key := range m {
			if s.Missing(key) && !yield(key) {
				return
			}
		}
	}
}

// Overlap returns the sizes of the intersection and differences:
// left only, both, right only.
//
// Similarity measures:
//   - overlap coefficient: both / (min(left, right) + both)
//   - Jaccard index: both / (left + both + right)
//
// Performance:
//   - time: Θ(k)
//   - space: Θ(k)
func (m MapSet[K, V]) Overlap(keys iter.Seq[K]) (int, int, int) {
	inter, diff := Set[K](), Set[K]()
	for key := range keys {
		if m.Contains(key) {
			inter.add(key)
		} else {
			diff.add(key)
		}
	}
	return len(m) - len(inter), len(inter), len(diff)
}

// Cast returns a zero-copy [MapSet].
// Equivalent to `MapSet[K, V](m)` without having to specify concrete types.
//
// An instantiated type alias would have the same functionality.
// Methods can also be called as unbound functions: `MapSet[K, V].Method(m, ...)`.
func Cast[K comparable, V any](m map[K]V) MapSet[K, V] {
	return m
}

// Collect returns unique keys with a default value.
// Equivalent to [Set] when value is `struct{}{}`.
//
// Related:
//   - [maps.Collect] for an iter.Seq2
func Collect[K comparable, V any](keys iter.Seq[K], value V) MapSet[K, V] {
	m := MapSet[K, V]{}
	m.Insert(keys, value)
	return m
}

// Set returns unique keys with an empty struct value.
//
// Related:
//   - [Collect] for an iter.Seq
func Set[K comparable](keys ...K) MapSet[K, struct{}] {
	s := make(MapSet[K, struct{}], len(keys))
	s.Add(keys...)
	return s
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
		if m.Missing(key) {
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

// Group returns values grouped by keys.
//
// Related:
//   - [GroupBy] for key function
func Group[K comparable, V any](seq iter.Seq2[K, V]) MapSet[K, []V] {
	m := MapSet[K, []V]{}
	for key, value := range seq {
		m[key] = append(m[key], value)
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

// Reduce combines values grouped by keys with binary function.
//
// Related:
//   - [Group] to collect into a slice
func Reduce[K comparable, V any](seq iter.Seq2[K, V], f func(V, V) V) MapSet[K, V] {
	m := MapSet[K, V]{}
	for key, value := range seq {
		v, ok := m[key]
		if ok {
			value = f(v, value)
		}
		m[key] = value
	}
	return m
}

// Memoize caches function call.
func Memoize[K comparable, V any](keys iter.Seq[K], f func(K) V) MapSet[K, V] {
	m := MapSet[K, V]{}
	for key := range keys {
		if m.Missing(key) {
			m[key] = f(key)
		}
	}
	return m
}
