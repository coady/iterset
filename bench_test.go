package iterset

import (
	"iter"
	"math/rand"
	"slices"
	"testing"
)

func setup(b *testing.B) (MapSet[int, struct{}], iter.Seq[int]) {
	defer b.ResetTimer()
	s := Set[int]()
	for range b.N * 2 {
		s.Add(rand.Intn(b.N))
	}
	k := make([]int, b.N/2)
	for i := range k {
		k[i] = rand.Intn(b.N)
	}
	return s, slices.Values(k)
}

func BenchmarkEqual(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.Equal(k)
	}
}

func BenchmarkIsSubset(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.IsSubset(k)
	}
}

func BenchmarkIsSuperset(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		s.IsSuperset(k)
	}
}

func BenchmarkIntersect(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.Intersect(k) {
		}
	}
}

func BenchmarkDifference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.Difference(k) {
		}
	}
}

func BenchmarkReverseDifference(b *testing.B) {
	s, k := setup(b)
	for range b.N {
		for range s.ReverseDifference(k) {
		}
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
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

func BenchmarkCompact(b *testing.B) {
	_, k := setup(b)
	s := slices.Values(slices.Sorted(k))
	b.ResetTimer()
	for range b.N {
		for range Compact(s) {
		}
	}
}
