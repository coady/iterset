package iterset

import (
	"context"
	"iter"
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestBreak(t *testing.T) {
	k := slices.Values([]string{"a", "A"})
	for range Unique(k) {
		break
	}
	for range UniqueBy(k, strings.ToLower) {
		break
	}
	for range Compact(k) {
		break
	}
	for range CompactBy(k, strings.TrimSpace) {
		break
	}
	for range Set("a").Intersect(k) {
		break
	}
	for range Intersect(k, slices.Values([]string{})) {
	}
	for range Set("b").Difference(k) {
		break
	}
	for range Set("b").SymmetricDifference(k) {
		break
	}
	for c := range Set("b").SymmetricDifference(k) {
		if c == "b" {
			break
		}
	}
	for range SortedUnion(k, slices.Values([]string{})) {
		break
	}
	for range SortedUnion(k, slices.Values([]string{""})) {
		break
	}
	for range SortedIntersect(k, slices.Values([]string{"a"})) {
		break
	}
	for range SortedDifference(k, slices.Values([]string{"a"})) {
		break
	}
	for range GoIter(context.Background(), slices.Values([]string{"", ""}), 0) {
		break
	}
}

func TestExit(t *testing.T) {
	k := slices.Values([]string{"a", "A"})
	for c := range Index(k).Difference(k) {
		t.Errorf("should be empty: %s", c)
	}
	for c := range Index(k).SymmetricDifference(k) {
		t.Errorf("should be empty: %s", c)
	}
	if IsSubset(slices.Values([]string{"b"}), k) {
		t.Errorf("should be false")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for c := range goChan(ctx, k, 1) {
		t.Errorf("should be canceled: %s", c)
	}
}

func TestEmpty(t *testing.T) {
	var m MapSet[string, struct{}]
	if m.Union() == nil {
		t.Error("should not be nil")
	}
	for range m.Intersect(nil) {
		t.Error("should be empty")
	}
	for range Intersect(maps.Keys(m), maps.Keys(m)) {
		t.Error("should be empty")
	}
	for range m.ReverseDifference(maps.Keys(m)) {
		t.Error("should be empty")
	}
	for range m.SymmetricDifference(maps.Keys(m)) {
		t.Error("should be empty")
	}
	m.Delete("")
	m.Remove(slices.Values([]string{""}))
}

func assertMulti[K any](t *testing.T, seq iter.Seq[K]) {
	count := Size(seq)
	if Size(seq) != count {
		t.Error("should not be single-use")
	}
}

func TestMulti(t *testing.T) {
	k := slices.Values([]string{"a", "A", "b"})
	assertMulti(t, difference(k, k))
	assertMulti(t, intersect(k, k))
	assertMulti(t, Unique(k))
	assertMulti(t, Keys(UniqueBy(k, strings.ToLower)))
	assertMulti(t, Keys(Compact(k)))
	assertMulti(t, Keys(CompactBy(k, strings.ToLower)))
	assertMulti(t, Keys(Set("b").Difference(k)))
	assertMulti(t, Set("b").SymmetricDifference(k))
}
