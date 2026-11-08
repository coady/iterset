//go:build !go1.27

package iterset

import (
	"fmt"
	"maps"
	"slices"
)

func ExampleMapSet_Equal() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").Equal(k), Set("a").Equal(k))
	// Output: true false
}

func ExampleMapSet_SubsetOf() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a").SubsetOf(k), Set("a", "c").SubsetOf(k))
	// Output: true false
}

func ExampleMapSet_Intersects() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("c").Intersects(k), Set("a").Intersects(k))
	// Output: false true
}

func ExampleMapSet_IntersectionWith() {
	s := Set("a", "b", "c")
	s.IntersectionWith(slices.Values([]string{"b", "c", "d"}))
	fmt.Println(s)
	// Output: map[b:{} c:{}]
}

func ExampleMapSet_Intersection() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	s := slices.Values([]string{"b", "c"})
	for key, value := range m.Intersection(s) {
		fmt.Println(key, value)
	}
	// Output: b 1
}

func ExampleMapSet_IntersectionCount() {
	s, k := Set("a", "b"), Set("b", "c", "d")
	fmt.Println(s.IntersectionCount(maps.Keys(k)))
	// Output: 1
}

func ExampleMapSet_Difference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(maps.Collect(Set("a", "b").Difference(k)))
	// Output: map[a:{}]
}

func ExampleMapSet_DifferenceWith() {
	s := Set("a", "b")
	s.DifferenceWith(slices.Values([]string{"b", "c"}))
	fmt.Println(s)
	// Output: map[a:{}]
}

func ExampleMapSet_SymmetricDifference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Set("a", "b").SymmetricDifference(k)))
	// Output: [c a]
}

func ExampleMapSet_Union() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	fmt.Println(m.Union(maps.All(map[string]int{"b": 2, "c": 3})))
	// Output: map[a:0 b:2 c:3]
}

func ExampleMapSet_UnionWith() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	m.UnionWith(maps.All(map[string]int{"b": 2, "c": 3}))
	fmt.Println(m)
	// Output: map[a:0 b:2 c:3]
}

func ExampleMapSet_Overlap() {
	s, k := Set("a", "b", "c"), []string{"b", "c", "d"}
	fmt.Println(s.Overlap(slices.Values(k)))
	// Output: 1 2 1
}
