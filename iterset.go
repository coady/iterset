// Package iterset is a set library based on maps and iterators.
package iterset

import (
	"cmp"
	"context"
	"iter"
	"maps"
	"slices"
)

func allFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for value := range seq {
		if !f(value) {
			return false
		}
	}
	return true
}

func difference[K comparable](keys, seq iter.Seq[K]) iter.Seq[K] {
	s := Set[K]()
	return func(yield func(K) bool) {
		next, stop := iter.Pull(seq)
		defer stop()
		k, ok := next()
		for key := range keys {
			for ok && s.Missing(key) {
				s.add(k)
				k, ok = next()
			}
			if s.Missing(key) && !yield(key) {
				return
			}
		}
	}
}

type zipSource struct {
	index int8 // which sequence the value is from
	empty bool // whether the other sequence is empty
}

func zip[K comparable](keys, seq iter.Seq[K]) iter.Seq2[K, zipSource] {
	return func(yield func(K, zipSource) bool) {
		next, stop := iter.Pull(seq)
		defer stop()
		for key := range keys {
			k, ok := next()
			if !yield(key, zipSource{empty: !ok}) || (ok && !yield(k, zipSource{index: 1})) {
				return
			}
		}
		source := zipSource{index: 1, empty: true}
		for k, ok := next(); ok && yield(k, source); k, ok = next() {
		}
	}
}

func intersect[K comparable](keys, seq iter.Seq[K]) iter.Seq[K] {
	return func(yield func(K) bool) {
		sets := [2]MapSet[K, struct{}]{Set[K](), Set[K]()}
		for key, source := range zip(keys, seq) {
			if sets[1-source.index].pop(key) {
				if !yield(key) {
					return
				}
			} else if !source.empty {
				sets[source.index].add(key)
			} else if len(sets[1-source.index]) == 0 {
				return
			}
		}
	}
}

// MapSet is a `map` extended with set methods.
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

// Equal returns whether the sets of keys are equal.
//
// Related:
//   - [MapSet.Equal] if either sequence was a map
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func Equal[K comparable](keys, seq iter.Seq[K]) bool {
	sets := [3]MapSet[K, struct{}]{Set[K](), Set[K](), Set[K]()}
	for key, source := range zip(keys, seq) {
		if sets[1-source.index].pop(key) {
			sets[2].add(key)
		} else if sets[2].Missing(key) {
			sets[source.index].add(key)
		}
		if source.empty && len(sets[source.index]) > 0 {
			return false
		}
	}
	return len(sets[0])+len(sets[1]) == 0
}

