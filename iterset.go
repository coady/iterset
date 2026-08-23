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

func (m MapSet[K, V]) intersectCount(keys iter.Seq[K]) int {
	if len(m) == 0 {
		return 0
	}
	count := 0
	for key := range keys {
		if m.Contains(key) {
			count += 1
		}
	}
	return count
}

func (m MapSet[K, V]) equal(keys iter.Seq[K]) bool {
	s := Set[K]()
	superset := allFunc(keys, func(key K) bool {
		s.add(key)
		return m.Contains(key)
	})
	return superset && len(m) == len(s)
}

func (m MapSet[K, V]) filter(f func(K) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, value := range m {
			if f(key) && !yield(key, value) {
				return
			}
		}
	}
}

func (m MapSet[K, V]) difference(keys iter.Seq[K]) iter.Seq2[K, V] {
	s := m.intersect(keys)
	switch len(s) {
	case 0:
		return maps.All(m)
	case len(m):
		return func(func(K, V) bool) {}
	}
	return m.filter(s.Missing)
}

func (m MapSet[K, V]) mask(keys iter.Seq[K]) iter.Seq2[K, V] {
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

func keep[K comparable, V1, V2 any](m MapSet[K, V1], keys MapSet[K, V2]) {
	if len(keys) == 0 {
		clear(m)
	} else {
		maps.DeleteFunc(m, func(key K, _ V1) bool { return !keys.Contains(key) })
	}
}

// Contains returns whether the key is present.
//
// Related:
//   - [MapSet.IsSuperset] for multiple keys
func (m MapSet[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}

// Missing returns whether the key is not present. ![MapSet.Contains] is preferred;
// Missing exists to pass as a function value, e.g. to [slices.DeleteFunc].
//
// Related:
//   - [MapSet.IsDisjoint] for multiple keys
func (m MapSet[K, V]) Missing(key K) bool {
	_, ok := m[key]
	return !ok
}

// IsSuperset returns whether all keys are present.
//
// Performance:
//   - time: O(k)
func (m MapSet[K, V]) IsSuperset(keys iter.Seq[K]) bool {
	return allFunc(keys, m.Contains)
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

func (m MapSet[K, V]) remove(keys iter.Seq[K]) {
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

func (m MapSet[K, V]) symmetricDifference(keys iter.Seq[K]) iter.Seq[K] {
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
			if !s.Contains(key) && !yield(key) {
				return
			}
		}
	}
}

func (m MapSet[K, V]) overlap(keys iter.Seq[K]) (int, int, int) {
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
		if !m.Contains(key) {
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
func IndexBy[K comparable, V any, S iter.Seq[V] | []V](values S, key func(V) K) MapSet[K, V] {
	m := sized[K, V, V](values)
	for value := range sequence[V](values) {
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
func GroupBy[K comparable, V any, S iter.Seq[V] | []V](values S, key func(V) K) MapSet[K, []V] {
	m := sized[K, []V, V](values)
	for value := range sequence[V](values) {
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
func Memoize[K comparable, V any, S iter.Seq[K] | []K](keys S, f func(K) V) MapSet[K, V] {
	m := sized[K, V, K](keys)
	for key := range sequence[K](keys) {
		if !m.Contains(key) {
			m[key] = f(key)
		}
	}
	return m
}
