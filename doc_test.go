package iterset

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
)

type item struct{ id string }

// Intersect a map with a slice of keys, retaining original order.
func Example_map() {
	data := map[string]int{"a": 0, "b": 1, "c": 2}
	keys := []string{"d", "c", "b"}
	// With no sets.
	for _, key := range keys {
		value, ok := data[key]
		if ok {
			fmt.Println(key, value)
		}
	}
	// Using a typical `mapset` would be inefficient and only change the `Contains` lookup.
	// Whereas `iterset` methods have in-lined logic with zero copying.
	for key, value := range Cast(data).Intersect(slices.Values(keys)) {
		fmt.Println(key, value)
	}
	// Output:
	// c 2
	// b 1
	// c 2
	// b 1
}

// Difference between two slices while, retaining original order.
func Example_slice() {
	x, y := []string{"a", "b", "c"}, []string{"e", "d", "c"}
	// With no sets.
	m := map[string]bool{}
	for _, c := range x {
		m[c] = true
	}
	for _, c := range y {
		if !m[c] {
		}
	}
	// Using a typical `mapset` only solves half the problem.
	s := Set(x...)
	for _, c := range y {
		if !s.Contains(c) {
		}
	}
	// Whereas `iterset` methods have in-lined logic with minimal copying.
	for c := range Set(x...).Difference(slices.Values(y)) {
		fmt.Println(c)
	}
	// Output:
	// e
	// d
}

func ExampleMapSet_Get() {
	m := Cast(map[string]int{"a": 1})
	fmt.Println(m.Get("a"))
	// Output: 1
}

func ExampleMapSet_Contains() {
	s := Set("b", "a", "b")
	fmt.Println(s.Contains("a"), s.Contains("b", "c"))
	// Output: true false
}

func ExampleMapSet_Missing() {
	s := Set("b", "a", "b")
	fmt.Println(s.Missing("c"), s.Missing("b", "c"))
	// Output: true false
}

func ExampleMapSet_Equal() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").Equal(k), Set("a").Equal(k))
	// Output: true false
}

func ExampleMapSet_IsSubset() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a").IsSubset(k), Set("a", "c").IsSubset(k))
	// Output: true false
}

func ExampleMapSet_IsSuperset() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").IsSuperset(k), Set("a").IsSuperset(k))
	// Output: true false
}

func ExampleMapSet_IsDisjoint() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("c").IsDisjoint(k), Set("a").IsDisjoint(k))
	// Output: true false
}

func ExampleMapSet_Add() {
	s := Set("a", "b")
	s.Add("b", "c")
	fmt.Println(len(s))
	// Output: 3
}

func ExampleMapSet_Delete() {
	s := Set("a", "b")
	s.Delete("b", "c")
	fmt.Println(len(s))
	// Output: 1
}

func ExampleMapSet_Union() {
	m := map[string]int{"a": 0, "b": 1}
	n := map[string]int{"b": 2, "c": 3}
	fmt.Println(Cast(m).Union(maps.All(n)))
	// Output: map[a:0 b:2 c:3]
}

func ExampleMapSet_SymmetricDifference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Set("a", "b").SymmetricDifference(k)))
	// Output: [c a]
}

func ExampleUnique() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(slices.Collect(Unique(k)))
	// Output: [b a]
}

func ExampleCollect() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Collect(k, true))
	// Output: map[a:true b:true]
}

func ExampleIndex() {
	fmt.Println(Index(slices.Values([]string{"b", "a", "b"})))
	// Output: map[a:1 b:0]
}

func ExampleCount() {
	fmt.Println(Count(slices.Values([]string{"b", "a", "b"})))
	// Output: map[a:1 b:2]
}

func ExampleIndexBy() {
	items := []item{{id: "b"}, {id: "a"}, {id: "b"}}
	fmt.Println(IndexBy(slices.Values(items), func(it item) string { return it.id }))
	// Output: map[a:{a} b:{b}]
}

func ExampleGroupBy() {
	items := []item{{id: "b"}, {id: "a"}, {id: "b"}}
	fmt.Println(GroupBy(slices.Values(items), func(it item) string { return it.id }))
	// Output: map[a:[{a}] b:[{b} {b}]]
}

func ExampleMemoize() {
	fmt.Println(Memoize(slices.Values([]string{"b", "a", "b"}), strings.ToUpper))
	// Output: map[a:A b:B]
}

func ExampleSorted() {
	fmt.Println(Sorted(Index(slices.Values([]string{"b", "a", "b"}))))
	// Output: [b a]
}

func TestIter(t *testing.T) {
	k := slices.Values([]string{"a"})
	for range Set("a").Intersect(k) {
		break
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
}