// EqualCounts returns whether the multisets of keys are equal.
//
// Related:
//   - [Equal] to ignore counts
//   - [Count] and [maps.Equal] if either sequence were counts
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func EqualCounts[K comparable](keys, seq iter.Seq[K]) bool {
	m := MapSet[K, int]{}
	for key, source := range zip(keys, seq) {
		if source.empty {
			return false
		}
		m[key] += int(cmp.Or(source.index, -1))
		if m[key] == 0 {
			delete(m, key)
		}
	}
	return len(m) == 0
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
//   - time: O(k)
//   - space: O(k)
func IsSubset[K comparable](keys, seq iter.Seq[K]) bool {
	return IsEmpty(difference(keys, seq))
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

// IsDisjoint returns whether no keys are present in the sequence.
//
// Related:
//   - [MapSet.IsDisjoint] if the sequence was a map
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func IsDisjoint[K comparable](keys, seq iter.Seq[K]) bool {
	return IsEmpty(intersect(keys, seq))
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

// Intersect returns the ordered keys which are present in the sequence(s).
//
// Related:
//   - [MapSet.Intersect] if the sequence was a map
//   - [SortedIntersect] if sequences are sorted
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func Intersect[K comparable](keys iter.Seq[K], seqs ...iter.Seq[K]) iter.Seq[K] {
	for _, seq := range seqs {
		keys = intersect(keys, seq)
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

// Difference returns the ordered keys which are not present in the sequence(s).
//
// Related:
//   - [MapSet.ReverseDifference] if the sequence was a map
//   - [SortedDifference] if sequences are sorted
//
// Performance:
//   - time: O(k)
//   - space: O(k)
func Difference[K comparable](keys iter.Seq[K], seqs ...iter.Seq[K]) iter.Seq[K] {
	for _, seq := range seqs {
		keys = difference(keys, seq)
	}
	return keys
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
	return func(yield func(K) bool) {
		s := Set[K]()
		for key := range keys {
			if s.Missing(key) && !yield(key) {
				return
			}
			s.add(key)
		}
	}
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
	return func(yield func(K, V) bool) {
		s := Set[K]()
		for value := range values {
			k := key(value)
			if s.Missing(k) && !yield(k, value) {
				return
			}
			s.add(k)
		}
	}
}

// Compact returns consecutive runs of deduplicated keys, with counts.
//
// Related:
//   - [Unique] to ignore adjacency
//   - [Count] to return a map
func Compact[K comparable](keys iter.Seq[K]) iter.Seq2[K, int] {
	return func(yield func(K, int) bool) {
		var current K
		count := 0
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
	return func(yield func(K, []V) bool) {
		var current K
		var group []V
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

// CompareValues returns a function which compares by value.
//
// Related:
//   - [Sorted] to sort by value
//   - [slices] functions with a custom [cmp.Compare]
func CompareValues[K comparable, V cmp.Ordered](m map[K]V) func(K, K) int {
	return func(a, b K) int { return cmp.Compare(m[a], m[b]) }
}

// Sorted returns keys ordered by corresponding value.
//
// Related:
//   - [Index] to retain original order
//   - [Count] to rank by frequency
//   - [slices.SortedFunc] with [CompareValues]
func Sorted[K comparable, V cmp.Ordered](m map[K]V) []K {
	return slices.SortedFunc(maps.Keys(m), CompareValues(m))
}

func minFunc[K any, V comparable](seq iter.Seq2[K, V], less func(V, V) bool) []K {
	keys := []K{}
	var current V
	for key, value := range seq {
		if len(keys) == 0 || less(value, current) {
			keys, current = []K{key}, value
		} else if value == current {
			keys = append(keys, key)
		}
	}
	return keys
}

// Min returns the key(s) with the minimum corresponding value.
// Will be empty only if the sequence is empty.
//
// Related:
//   - [Count] to rank by frequency
//   - [slices.MinFunc] with [CompareValues]
func Min[K any, V cmp.Ordered](seq iter.Seq2[K, V]) []K {
	return minFunc(seq, cmp.Less)
}

// Max returns the key(s) with the maximum corresponding value.
// Will be empty only if the sequence is empty.
//
// Related:
//   - [Count] to rank by frequency
//   - [slices.MaxFunc] with [CompareValues]
func Max[K any, V cmp.Ordered](seq iter.Seq2[K, V]) []K {
	return minFunc(seq, func(a, b V) bool { return cmp.Less(b, a) })
}

// IsEmpty returns where there are no values in a sequence.
func IsEmpty[V any](seq iter.Seq[V]) bool {
	for range seq {
		return false
	}
	return true
}

// Size returns the number of values in a sequence.
func Size[V any](seq iter.Seq[V]) int {
	count := 0
	for range seq {
		count += 1
	}
	return count
}

// Keys returns the keys from a sequence of pairs.
//
// Related:
//   - [maps.Keys] for a map
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		seq(func(key K, _ V) bool { return yield(key) })
	}
}

// SortedUnion returns the union of sorted keys.
//
// Related:
//   - [Compact] to deduplicate
//
// Performance:
//   - time: O(k)
func SortedUnion[K cmp.Ordered](keys, seq iter.Seq[K]) iter.Seq[K] {
	return sortedUnionFunc(keys, seq, cmp.Compare)
}

func sortedUnionFunc[V any](keys, values iter.Seq[V], compare func(V, V) int) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull(values)
		defer stop()
		value, ok := next()
		for key := range keys {
			for ok && compare(key, value) > 0 {
				if !yield(value) {
					return
				}
				value, ok = next()
			}
			if !yield(key) {
				return
			}
		}
		for ok && yield(value) {
			value, ok = next()
		}
	}
}

// SortedIntersect returns the intersection of sorted keys.
//
// Performance:
//   - time: O(k)
func SortedIntersect[K cmp.Ordered](keys, seq iter.Seq[K]) iter.Seq[K] {
	return Keys(sortedIntersectFunc(keys, seq, cmp.Compare))
}

func sortedIntersectFunc[K, V any](
	keys iter.Seq[K], values iter.Seq[V], compare func(K, V) int,
) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull(values)
		defer stop()
		value, ok := next()
		for key := range keys {
			c := 1
			for ok {
				c = compare(key, value)
				if c <= 0 {
					break
				}
				value, ok = next()
			}
			if !ok || (c == 0 && !yield(key, value)) {
				return
			}
		}
	}
}

// SortedDifference returns the difference of sorted keys.
//
// Performance:
//   - time: O(k)
func SortedDifference[K cmp.Ordered](keys, seq iter.Seq[K]) iter.Seq[K] {
	return sortedDifferenceFunc(keys, seq, cmp.Compare)
}

func sortedDifferenceFunc[K, V any](
	keys iter.Seq[K], values iter.Seq[V], compare func(K, V) int,
) iter.Seq[K] {
	return func(yield func(K) bool) {
		next, stop := iter.Pull(values)
		defer stop()
		value, ok := next()
		for key := range keys {
			c := 1
			for ok {
				c = compare(key, value)
				if c <= 0 {
					break
				}
				value, ok = next()
			}
			if c != 0 && !yield(key) {
				return
			}
		}
	}
}

func goChan[V any](ctx context.Context, seq iter.Seq[V], size int) <-chan V {
	ch := make(chan V, size)
	go func() {
		defer close(ch)
		for value := range seq {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case ch <- value:
			}
		}
	}()
	return ch
}

// GoIter iterates the sequence in a background goroutine and channel.
// An unbuffered channel (size 0) is sufficient for parallelism,
// but channels introduce overhead. As always, benchmark first.
func GoIter[V any](seq iter.Seq[V], size int) iter.Seq[V] {
	return func(yield func(V) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for value := range goChan(ctx, seq, size) {
			if !yield(value) {
				return
			}
		}
	}
}
