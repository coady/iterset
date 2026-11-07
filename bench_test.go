package iterset

import (
	"iter"
	"maps"
	"math/rand"
	"slices"
	"testing"
)

const size = 100_000

func identity[V any](v V) V { return v }

func setup(b *testing.B) (MapSet[int, struct{}], iter.Seq[int]) {
	s := Set[int]()
	for range size / 4 {
		s.Insert(rand.Intn(size))
	}
	k := make([]int, size/2)
	for i := range k {
		k[i] = rand.Intn(size)
	}
	return s, slices.Values(k)
}

func BenchmarkMapSet_Equal(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.Equal(k)
	}
}

func BenchmarkMapSet_EqualFull(b *testing.B) {
	s, _ := setup(b)
	for b.Loop() {
		s.Equal(maps.Keys(s))
	}
}

func BenchmarkEqual(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		Equal(maps.Keys(s), k)
	}
}

func BenchmarkEqualFull(b *testing.B) {
	s, _ := setup(b)
	for b.Loop() {
		Equal(maps.Keys(s), maps.Keys(s))
	}
}

func BenchmarkEqualSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		Equal(maps.Keys(s), slc)
	}
}

func BenchmarkEqualCounts(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		EqualCounts(maps.Keys(s), k)
	}
}

func BenchmarkEqualCountsSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		EqualCounts(maps.Keys(s), slc)
	}
}

func BenchmarkMapSet_ContainsAll(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.ContainsAll(k)
	}
}

func BenchmarkMapSet_SubsetOf(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.SubsetOf(k)
	}
}

func BenchmarkSubsetOf(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		SubsetOf(maps.Keys(s), k)
	}
}

func BenchmarkSubsetOfSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		SubsetOf(maps.Keys(s), slc)
	}
}

func BenchmarkMapSet_Intersects(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.Intersects(k)
	}
}

func BenchmarkIntersects(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		Intersects(maps.Keys(s), k)
	}
}

func BenchmarkIntersectsSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		Intersects(maps.Keys(s), slc)
	}
}

func BenchmarkMapSet_Intersection(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range s.Intersection(k) {
		}
	}
}

func BenchmarkIntersection(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range Intersection(maps.Keys(s), k) {
		}
	}
}

func BenchmarkIntersectSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		for range Intersection(maps.Keys(s), slc) {
		}
	}
}

func BenchmarkMapSet_Difference(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range s.Difference(k) {
		}
	}
}

func BenchmarkDifference(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range Difference(maps.Keys(s), k) {
		}
	}
}

func BenchmarkDifferenceSlice(b *testing.B) {
	s, k := setup(b)
	slc := slices.Collect(k)
	for b.Loop() {
		for range Difference(maps.Keys(s), slc) {
		}
	}
}

func BenchmarkMapSet_ReverseDifference(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range s.ReverseDifference(k) {
		}
	}
}

func BenchmarkMapSet_SymmetricDifference(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		for range s.SymmetricDifference(k) {
		}
	}
}

func BenchmarkMapSet_Overlap(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.Overlap(k)
	}
}

func BenchmarkMapSet_IntersectionCount(b *testing.B) {
	s, k := setup(b)
	for b.Loop() {
		s.IntersectionCount(k)
	}
}

func BenchmarkUnique(b *testing.B) {
	_, k := setup(b)
	for b.Loop() {
		for range Unique(k) {
		}
	}
}

func BenchmarkUniqueBy(b *testing.B) {
	_, k := setup(b)
	for b.Loop() {
		for range UniqueBy(k, identity[int]) {
		}
	}
}

func BenchmarkCompact(b *testing.B) {
	_, k := setup(b)
	s := slices.Values(slices.Sorted(k))
	for b.Loop() {
		for range Compact(s) {
		}
	}
}

func BenchmarkCompactBy(b *testing.B) {
	_, k := setup(b)
	for b.Loop() {
		for range CompactBy(k, identity[int]) {
		}
	}
}

func BenchmarkSet(b *testing.B) {
	_, k := setup(b)
	s := slices.Collect(k)
	for b.Loop() {
		Set(s...)
	}
}

func BenchmarkIndexBy(b *testing.B) {
	_, k := setup(b)
	for b.Loop() {
		IndexBy(k, identity[int])
	}
}

func BenchmarkGroupBy(b *testing.B) {
	_, k := setup(b)
	for b.Loop() {
		GroupBy(k, identity[int])
	}
}

func BenchmarkSorted(b *testing.B) {
	s, k := setup(b)
	v := slices.Values(slices.Sorted(maps.Keys(s)))
	k = slices.Values(slices.Sorted(k))
	for b.Loop() {
		for range SortedUnion(k, v) {
		}
		for range SortedIntersection(k, v) {
		}
		for range SortedDifference(k, v) {
		}
	}
}

func BenchmarkSortedSlice(b *testing.B) {
	s, k := setup(b)
	v := slices.Sorted(maps.Keys(s))
	k = slices.Values(slices.Sorted(k))
	for b.Loop() {
		for range SortedUnion(k, v) {
		}
		for range SortedIntersection(k, v) {
		}
		for range SortedDifference(k, v) {
		}
	}
}
