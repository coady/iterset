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
	defer b.ResetTimer()
	s := Set[int]()
	for range size / 4 {
		s.Add(rand.Intn(size))
	}
	k := make([]int, size/2)
	for i := range k {
		k[i] = rand.Intn(size)
	}
	return s, slices.Values(k)
}

func BenchmarkMapSet_Equal(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.Equal(k)
	}
}

func BenchmarkEqual(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		Equal(maps.Keys(s), k)
	}
}
func BenchmarkEqualCounts(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		EqualCounts(maps.Keys(s), k)
	}
}

func BenchmarkMapSet_IsSubset(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.IsSubset(k)
	}
}

func BenchmarkIsSubset(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		IsSubset(maps.Keys(s), k)
	}
}

func BenchmarkMapSet_IsSuperset(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.IsSuperset(k)
	}
}

func BenchmarkMapSet_IsDisjoint(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.IsDisjoint(k)
	}
}

func BenchmarkIsDisjoint(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		IsDisjoint(maps.Keys(s), k)
	}
}

func BenchmarkMapSet_Intersect(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.Intersect(k) {
		}
	}
}

func BenchmarkIntersect(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range Intersect(maps.Keys(s), k) {
		}
	}
}

func BenchmarkMapSet_Difference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.Difference(k) {
		}
	}
}

func BenchmarkDifference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range Difference(maps.Keys(s), k) {
		}
	}
}

func BenchmarkMapSet_ReverseDifference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.ReverseDifference(k) {
		}
	}
}

func BenchmarkMapSet_SymmetricDifference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.SymmetricDifference(k) {
		}
	}
}

func BenchmarkUnique(b *testing.B) {
	_, k := setup(b)
	for range b.N {
		for range Unique(k) {
		}
	}
}

func BenchmarkUniqueBy(b *testing.B) {
	_, k := setup(b)
	for range b.N {
		for range UniqueBy(k, identity) {
		}
	}
}

func BenchmarkCompact(b *testing.B) {
	_, k := setup(b)
	s := slices.Values(slices.Sorted(k))
	b.ResetTimer()
	for range b.N {
		for range Compact(s) {
		}
	}
}

func BenchmarkCompactBy(b *testing.B) {
	_, k := setup(b)
	for range b.N {
		for range CompactBy(k, identity) {
		}
	}
}

func BenchmarkSet(b *testing.B) {
	_, k := setup(b)
	s := slices.Collect(k)
	b.ResetTimer()
	for range b.N {
		Set(s...)
	}
}

func BenchmarkIndexBy(b *testing.B) {
	_, k := setup(b)
	for range b.N {
		IndexBy(k, identity)
	}
}

func BenchmarkGroupBy(b *testing.B) {
	_, k := setup(b)
	for range b.N {
		GroupBy(k, identity)
	}
}

func BenchmarkSorted(b *testing.B) {
	s, k := setup(b)
	v := slices.Values(slices.Sorted(maps.Keys(s)))
	k = slices.Values(slices.Sorted(k))
	b.ResetTimer()
	for range b.N {
		for range SortedUnion(k, v) {
		}
		for range SortedIntersect(k, v) {
		}
		for range SortedDifference(k, v) {
		}
	}
}
